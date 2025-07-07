package main

import (
	"net/http"

	"golang.org/x/net/html"
)

// represents the result of a single crawl operation
type CrawlResult struct {
	Url        string
	StatusCode int
	Routes     []string
	Error      error
}

// holds the state and configuration for the web crawler
type Crawler struct {
	BaseUrl      string
	MaxDepth     int
	Visited      map[string]bool
	UniqueRoutes map[string]bool
	Client       *http.Client
}

// represents a parsed HTML page
type Page struct {
	Url      string
	Document *html.Node
	Links    []Link
}

// represents an extracted link from a page
type Link struct {
	Href       string
	Normalized string
	Text       string
	SameDomain string
}

// holds crawler configuration
type CrawlConfig struct {
	BaseUrl   string
	MaxDepth  int
	UserAgent string
	Timeout   int // seconds
}

// tracks statistics for a Url
type UrlStats struct {
	Url         string
	Visits      int
	LastVisited string
	StatusCode  int
}
