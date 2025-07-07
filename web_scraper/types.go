package main

import (
	"net/http"
)

// holds the state and configuration for the web crawler
type Crawler struct {
	BaseUrl      string
	MaxDepth     int
	Visited      map[string]bool
	UniqueRoutes map[string]bool
	Client       *http.Client
}
