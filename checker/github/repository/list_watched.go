package repository

import (
	"context"
	"sync"

	"github.com/google/go-github/github"
	"github.com/leonklingele/github-release-checker/logging"
	"github.com/pkg/errors"
)

func listWatchedAll(
	activity *github.ActivityService,
	repoChan chanW,
	wg *sync.WaitGroup,
) error {
	firstPage := 1
	res, err := listWatched(activity, firstPage, repoChan)
	if err != nil {
		return errors.Wrap(err, "failed to get initial watched page")
	}
	for page := firstPage + 1; page <= res.LastPage; page++ {
		wg.Add(1)
		go func(page int) {
			defer func() {
				wg.Done()
				logging.Debug("done getting watched repos, page", page)
			}()
			logging.Debug("start getting watched repos, page", page)
			if _, err := listWatched(activity, page, repoChan); err != nil {
				logging.Error(errors.Wrapf(err, "failed to get watched repos, page %d", page))
			}
		}(page)
	}
	return nil
}

func listWatched(
	activity *github.ActivityService,
	page int,
	repoChan chanW,
) (*github.Response, error) {
	opt := &github.ListOptions{
		Page:    page,
		PerPage: perPage,
	}
	// TODO(leon): Pass in context
	watched, res, err := activity.ListWatched(context.TODO(), "", opt)
	if err != nil {
		return nil, err
	}
	for _, repo := range watched {
		repoChan <- newRepository(repo)
	}
	return res, nil
}
