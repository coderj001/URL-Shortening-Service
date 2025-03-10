package helpers

import (
	"crypto/rand"
	"strings"

	"github.com/coderj001/URL-shortener/config"
)

// EnforceHTTP ...
func EnforceHTTP(url string) string {
	// make every url https
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// RemoveDomainError ...
func RemoveDomainError(url string) bool {
	// basically this functions removes all the commonly found
	// prefixes from URL such as http, https, www
	// then checks of the remaining string is the DOMAIN itself
	if url == config.GetConfig().GetDomain() {
		return false
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == config.GetConfig().GetDomain() {
		return false
	}
	return true
}

func GenerateID(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var id string
	for _, byte := range b {
		// Use modulo to map random bytes to our charset
		idx := int(byte) % len(charset)
		id += string(charset[idx])
	}

	return id, nil
}
