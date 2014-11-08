package mail

import (
	"strings"

	"github.com/leonklingele/github-release-checker/checker/github/tag"
	"github.com/leonklingele/github-release-checker/logging"
	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

type mailHandler struct {
	config *Config
}

func newWorker(d *gomail.Dialer, mailq <-chan *gomail.Message) error {
	s, err := d.Dial()
	if err != nil {
		return errors.Wrap(err, "failed to connect to mail server")
	}
	defer func() {
		if err := s.Close(); err != nil {
			logging.Error(errors.Wrap(err, "failed to close connection to mail server"))
		}
	}()
	for m := range mailq {
		if err := gomail.Send(s, m); err != nil {
			logging.Error(errors.Wrap(err, "failed to send mail"))
		}
	}
	return nil
}

func (mh *mailHandler) Handle(tagChan tag.Chan) {
	c := mh.config
	mailq := make(chan *gomail.Message, c.NumWorkers)

	d := gomail.NewDialer(c.Host, c.Port, c.User, c.Pswd)
	d.TLSConfig = c.TLSConfig
	for i := 0; i < c.NumWorkers; i++ {
		go func() {
			if err := newWorker(d, mailq); err != nil {
				logging.Error(errors.Wrap(err, "failed to spawn worker"))
			}
		}()
	}

	for tag := range tagChan {
		name := tag.Repository.GetName()
		fullName := tag.Repository.GetFullName()
		url := tag.Repository.GetHTMLURL()
		version := tag.Version

		to := c.To
		if tag.Repository.IsImportant {
			to = append(to, c.ImportantTo...)
		}
		annotate := func(s string) string {
			s = strings.Replace(s, "$name", name, -1)
			s = strings.Replace(s, "$fullName", fullName, -1)
			s = strings.Replace(s, "$url", url, -1)
			s = strings.Replace(s, "$version", version, -1)
			return s
		}
		subject := annotate(c.Subject)
		body := annotate(c.Body)

		for _, r := range to {
			m := gomail.NewMessage()
			m.SetHeader("From", c.From)
			m.SetAddressHeader("To", r, "")
			m.SetHeader("Subject", subject)
			m.SetBody("text/plain", body)
			mailq <- m
		}
	}
	close(mailq)
}
