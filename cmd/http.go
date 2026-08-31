package main

import (
	"net/http"
	"time"
)

// http client for reuse
var httpClient = &http.Client{
	Timeout: 10 * time.Second, // Set a 10-second timeout for HTTP requests
}
