package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/alexflint/go-arg"
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
		t.Fatalf("could not parse urls: %v", err)
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
			t.Fatalf("url %s should %s be found", test.url, shouldBeFoundStr)
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
		t.Fatalf("could not get input text: %v", err)
	}
	if expected != inputText {
		t.Fatalf("expected inputText \"%s\" but got \"%s\"", expected, inputText)
	}
}

func TestGetReaders(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	tmpFile, err := ioutil.TempFile(os.TempDir(), "urlook_test")
	if err != nil {
		t.Fatalf("could not make temp file: %v", err)
	}
	defer func(fname string) {
		err := os.Remove(fname)
		if err != nil {
			t.Logf("could not remove temp file %s: %v", fname, err)
		}
	}(tmpFile.Name())
	t.Run(
		"one input file only",
		func(t *testing.T) {
			os.Args = []string{os.Args[0], tmpFile.Name()}
			var args argsType
			arg.MustParse(&args)

			readers, err := getReaders(args)
			if err != nil {
				t.Fatalf("could not get readers: %v", err)
			}
			defer closeReaders(readers)

			readersCnt := len(readers)
			if readersCnt != 1 {
				t.Fatalf("expected readers count 1 but got %d", readersCnt)
			}
		},
	)
	t.Run(
		"stdin only",
		func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("could not open pipe: %v", err)
			}
			os.Stdin = r
			w.WriteString("test input string\n")
			defer w.Close()

			var args argsType
			arg.MustParse(&args)

			readers, err := getReaders(args)
			if err != nil {
				t.Fatalf("could not get readers: %v", err)
			}
			defer closeReaders(readers)

			readersCnt := len(readers)
			if readersCnt != 1 {
				t.Fatalf("expected readers count 1 but got %d", readersCnt)
			}
		},
	)
	t.Run(
		"stdin and one input file",
		func(t *testing.T) {
			os.Args = []string{os.Args[0], "someFile.md"}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("could not open pipe: %v", err)
			}
			os.Stdin = r
			w.WriteString("test input string\n")
			defer w.Close()

			var args argsType
			arg.MustParse(&args)

			_, err = getReaders(args)
			if err == nil {
				t.Fatalf("expected getReaders returns error but got nil")
			}
		},
	)
	t.Run(
		"neither stdin nor input file",
		func(t *testing.T) {
			os.Args = []string{os.Args[0]}

			var args argsType
			arg.MustParse(&args)

			_, err = getReaders(args)
			if err == nil {
				t.Fatalf("expected getReaders returns error but got nil")
			}
		},
	)
}
