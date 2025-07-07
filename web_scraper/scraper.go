package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

func printNode(n *html.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%sType: %d", indent, n.Type)

	if n.Type == html.ElementNode {
		fmt.Printf(", Tag: %s", n.Data)
	} else if n.Type == html.TextNode {
		fmt.Printf(", Text: %q", strings.TrimSpace(n.Data))
	}

	if len(n.Attr) > 0 {
		fmt.Printf(", Attr: %v", n.Attr)
	}
	fmt.Println()

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		printNode(c, depth+1)
	}
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

func printAnchorTags(anchorTags []*html.Node) {
	for i, anchor := range anchorTags {
		fmt.Printf("Anchor #%d:\n", i+1)

		for _, attr := range anchor.Attr {
			if attr.Key == "href" {
				fmt.Printf("	Link: %s\n", attr.Val)
			}
		}

		for child := anchor.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				fmt.Printf("  Text: %q\n", strings.TrimSpace(child.Data))
			}
		}
		fmt.Println()
	}
}

func normalizeURL(href, base string) string {
	if strings.HasPrefix(href, "http") {
		return href // already absolute
	}
	if strings.HasPrefix(href, "/") {
		return base + href
	}
	return base + "/" + href
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

func crawl(url string, baseURL string, maxDepth int, visited map[string]bool) []string {
	if maxDepth <= 0 || visited[url] {
		return nil
	}
	visited[url] = true

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

	maxDepth := 4

	visited := make(map[string]bool)

	allRoutes := crawl(baseURL, baseURL, maxDepth, visited)

	fmt.Println("Discovered routes: ")
	for i, route := range allRoutes {
		fmt.Printf("%d. %s\n", i+1, route)
	}
}
