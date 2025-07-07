package main

import (
	"strings"

	"golang.org/x/net/html"
)

func crawl(url string, baseURL string, maxDepth int, visited *map[string]bool) []string {
	if maxDepth <= 0 || (*visited)[url] {
		return nil
	}
	(*visited)[url] = true

	_, body := requestHandler(url)
	if body == nil {
		return nil
	}

	doc, err := htmlParser(body)
	if err != nil {
		return nil
	}

	routes := extractRoutes(doc, baseURL)
	var allRoutes []string
	allRoutes = append(allRoutes, routes...)

	for _, route := range routes {
		if isSameDomain(route, baseURL) {
			deeperRoutes := crawl(route, baseURL, maxDepth-1, visited)
			allRoutes = append(allRoutes, deeperRoutes...)
		}
	}
	return allRoutes
}

func findElements(n *html.Node, tag string) []*html.Node {
	var nodes []*html.Node
	if n.Type == html.ElementNode && n.Data == tag {
		nodes = append(nodes, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, findElements(c, tag)...)
	}
	return nodes
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
