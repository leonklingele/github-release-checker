package github

import (
	"github.com/google/go-github/github"
)

type Github struct {
	config *Config
	Client *github.Client
}

func New(c *Config) *Github {
	bat := github.BasicAuthTransport{
		Username: c.User,
		Password: c.Token,
	}
	return &Github{
		config: c,
		Client: github.NewClient(bat.Client()),
	}
}
