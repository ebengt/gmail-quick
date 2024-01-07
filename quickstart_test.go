package main

import (
	"os"
	"testing"
)

func Test_config(t *testing.T) {
	args := []string{"asd", "-infile", "afile", "-subject", "asubject", "areceiver"}
	os.Setenv("from", "sender")
	c := configuration(args, "from")
	if c.infile != "afile" {
		t.Errorf("config = %s; want afile", c.infile)
	}
	if c.subject != "asubject" {
		t.Errorf("config = %s; want asubject", c.subject)
	}
	if c.receiver != "areceiver" {
		t.Errorf("config = %s; want areceiver", c.receiver)
	}
	if c.from != "sender" {
		t.Errorf("config = %s; want sender", c.from)
	}
}

func Test_parse_request(t *testing.T) {
	args := "http://localhost:4000/?state=state-token&code=4/0AfJohXmVHC-A5HOGVI4MDvXUyLc4plBIkjLzcPHxfyuXxQPaikecCRsVWTSDZcE4sVXzdg&scope=https://www.googleapis.com/auth/gmail.readonly%20https://www.googleapis.com/auth/gmail.compose"
	code := getTokenFromWebParse(args)
	if code != "4/0AfJohXmVHC-A5HOGVI4MDvXUyLc4plBIkjLzcPHxfyuXxQPaikecCRsVWTSDZcE4sVXzdg" {
		t.Errorf("%s; want 4/0AfJohXmVHC-A5HOGVI4MDvXUyLc4plBIkjLzcPHxfyuXxQPaikecCRsVWTSDZcE4sVXzdg", code)
	}
}
