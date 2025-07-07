package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func requestHandler(url string) (resp *http.Response, body []byte) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}
	body, _ = io.ReadAll(resp.Body)

	return resp, body
}

func htmlParser(body []byte) (*html.Node, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doc, nil
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

func getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isSameDomain(url, base string) bool {
	baseDomain := strings.Split(base, "/")[2]
	urlDomain := strings.Split(url, "/")[2]
	return baseDomain == urlDomain
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

func check404(doc *html.Node) bool {
	return true
}

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

		_, body := requestHandler(currentURL)
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
