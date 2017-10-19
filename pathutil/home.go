package pathutil

import (
	"os/user"
	"strings"

	"github.com/pkg/errors"
)

const (
	homeVar = "$HOME"
)

func ReplaceHome(s string) (string, error) {
	if !strings.Contains(s, homeVar) {
		return s, nil
	}
	u, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "failed to get current user")
	}
	return strings.Replace(s, homeVar, u.HomeDir, -1), nil
}
