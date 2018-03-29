package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/olshevskiy87/urlook/urlook"
	"mvdan.cc/xurls"
)

type argsType struct {
	Filenames       []string `arg:"positional" help:"filenames with links to check"`
	FailOnDuplicate bool     `arg:"--fail-on-duplicate" help:"fail if there is a duplicate url"`
	RequestTimeout  int      `arg:"--timeout,-t" help:"request timeout in seconds"`
	WhiteList       []string `arg:"--white,-w,separate" help:"white list url (can be specified multiple times)"`
}

func parseURLs(source string) ([]string, error) {
	re, err := xurls.StrictMatchingScheme("https?://")
	if err != nil {
		return nil, fmt.Errorf("could not prepare search regex: %v", err)
	}
	return re.FindAllString(source, -1), nil
}

func getInputText(rr []io.ReadCloser) (string, error) {
	var inputBuffer bytes.Buffer
	for _, r := range rr {
		content, err := ioutil.ReadAll(r)
		if err != nil {
			return "", fmt.Errorf("could not read content: %v", err)
		}
		inputBuffer.Write(content)
	}
	return inputBuffer.String(), nil
}

func closeReaders(rr []io.ReadCloser) {
	for _, r := range rr {
		err := r.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not close reader\n")
		}
	}
}

func getReaders(args argsType) ([]io.ReadCloser, error) {
	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not get stdin stats: %v", err)
	}
	isPipe := (stdinStat.Mode() & os.ModeCharDevice) == 0
	filenamesCnt := len(args.Filenames)
	if isPipe && filenamesCnt > 0 {
		return nil, fmt.Errorf("please specify at least one filename or pass text from standard input")
	}
	if isPipe {
		return []io.ReadCloser{os.Stdin}, nil
	}
	var readers []io.ReadCloser
	if filenamesCnt == 0 {
		return nil, fmt.Errorf("specify at least one filename")
	}
	for _, fname := range args.Filenames {
		file, err := os.Open(fname)
		if err != nil {
			return nil, fmt.Errorf("could not open file %s: %v", fname, err)
		}
		readers = append(readers, file)
	}
	return readers, nil
}

func main() {
	var args argsType
	args.RequestTimeout = 10
	arg.MustParse(&args)

	readers, err := getReaders(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get readers: %v\n", err)
		os.Exit(1)
	}
	defer closeReaders(readers)
	inputText, err := getInputText(readers)
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
		fmt.Println("no URLs found")
		os.Exit(0)
	}
	app := urlook.New(urls)
	app.SetIsFailOnDuplicates(args.FailOnDuplicate)
	if err := app.SetTimeout(args.RequestTimeout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app.SetWhiteList(args.WhiteList)
	if err := app.CheckAllURLs(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("no issues found")
}
