package main

import (
	"net/http"
	"testing"
	"time"
)

const (
	expected = "4/0AfJohXmVHC-A5HOGVI4MDvXUyLc4plBIkjLzcPHxfyuXxQPaikecCRsVWTSDZcE4sVXzdg"
)

func Test_serve(t *testing.T) {
	c := make(chan string, 1)

	go serve_HTTP(c, "localhost:12345")
	_, err := http.Get("http://localhost:12345/?state=state-token&code=" + expected + "&scope=https://www.googleapis.com/auth/gmail.readonly%20https://www.googleapis.com/auth/gmail.compose")
	if err != nil {
		t.Errorf("Get = %v", err)
	}
	select {
	case code := <-c:
		if code != expected {
			t.Errorf("%s; want %s", code, expected)
		}
	case <-time.After(time.Second):
		t.Errorf("timeout")
	}
}
