package main

import "net/http"

// NewHTTPClient create http client
func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	return client
}
