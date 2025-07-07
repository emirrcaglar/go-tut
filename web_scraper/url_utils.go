package main

import (
	"net/url"
	"strings"
)

func normalizeURL(href, base string) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}

	if !u.IsAbs() {
		baseU, _ := url.Parse(base)
		u = baseU.ResolveReference(u)
	}

	u.Fragment = ""
	u.RawQuery = ""
	return u.String()
}

func isSameDomain(url, base string) bool {
	baseDomain := strings.Split(base, "/")[2]
	urlDomain := strings.Split(url, "/")[2]
	return baseDomain == urlDomain
}
