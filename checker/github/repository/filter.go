package repository

import (
	"github.com/leonklingele/github-release-checker/logging"
)

type filterFunc func(*Repository) bool

func MakeChanFilter(filter filterFunc) func(Chan) Chan {
	return func(in Chan) Chan {
		return filterChan(in, filter)
	}
}

func filterChan(in Chan, isOK filterFunc) Chan {
	out := make(chanRW)
	go func() {
		for r := range in {
			if isOK(r) {
				out <- r
			}
		}
		close(out)
	}()
	return onlyReadable(out)
}

func CleanupChan(in Chan) Chan {
	filterIgnored := MakeChanFilter(func(repo *Repository) bool {
		if repo.IsIgnored {
			logging.Debug("skipping ignored repo", repo.GetFullName())
		}
		return !repo.IsIgnored
	})
	checked := make(map[int64]struct{})
	dedup := MakeChanFilter(func(repo *Repository) bool {
		// No need to lock the map, we're not checking concurrently
		if _, ok := checked[repo.GetID()]; ok {
			logging.Debug("already checked repo", repo.GetFullName())
			return false
		}
		checked[repo.GetID()] = struct{}{}
		return true
	})
	return dedup(filterIgnored(in))
}
