package repository

import (
	"context"
	"sync"

	"github.com/google/go-github/github"
	"github.com/leonklingele/github-release-checker/logging"
	"github.com/pkg/errors"
)

func listStarredAll(
	activity *github.ActivityService,
	repoChan chanW,
	wg *sync.WaitGroup,
) error {
	firstPage := 1
	res, err := listStarred(activity, firstPage, repoChan)
	if err != nil {
		return errors.Wrap(err, "failed to get initial starred page")
	}
	for page := firstPage + 1; page <= res.LastPage; page++ {
		wg.Add(1)
		go func(page int) {
			defer func() {
				wg.Done()
				logging.Debug("done getting starred repos, page", page)
			}()
			logging.Debug("start getting starred repos, page", page)
			if _, err := listStarred(activity, page, repoChan); err != nil {
				logging.Error(errors.Wrapf(err, "failed to get starred repos, page %d", page))
			}
		}(page)
	}
	return nil
}

func listStarred(
	activity *github.ActivityService,
	page int,
	repoChan chanW,
) (*github.Response, error) {
	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}
	// TODO(leon): Pass in context
	starred, res, err := activity.ListStarred(context.TODO(), "", opt)
	if err != nil {
		return nil, err
	}
	for _, star := range starred {
		repoChan <- newRepository(star.Repository)
	}
	return res, nil
}
