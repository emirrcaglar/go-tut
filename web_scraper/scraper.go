package main

import (
	"fmt"
	"net/http"
)

func requestHandler(url string) (resp *http.Response) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(err)
	}
	return resp
}

func main() {
	//... scraping logic

	resp := requestHandler("https://scrape-me.dreamsofcode.io/")

	fmt.Println(resp)
}
