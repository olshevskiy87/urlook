package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestParseURLs(t *testing.T) {
	tests := []struct {
		url           string
		shouldBeFound bool
	}{
		{"https://domain1.com", true},
		{"http://domain2.com", true},
		{"ftp://domain3.com", false},
	}
	urls, err := parseURLs("test urls: https://domain1.com http://domain2.com ftp://domain3.com")
	if err != nil {
		t.Errorf("could not parse urls: %v", err)
	}
	for _, test := range tests {
		found := false
		for _, url := range urls {
			if url == test.url {
				found = true
				break
			}
		}
		var shouldBeFoundStr string
		if !test.shouldBeFound {
			shouldBeFoundStr = "not"
		}
		if found != test.shouldBeFound {
			t.Errorf("url %s should %s be found", test.url, shouldBeFoundStr)
		}
	}
}

func TestGetInputText(t *testing.T) {
	expected := "test urls list"
	rr := []io.ReadCloser{
		ioutil.NopCloser(
			bytes.NewBufferString(expected),
		),
	}
	inputText, err := getInputText(rr)
	if err != nil {
		t.Errorf("could not get input text: %v", err)
	}
	if expected != inputText {
		t.Errorf("expected inputText \"%s\" but got \"%s\"", expected, inputText)
	}
}
