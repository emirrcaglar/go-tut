package main

import (
	"fmt"
)

func main() {
	baseURL := "https://scrape-me.dreamsofcode.io/"
	maxDepth := 5

	visited := make(map[string]bool)
	uniqueRoutes := make(map[string]bool)

	var crawlRecursive func(string, int)
	crawlRecursive = func(currentURL string, depth int) {
		if depth > maxDepth || visited[currentURL] {
			return
		}
		visited[currentURL] = true

		resp, body := requestHandler(currentURL)

		if isErrorResponse(resp) {
			return
		}

		doc, _ := htmlParser(body)

		anchors := findElements(doc, "a")
		for _, anchor := range anchors {
			href := getAttribute(anchor, "href")
			normalized := normalizeURL(href, baseURL)

			if normalized != "" && isSameDomain(normalized, baseURL) {
				uniqueRoutes[normalized] = true
				crawlRecursive(normalized, depth+1)
			}
		}
	}

	crawlRecursive(baseURL, 0)

	fmt.Println("Unique routes:")
	for route := range uniqueRoutes {
		fmt.Println("-", route)
	}
}
