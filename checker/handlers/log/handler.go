package log

import (
	"github.com/leonklingele/github-release-checker/checker/github/tag"
	"github.com/leonklingele/github-release-checker/logging"
)

type logHandler struct {
}

func (lh *logHandler) Handle(tagChan tag.Chan) {
	for tag := range tagChan {
		name := tag.Repository.GetFullName()
		version := tag.Version
		msg := "found new tag"
		if tag.Repository.IsImportant {
			msg = "found new important tag"
		}
		logging.Info(msg, name, version)
	}
}
