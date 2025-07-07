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

func main() {
	baseURL := "https://scrape-me.dreamsofcode.io/"
	_, body := requestHandler(baseURL)

	doc, err := htmlParser(body)
	if err != nil {
		fmt.Println("Error parsing html: ", err)
		return
	}

	routes := extractRoutes(doc, baseURL)
	fmt.Println("Discovered routes: ")
	for _, route := range routes {
		fmt.Println("-", route)
	}
}
