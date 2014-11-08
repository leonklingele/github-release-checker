package stringutil

import (
	"strings"
)

func SliceContainsCaseInsensitive(sl []string, str string) bool {
	for _, s := range sl {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}
