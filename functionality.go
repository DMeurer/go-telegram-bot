package main

import (
	"net/http"
)

type headerEntry struct {
	key   string
	value string
}

func sendRequest(method string, url string, headers []headerEntry) (*http.Response, error) {
	// create a new http client
	client := &http.Client{}

	// build the request with the url and method
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// add headers
	for _, header := range headers {
		req.Header.Add(header.key, header.value)
	}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
