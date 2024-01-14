package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	var port string
	flag.StringVar(&port, "port", "12345", "help message for flagname")
	c := make(chan string, 1)
	go serve_HTTP(c, "localhost:"+port)
	config.RedirectURL = "http://localhost:" + port
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)
	authCode := <-c
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

type handle struct {
	channel chan string
}

func serve_HTTP(c chan string, redirect_address string) {
	h := handle{channel: c}
	http.HandleFunc("/", h.code_from_query)
	log.Fatal(http.ListenAndServe(redirect_address, nil))
}

func (h *handle) code_from_query(_ http.ResponseWriter, r *http.Request) {
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatalf("Unable to parse query: %v", err)
	}
	if len(q["code"]) > 0 {
		h.channel <- q["code"][0]
	}
}
