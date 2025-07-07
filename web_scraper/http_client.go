package main

import (
	"io"
	"net/http"
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

func isErrorResponse(resp *http.Response) bool {
	if resp == nil {
		return true
	}
	status := resp.StatusCode
	return status >= 400 && status < 600
}
