package repository

import (
	"sync"

	"github.com/google/go-github/github"
	"github.com/leonklingele/github-release-checker/logging"
	"github.com/pkg/errors"
)

const (
	perPage = 1000
)

type checkFunc func(
	*github.ActivityService,
	chanW,
	*sync.WaitGroup,
) error

func List(activity *github.ActivityService) Chan {
	repoChan := make(chanRW)

	go func() {
		checks := map[string]checkFunc{
			"listStarred": listStarredAll,
			"listWatched": listWatchedAll,
		}
		wg := &sync.WaitGroup{}
		defer func() {
			wg.Wait()
			close(repoChan)
			logging.Debug("done listing repos")
		}()

		for name, check := range checks {
			wg.Add(1)
			go func(name string, check checkFunc) {
				defer wg.Done()
				if err := check(activity, onlyWritable(repoChan), wg); err != nil {
					logging.Error(errors.Wrapf(err, "failed to run '%s' check", name))
				}
			}(name, check)
		}
	}()

	return onlyReadable(repoChan)
}
