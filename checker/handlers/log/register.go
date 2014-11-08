package log

import (
	"github.com/leonklingele/github-release-checker/checker/handlers"
)

const (
	handlerName = "log"
)

func Register() error {
	return handlers.Register(handlerName, func() (handlers.Handler, error) {
		return &logHandler{}, nil
	})
}
