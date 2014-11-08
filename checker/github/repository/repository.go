package repository

import (
	"github.com/google/go-github/github"
)

type Repository struct {
	*github.Repository
	IsIgnored   bool
	IsImportant bool
}

func newRepository(repo *github.Repository) *Repository {
	return &Repository{
		Repository: repo,
	}
}
