package requests

import (
	"net/http"
)

type HeaderEntry struct {
	ParamKey   string
	ParamValue string
}

func SendRequest(method string, url string, headers []HeaderEntry) (*http.Response, error) {
	// create a new http client
	client := &http.Client{}

	// build the request with the url and method
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// add headers
	for _, header := range headers {
		req.Header.Add(header.ParamKey, header.ParamValue)
	}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
