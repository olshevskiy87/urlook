package urlook

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/olshevskiy87/urlook/urlook/status"
)

// Bot is a main app object
type Bot struct {
	urls               map[string]int
	clientHTTP         *http.Client
	workersCount       int
	isFailOnDuplicates bool
	whiteList          []string
}

// URLChanItem contains an url and it's index
// in results array
type URLChanItem struct {
	URL   string
	Index int
}

// DefaultHTTPRequestTimeout is the max interval
// for http request in seconds
const DefaultHTTPRequestTimeout = 10

var userAgent = fmt.Sprintf(
	"%s_%s:urlook:v0.1 (by /u/olshevskiy87)",
	runtime.GOOS, runtime.GOARCH,
)

// New returns new Bot object
func New(urls []string) *Bot {
	var urlsMap = make(map[string]int, len(urls))
	for _, url := range urls {
		if _, ok := urlsMap[url]; !ok {
			urlsMap[url] = 1
			continue
		}
		urlsMap[url]++
	}
	return &Bot{
		urls: urlsMap,
		clientHTTP: &http.Client{
			Timeout: time.Duration(
				time.Duration(DefaultHTTPRequestTimeout) * time.Second,
			),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		workersCount:       runtime.NumCPU() * 2, // TODO: set a proper value
		isFailOnDuplicates: true,
	}
}

// SetTimeout sets the new http request timeout value
func (b *Bot) SetTimeout(newTimeout int) error {
	// timeout must be greater than 0
	if newTimeout < 0 {
		return fmt.Errorf("invaild timeout value: %d", newTimeout)
	}
	b.clientHTTP.Timeout = time.Duration(
		time.Duration(newTimeout) * time.Second,
	)
	return nil
}

// SetWorkersCount sets the new workers count value
func (b *Bot) SetWorkersCount(newWorkersCount int) {
	// there must be at least one worker
	if newWorkersCount < 1 {
		return
	}
	b.workersCount = newWorkersCount
}

// SetIsFailOnDuplicates defines if duplicate urls will be ignored
func (b *Bot) SetIsFailOnDuplicates(isFail bool) {
	b.isFailOnDuplicates = isFail
}

// SetWhiteList sets the white list urls
func (b *Bot) SetWhiteList(wl []string) {
	var newWhiteList = make([]string, 0, len(wl))
	for _, w := range wl {
		var trimmedW = strings.TrimSpace(w)
		if trimmedW == "" {
			continue
		}
		newWhiteList = append(newWhiteList, trimmedW)
	}
	b.whiteList = newWhiteList
}

// CheckAllURLs checks all given URLs
func (b *Bot) CheckAllURLs() error {
	urls := make([]string, 0, len(b.urls))
	duplicateUrls := make(map[string]int, len(b.urls))
	whiteListUrls := make([]string, 0, len(b.urls))
	for url, cnt := range b.urls {
		if cnt > 1 {
			duplicateUrls[url] = cnt
		}
		if b.checkURLInWhiteList(url) {
			whiteListUrls = append(whiteListUrls, url)
			continue
		}
		urls = append(urls, url)
	}
	var (
		urlsCount = len(urls)
		results   = make([]*Result, urlsCount)
	)
	fmt.Printf("URLs to check: %d\n", urlsCount)
	for i, url := range urls {
		fmt.Printf("%3d. %s\n", i+1, url)
	}

	var (
		wg      sync.WaitGroup
		urlChan = make(chan *URLChanItem)
	)
	wg.Add(b.workersCount)

	for i := 0; i < b.workersCount; i++ {
		go func() {
			defer wg.Done()
			for urlChanItem := range urlChan {
				res, err := b.checkURL(urlChanItem.URL)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				fmt.Print(res.Status.GetSign())
				results[urlChanItem.Index] = res
			}
		}()
	}
	for i, url := range urls {
		urlChan <- &URLChanItem{
			URL:   url,
			Index: i,
		}
	}
	close(urlChan)
	wg.Wait()

	var whiteListNum = len(whiteListUrls)
	if whiteListNum > 0 {
		fmt.Printf("\nWhite listed URLs (%d):\n", whiteListNum)
		for _, url := range whiteListUrls {
			fmt.Printf(" - %s\n", url)
		}
	}

	fmt.Println()
	var errorsCount int
	for _, res := range results {
		if res.Status.IsSuccess() {
			continue
		}
		errorsCount++
		fmt.Printf("%3d. %s\n", errorsCount, res)
	}
	if errorsCount > 0 {
		return fmt.Errorf("issues found: %d", errorsCount)
	}
	var duplicatesNum = len(duplicateUrls)
	if b.isFailOnDuplicates && duplicatesNum > 0 {
		fmt.Println("Duplicates:")
		for url, cnt := range duplicateUrls {
			fmt.Printf(" - %s (%d)\n", url, cnt)
		}
		return fmt.Errorf("duplicates found: %d", duplicatesNum)
	}
	return nil
}

// checkURL checks one url and returns pointer to Result
func (b *Bot) checkURL(url string) (*Result, error) {
	if b.clientHTTP == nil {
		return nil, errors.New("http client is not defined")
	}
	// TODO: use Head on 405 (Method Not Allowed)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not prepare new http request: %v", err)
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := b.clientHTTP.Do(req)
	if err != nil {
		return &Result{
			URL:     url,
			Status:  status.New(0),
			Message: err.Error(),
		}, nil
	}
	defer func(resp *http.Response) {
		err := resp.Body.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "could not close response body")
		}
	}(resp)

	res := &Result{
		URL:    url,
		Status: status.New(resp.StatusCode),
	}

	if res.Status.IsRedirect() {
		locationURL, err := resp.Location()
		if err != nil {
			res.Message = fmt.Sprintf("could not retrieve Location: %s", err.Error())
		} else {
			res.Message = locationURL.String()
		}
	}
	return res, nil
}

// checkURLInWhiteList checks if the given url is in white-list
func (b *Bot) checkURLInWhiteList(url string) bool {
	for _, whiteListURL := range b.whiteList {
		if strings.Contains(url, whiteListURL) {
			return true
		}
	}
	return false
}
