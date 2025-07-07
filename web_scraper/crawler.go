package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func (c *Crawler) extractLinks(doc *html.Node, BaseUrl string) []string {
	anchors := findElements(doc, "a")
	uniqueRoutes := make(map[string]bool)
	var routes []string

	for _, anchor := range anchors {
		href := getAttribute(anchor, "href")
		if href == "" {
			continue
		}

		normalized := normalizeURL(href, BaseUrl)

		if !strings.HasPrefix(normalized, "http") || uniqueRoutes[normalized] {
			continue
		}

		if isSameDomain(normalized, BaseUrl) {
			uniqueRoutes[normalized] = true
			routes = append(routes, normalized)
		}

	}
	return routes
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

// Start begins the crawling process
func (c *Crawler) Start() []string {
	c.crawlRecursive(c.BaseUrl, 0)

	// Convert map to slice for return
	routes := make([]string, 0, len(c.UniqueRoutes))
	for route := range c.UniqueRoutes {
		routes = append(routes, route)
	}

	return routes
}

func (c *Crawler) worker(id int, jobs <-chan string, results chan<- []string) {
	for url := range jobs {
		fmt.Printf("Worker %d processing: %s\n", id, url)

		resp, body := requestHandler(url)
		if isErrorResponse(resp) {
			results <- []string{} // Send empty slice on error
			continue
		}

		doc, err := htmlParser(body)
		if err != nil {
			results <- []string{} // Send empty slice on error
			continue
		}

		routes := c.extractLinks(doc, c.BaseUrl)
		results <- routes
	}
}

// StartConcurrent begins a concurrent crawling process
func (c *Crawler) StartConcurrent(maxWorkers int) []string {
	jobs := make(chan string, maxWorkers)
	results := make(chan []string, maxWorkers)

	// Start workers
	for i := range maxWorkers {
		go c.worker(i+1, jobs, results)
	}

	// Keep track of all URLs (queued or processed) to avoid duplicates
	seen := make(map[string]bool)

	// Start with the base URL
	jobs <- c.BaseUrl
	seen[c.BaseUrl] = true
	c.UniqueRoutes[c.BaseUrl] = true
	activeJobs := 1

	for activeJobs > 0 {
		foundRoutes := <-results
		activeJobs--

		for _, route := range foundRoutes {
			if !seen[route] {
				seen[route] = true
				c.UniqueRoutes[route] = true
				activeJobs++
				jobs <- route
			}
		}
	}

	close(jobs) // All jobs are done, close channel to terminate workers

	// Convert map to slice for return
	finalRoutes := make([]string, 0, len(c.UniqueRoutes))
	for route := range c.UniqueRoutes {
		finalRoutes = append(finalRoutes, route)
	}

	return finalRoutes
}

// crawlRecursive performs recursive crawling
func (c *Crawler) crawlRecursive(currentURL string, depth int) {
	if depth > c.MaxDepth || c.Visited[currentURL] {
		return
	}

	c.Visited[currentURL] = true

	resp, body := requestHandler(currentURL)
	if isErrorResponse(resp) {
		return
	}

	doc, _ := htmlParser(body)
	anchors := findElements(doc, "a")

	for _, anchor := range anchors {
		href := getAttribute(anchor, "href")
		normalized := normalizeURL(href, c.BaseUrl)

		if normalized != "" && isSameDomain(normalized, c.BaseUrl) {
			c.UniqueRoutes[normalized] = true
			c.crawlRecursive(normalized, depth+1)
		}
	}
}

// PrintRoutes prints all discovered routes
func (c *Crawler) PrintRoutes() {
	fmt.Println("Unique routes:")
	for route := range c.UniqueRoutes {
		fmt.Println("-", route)
	}
}

// NewCrawler creates a new crawler instance
func NewCrawler(BaseUrl string, maxDepth int) *Crawler {
	return &Crawler{
		BaseUrl:      BaseUrl,
		MaxDepth:     maxDepth,
		Visited:      make(map[string]bool),
		UniqueRoutes: make(map[string]bool),
	}
}
