package main

import (
	"bytes"
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
		return nil, fmt.Errorf("could not prepare search regex: %v", err)
	}
	return re.FindAllString(source, -1), nil
}

func main() {
	var args struct {
		Filenames       []string `arg:"positional" help:"filenames with links to check"`
		FailOnDuplicate bool     `arg:"--fail-on-duplicate" help:"fail if there is a duplicate url"`
		RequestTimeout  int      `arg:"--timeout,-t" help:"request timeout in seconds"`
		WhiteList       []string `arg:"--white,-w,separate" help:"white list url (can be specified multiple times)"`
	}
	args.RequestTimeout = 10
	arg.MustParse(&args)

	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get stdin stats: %v\n", err)
		os.Exit(1)
	}
	isPipe := (stdinStat.Mode() & os.ModeCharDevice) == 0
	filenamesCnt := len(args.Filenames)
	if isPipe && filenamesCnt > 0 {
		fmt.Fprintln(os.Stderr, "please specify at least one filename or pass text from standard input")
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
		if filenamesCnt == 0 {
			fmt.Fprintln(os.Stderr, "specify at least one filename")
			os.Exit(1)
		}
		var inputBuffer bytes.Buffer
		for _, fname := range args.Filenames {
			fileContent, err := ioutil.ReadFile(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not read file \"%s\": %v", fname, err)
				os.Exit(1)
			}
			inputBuffer.Write(fileContent)
		}
		inputText = inputBuffer.String()
	}
	urls, err := parseURLs(inputText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse urls: %v\n", err)
		os.Exit(1)
	}
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "no URLs found")
		os.Exit(0)
	}
	app := urlook.New(urls)
	app.SetIsFailOnDuplicates(args.FailOnDuplicate)
	if err := app.SetTimeout(args.RequestTimeout); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.SetWhiteList(args.WhiteList)
	if err := app.CheckAllURLs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("no issues found")
}
