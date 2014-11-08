package config

import (
	"github.com/leonklingele/github-release-checker/checker"
	"github.com/leonklingele/github-release-checker/checker/handlers/mail"
)

type Config struct {
	CheckerConfig *checker.Config `toml:"checker"`
	MailConfig    *mail.Config    `toml:"mail"`
}
