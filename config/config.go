package config

import (
	"github.com/BurntSushi/toml"
	"github.com/leonklingele/github-release-checker/checker"
	"github.com/leonklingele/github-release-checker/checker/handlers/mail"
	"github.com/pkg/errors"
)

type Config struct {
	CheckerConfig *checker.Config `toml:"checker"`
	MailConfig    *mail.Config    `toml:"mail"`
}

func Load(cfp string) (*Config, error) {
	var c Config
	if _, err := toml.DecodeFile(cfp, &c); err != nil {
		return nil, errors.Wrap(err, "failed to load or parse config file")
	}

	// TODO(leon): Set defaults, unelegantly
	if c.CheckerConfig.Workers <= 0 {
		c.CheckerConfig.Workers = 10
	}

	return &c, nil
}
