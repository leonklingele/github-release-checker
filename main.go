package main

import (
	"flag"

	"github.com/leonklingele/github-release-checker/checker"
	logHandler "github.com/leonklingele/github-release-checker/checker/handlers/log"
	mailHandler "github.com/leonklingele/github-release-checker/checker/handlers/mail"
	"github.com/leonklingele/github-release-checker/config"
	"github.com/leonklingele/github-release-checker/logging"
	"github.com/leonklingele/github-release-checker/pathutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const (
	defaultConfigFilePath = "$HOME/.github-release-checker/config.toml"
)

func main() {
	configFilePath := flag.String("config", defaultConfigFilePath, "optional, path where to find the config file")
	enableDebug := flag.Bool("debug", false, "optional, whether to enable debug mode")
	flag.Parse()

	if *enableDebug {
		logging.SetDebug()
	}

	if err := boot(*configFilePath); err != nil {
		logging.Fatal(err)
	}
}

func boot(configFilePath string) error {
	cfp, err := annotateConfigFilePath(configFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to annotate config file path")
	}

	conf, err := config.Load(cfp)
	if err != nil {
		return err
	}

	if err := logHandler.Register(); err != nil {
		logging.Error(errors.Wrap(err, "failed to register log handler"))
	}
	if err := mailHandler.Register(conf.MailConfig); err != nil {
		logging.Error(errors.Wrap(err, "failed to register mail handler"))
	}

	c, err := checker.New(conf.CheckerConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create checker")
	}
	return c.Start()
}

func annotateConfigFilePath(p string) (string, error) {
	return pathutil.ReplaceHome(p)
}
