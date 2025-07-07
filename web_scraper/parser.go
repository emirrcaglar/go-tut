package main

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func htmlParser(body []byte) (*html.Node, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func extractRoutes(doc *html.Node, baseURL string) []string {
	anchors := findElements(doc, "a")
	uniqueRoutes := make(map[string]bool)
	var routes []string
	for _, anchor := range anchors {
		href := getAttribute(anchor, "href")
		if href == "" {
			continue
		}
		normalized := normalizeURL(href, baseURL)
		if !strings.HasPrefix(normalized, "http") || uniqueRoutes[normalized] {
			continue
		}
		if isSameDomain(normalized, baseURL) {
			uniqueRoutes[normalized] = true
			routes = append(routes, normalized)
		}
	}
	return routes
}
