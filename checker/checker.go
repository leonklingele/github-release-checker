package checker

import (
	"database/sql"
	"os"
	"time"

	"github.com/leonklingele/github-release-checker/checker/github"
	"github.com/leonklingele/github-release-checker/checker/github/repository"
	"github.com/leonklingele/github-release-checker/checker/github/tag"
	"github.com/leonklingele/github-release-checker/checker/handlers"
	"github.com/leonklingele/github-release-checker/logging"
	"github.com/leonklingele/github-release-checker/pathutil"
	"github.com/leonklingele/github-release-checker/stringutil"
	"github.com/pkg/errors"
)

const (
	sqliteInitStmt = `
		begin;
		create table tags (
			id integer primary key not null,
			repo text,
			version text,
			found datetime default current_timestamp
		);
		create unique index unique_repo_version on tags (repo, version);
		commit;
	`
	sqliteInsertStmt = `
		insert into tags (
			repo, version
		) values (
			?, ?
		);
	`
)

type Checker struct {
	config *Config

	github *github.Github
	db     *sql.DB
}

func (c *Checker) Start() error {
	stmt, err := c.db.Prepare(sqliteInsertStmt)
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement")
	}

	runCheck := func() {
		done := make(chan struct{})
		go func() {
			start := time.Now()
			logging.Info("start running check")
			<-done
			delta := time.Since(start)
			logging.Infof("done running check in %fs", delta.Seconds())
		}()

		activity := c.github.Client.Activity
		listRepos := repository.List
		// Make sure to first annotate, then cleanup
		annotateRepos := repository.MakeChanFilter(func(repo *repository.Repository) bool {
			repo.IsIgnored = c.isIgnoredRepository(repo)
			repo.IsImportant = c.isImportantRepository(repo)
			return true
		})
		cleanupRepos := repository.CleanupChan
		listTags := tag.List
		filterTags := tag.MakeChanFilter(func(tag *tag.Tag) bool {
			name := tag.Repository.GetFullName()
			version := tag.Version
			if _, err := stmt.Exec(name, version); err != nil {
				// TODO(leon): Should only use .Debug if is unique constraint error
				logging.Debug(errors.Wrapf(err, "failed to insert row: (%s, %s)", name, version))
				return false
			}
			return true
		})

		// TODO(leon): This looks a bit ugly. Use pipeline slice instead?
		tagChan :=
			filterTags(
				listTags(
					cleanupRepos(
						annotateRepos(
							listRepos(activity),
						),
					),
					c.config.Workers,
				),
			)
		if err := handlers.Handle(tagChan, done); err != nil {
			logging.Error(errors.Wrap(err, "failed to handle tag chan"))
		}
	}

	if c.config.CheckInterval.Nanoseconds() > 0 {
		ticker := time.NewTicker(c.config.CheckInterval.Duration)
		// Defer ticker listener to ensure this order:
		// 1. Start the timer
		// 2. Run a timerless check
		// 3. Run timer checks
		defer func() {
			for range ticker.C {
				runCheck()
			}
		}()
	}
	runCheck()

	return nil
}

func (c *Checker) isIgnoredRepository(repo *repository.Repository) bool {
	return stringutil.SliceContainsCaseInsensitive(c.config.RepositoriesConfig.Ignored, repo.GetFullName())
}

func (c *Checker) isImportantRepository(repo *repository.Repository) bool {
	return stringutil.SliceContainsCaseInsensitive(c.config.RepositoriesConfig.Important, repo.GetFullName())
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(sqliteInitStmt)
	return err
}

func New(c *Config) (*Checker, error) {
	dbp, err := pathutil.ReplaceHome(c.DBConfig.Path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to replace home")
	}
	db, err := sql.Open("sqlite3", dbp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	// defer db.Close() // TODO(leon): Do we even need to close?

	if _, err := os.Stat(dbp); os.IsNotExist(err) {
		if err := initDB(db); err != nil {
			return nil, errors.Wrap(err, "failed to init database")
		}
		logging.Info("successfully initialized database")
	}

	return &Checker{
		config: c,
		github: github.New(c.GithubConfig),
		db:     db,
	}, nil
}
