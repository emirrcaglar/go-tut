package main

import (
	"fmt"
	"io"
	"net/http"

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

func htmlParser(resp *http.Response) (doc *html.Node, err error) {
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return
	}

	return doc, nil
}

func main() {
	//... scraping logic

	_, body := requestHandler("https://scrape-me.dreamsofcode.io/")

	fmt.Println(string(body))
}
