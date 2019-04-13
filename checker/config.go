package checker

import (
	"time"

	"github.com/leonklingele/github-release-checker/checker/github"
)

type Config struct {
	CheckInterval duration `toml:"interval"`
	Workers       int      `toml:"workers"`

	DBConfig           *DBConfig           `toml:"db"`
	RepositoriesConfig *RepositoriesConfig `toml:"repositories"`
	GithubConfig       *github.Config      `toml:"github"`
}

type DBConfig struct {
	Path string `toml:"path"`
}

type RepositoriesConfig struct {
	Ignored   []string `toml:"ignored"`
	Important []string `toml:"important"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func (d duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}
