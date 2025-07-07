package main

import "fmt"

func main() {
	crawler := NewCrawler("https://scrape-me.dreamsofcode.io/", 5)

	// Start the concurrent crawler
	routes := crawler.StartConcurrent(10)

	// Print the discovered routes
	fmt.Println("Discovered routes:")
	for _, route := range routes {
		fmt.Println("-", route)
	}
}
