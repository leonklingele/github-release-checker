package mail

import (
	"crypto/tls"
)

type Config struct {
	Enabled    bool `toml:"enabled"`
	NumWorkers int  `toml:"workers"`
	Insecure   bool `toml:"insecure"`
	TLSConfig  *tls.Config

	Host string `toml:"host"`
	Port int    `toml:"port"`
	User string `toml:"user"`
	Pswd string `toml:"pswd"`

	From        string   `toml:"from"`
	To          []string `toml:"to"`
	ImportantTo []string `toml:"important_to"`

	Subject string `toml:"subject"`
	Body    string `toml:"body"`
}
