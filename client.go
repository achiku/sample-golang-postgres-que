package main

import (
	"net/http"
	"time"
)

// NewHTTPClient create http client
func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{
		Timeout:   time.Duration(10) * time.Second,
		Transport: tr,
	}
	return client
}
