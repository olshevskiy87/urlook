package urlook

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/olshevskiy87/urlook/urlook/status"
)

// Bot is a main app object
type Bot struct {
	urls         []string
	clientHTTP   *http.Client
	workersCount int
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

// New returns new Bot object
func New(urls []string) *Bot {
	return &Bot{
		urls: urls,
		clientHTTP: &http.Client{
			Timeout: time.Duration(
				time.Duration(DefaultHTTPRequestTimeout) * time.Second,
			),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		workersCount: runtime.NumCPU() * 2, // TODO: set a proper value
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

// CheckAllURLs checks all given URLs
func (b *Bot) CheckAllURLs() error {
	var (
		urlsCount = len(b.urls)
		results   = make([]*Result, urlsCount)
	)
	fmt.Printf("URLs to check: %d\n", urlsCount)
	for i, url := range b.urls {
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
	for i, url := range b.urls {
		urlChan <- &URLChanItem{
			URL:   url,
			Index: i,
		}
	}
	close(urlChan)
	wg.Wait()

	fmt.Println()

	var errorsCount int
	for _, res := range results {
		if status.IsSuccess(res.Status.Code) {
			continue
		}
		errorsCount++
		fmt.Printf("%3d. %s\n", errorsCount, res)
	}
	if errorsCount > 0 {
		return fmt.Errorf("issues found: %d", errorsCount)
	}
	return nil
}

// checkURL checks one url and returns pointer to Result
func (b *Bot) checkURL(url string) (*Result, error) {
	if b.clientHTTP == nil {
		return nil, errors.New("http client is not defined")
	}
	// TODO: use Head on 405 (Method Not Allowed)
	resp, err := b.clientHTTP.Get(url)
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

	if status.IsRedirect(res.Status.Code) {
		locationURL, err := resp.Location()
		if err != nil {
			res.Message = fmt.Sprintf("could not retrieve Location: %s", err.Error())
		} else {
			res.Message = locationURL.String()
		}
	}
	return res, nil
}
