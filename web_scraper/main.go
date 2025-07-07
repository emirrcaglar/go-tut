package main

func main() {
	crawler := NewCrawler("https://scrape-me.dreamsofcode.io/", 5)
	crawler.Start()
	crawler.PrintRoutes()
}
