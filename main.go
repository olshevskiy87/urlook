package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/olshevskiy87/urlook/urlook"
	"mvdan.cc/xurls"
)

func parseURLs(source string) ([]string, error) {
	re, err := xurls.StrictMatchingScheme("https?://")
	if err != nil {
		return nil, fmt.Errorf("could not prepare search regex: %s", err)
	}
	result := re.FindAllString(source, -1)
	return result, nil
}

func main() {
	var args struct {
		Filename string `arg:"positional" help:"filename with links to check"`
	}
	arg.MustParse(&args)

	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get stdin stats: %v", err)
		os.Exit(1)
	}
	isPipe := (stdinStat.Mode() & os.ModeCharDevice) == 0
	if isPipe && args.Filename != "" {
		fmt.Fprintln(os.Stderr, "please specify a filename or pass text from standard input")
		os.Exit(1)
	}
	var inputText string
	if isPipe {
		stdinText, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read from stdin: %v", err)
			os.Exit(1)
		}
		inputText = string(stdinText)
	} else {
		if args.Filename == "" {
			fmt.Fprintf(os.Stderr, "specify a filename")
			os.Exit(1)
		}
		fileContent, err := ioutil.ReadFile(args.Filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read file \"%s\": %v", args.Filename, err)
			os.Exit(1)
		}
		inputText = string(fileContent)
	}
	urls, err := parseURLs(inputText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse urls: %v\n", err)
		os.Exit(1)
	}
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "no URLs found")
		os.Exit(1)
	}
	app := urlook.New(urls)
	/*
	 *err = app.SetTimeout(5)
	 *if err != nil {
	 *    fmt.Println(err)
	 *    os.Exit(1)
	 *}
	 */
	err = app.CheckAllURLs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("no issues found")
}
