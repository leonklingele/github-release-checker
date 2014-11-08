package tag

import (
	"fmt"
	"strings"
	"sync"

	"github.com/leonklingele/github-release-checker/checker/github/repository"
	"github.com/leonklingele/github-release-checker/logging"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

const (
	// TODO(leon): Make this configurable
	maxWorkers = 100
)

func newListWorker(repoChan repository.Chan, tagChan chanW) {
	for repo := range repoChan {
		htmlURL := repo.GetHTMLURL()
		tagsURL := htmlURL + "/tags.atom"
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(tagsURL)
		if err != nil {
			logging.Error(errors.Wrap(err, fmt.Sprintf("failed to parse URL %s", tagsURL)))
			continue
		}
		for _, item := range feed.Items {
			split := strings.Split(item.Link, "/")
			version := split[len(split)-1]
			tagChan <- newTag(repo, version)
		}
	}
}

func List(repoChan repository.Chan) Chan {
	tagChan := make(chanRW)

	go func() {
		wg := &sync.WaitGroup{}
		defer func() {
			wg.Wait()
			close(tagChan)
			logging.Debug("done listing tags")
		}()

		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				newListWorker(repoChan, onlyWritable(tagChan))
			}()
		}
	}()

	return onlyReadable(tagChan)
}
