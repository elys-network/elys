package main

import (
	"net"
	"net/http"
	"strings"
)

func isServiceRunning(url string) bool {
	// Remove the "tcp://" prefix if present
	if strings.HasPrefix(url, "tcp://") {
		url = strings.TrimPrefix(url, "tcp://")
	}

	// Attempt to make a TCP connection
	conn, err := net.Dial("tcp", url)
	if err == nil {
		conn.Close()
		return true
	}

	// If TCP connection fails, attempt an HTTP GET request
	resp, err := http.Get("http://" + url)
	if err == nil {
		resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}

	return false
}
