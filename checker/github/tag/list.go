package tag

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/leonklingele/github-release-checker/checker/github/repository"
	"github.com/leonklingele/github-release-checker/logging"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

const (
	// TODO(leon): Make this configurable
	maxWorkers = 100
)

func newListWorker(repoChan repository.Chan, tagChan chanW, c *fasthttp.Client) {
	ts := time.Now().Unix()
	fp := gofeed.NewParser()
	for repo := range repoChan {
		htmlURL := repo.GetHTMLURL()
		atom := "tags.atom"
		tagsURL := fmt.Sprintf("%s/%s?t=%d", htmlURL, atom, ts)
		status, body, err := c.Get(nil, tagsURL)
		if status != fasthttp.StatusOK {
			logging.Errorf("bad HTTP status code of %s: %d", tagsURL, status)
			continue
		}
		if err != nil {
			logging.Error(errors.Wrapf(err, "failed to get URL %s", tagsURL))
			continue
		}
		feed, err := fp.Parse(bytes.NewReader(body))
		if err != nil {
			logging.Error(errors.Wrapf(err, "failed to parse response of %s", tagsURL))
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

		c := &fasthttp.Client{}
		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				newListWorker(repoChan, onlyWritable(tagChan), c)
			}()
		}
	}()

	return onlyReadable(tagChan)
}
