package lib

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func checkPassword(saved, input string) bool {
	if strings.HasPrefix(saved, "{bcrypt}") {
		savedPassword := strings.TrimPrefix(saved, "{bcrypt}")
		return bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(input)) == nil
	}

	return saved == input
}

func isAllowedHost(allowedHosts []string, origin string) bool {
	for _, host := range allowedHosts {
		if host == origin {
			return true
		}
	}
	return false
}

var errPrefixMismatch = errors.New("webdav: prefix mismatch")

func stripPrefix(urlPath string, prefix string) (string, error) { // /abc /abcd/
	if len(prefix) == 0 {
		return urlPath, nil
	}
	if r := strings.TrimPrefix(urlPath, prefix); len(r) < len(urlPath) {
		if len(r) == 0 {
			r = `/`
		}
		return r, nil
	}
	return urlPath, errPrefixMismatch
}
