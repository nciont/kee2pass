package kee2pass

import (
	"net/url"
	"netelevate.com/kee2pass/entities"
	"regexp"
	"strings"
	"unicode/utf8"
)

var re = regexp.MustCompile("[[:space:]/\\.]+")

// PrettifyURL extracts hostname if available
func PrettifyURL(rawurl string) string {
	if addr, e := url.Parse(rawurl); e == nil {
		if host := addr.Hostname(); host != "" {
			return host
		}
	}

	return rawurl
}

// GetItemName returns name for an item
func GetItemName(item *entities.XMLEntry) string {
	return normalizeName(extractRawName(item))
}

// normalizeName lowercases, replaces whitespace and slashes with dots
func normalizeName(name string) string {
	return strings.ToLower(re.ReplaceAllLiteralString(strings.TrimSpace(name), "."))
}

func extractRawName(item *entities.XMLEntry) string {
	var url = PrettifyURL(item.URL)

	// prefer username@host
	if url != "" && item.User != "" {
		return item.User + "@" + url
	}

	// then title
	if utf8.RuneCountInString(item.Title) > 2 {
		return item.Title
	}

	// then just host
	if url != "" {
		return url
	}

	// user, warn with question marks
	if item.User != "" {
		return item.User + "@" + "???"
	}

	return item.UUID
}
