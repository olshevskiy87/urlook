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

var args struct {
	Filenames       []string `arg:"positional" help:"filenames with links to check"`
	FailOnDuplicate bool     `arg:"--fail-on-duplicate" help:"fail if there is a duplicate url"`
	RequestTimeout  int      `arg:"--timeout,-t" help:"request timeout in seconds"`
	WhiteList       []string `arg:"--white,-w,separate" help:"white list url (can be specified multiple times)"`
}

func init() {
	args.RequestTimeout = 10
	arg.MustParse(&args)
}

func parseURLs(source string) ([]string, error) {
	re, err := xurls.StrictMatchingScheme("https?://")
	if err != nil {
		return nil, fmt.Errorf("could not prepare search regex: %v", err)
	}
	return re.FindAllString(source, -1), nil
}

func getInputText() (string, error) {
	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("could not get stdin stats: %v", err)
	}
	isPipe := (stdinStat.Mode() & os.ModeCharDevice) == 0
	filenamesCnt := len(args.Filenames)
	if isPipe && filenamesCnt > 0 {
		return "", fmt.Errorf("please specify at least one filename or pass text from standard input")
	}
	if isPipe {
		stdinText, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("could not read from stdin: %v", err)
		}
		return string(stdinText), nil
	}
	if filenamesCnt == 0 {
		return "", fmt.Errorf("specify at least one filename")
	}
	var inputBuffer bytes.Buffer
	for _, fname := range args.Filenames {
		fileContent, err := ioutil.ReadFile(fname)
		if err != nil {
			return "", fmt.Errorf("could not read file \"%s\": %v", fname, err)
		}
		inputBuffer.Write(fileContent)
	}
	return inputBuffer.String(), nil
}

func main() {
	inputText, err := getInputText()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get input text: %v\n", err)
		os.Exit(1)
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
