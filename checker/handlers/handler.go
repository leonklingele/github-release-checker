package handlers

import (
	"github.com/leonklingele/github-release-checker/checker/github/tag"
)

type Handler interface {
	Handle(tag.Chan)
}
