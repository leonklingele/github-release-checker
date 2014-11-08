package tag

import (
	"github.com/leonklingele/github-release-checker/checker/github/repository"
)

type Tag struct {
	Repository *repository.Repository
	Version    string
}

func newTag(repo *repository.Repository, version string) *Tag {
	return &Tag{
		Repository: repo,
		Version:    version,
	}
}
