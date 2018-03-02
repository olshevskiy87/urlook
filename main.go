package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/olshevskiy87/urlook/urlook"
	"mvdan.cc/xurls"
)

func parseURLs(filepath string) ([]string, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file \"%s\": %v", filepath, err)
	}
	re, err := xurls.StrictMatchingScheme("https?://")
	if err != nil {
		return nil, fmt.Errorf("could not prepare search regex: %s", err)
	}
	result := re.FindAllString(string(file), -1)
	return result, nil
}

func main() {
	var args struct {
		Filename string `arg:"positional,required" help:"filename with links to check"`
	}
	arg.MustParse(&args)

	urls, err := parseURLs(args.Filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse file %s: %s\n", args.Filename, err)
		os.Exit(1)
	}
	if len(urls) == 0 {
		fmt.Fprintf(os.Stderr, "no URLs found in file %s\n", args.Filename)
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
	} else {
		fmt.Println("no issues found")
	}
}
