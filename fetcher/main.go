package fetcher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Fetcher interface {
	// Fetch returns a slice of URLs found
	// on the page of the input URL
	Fetch(url string) (urls []string, err error)
}

type fetcher struct {
	ratelimiter <-chan time.Time
	timeout     int
}

func New(ratelimiter <-chan time.Time, timeout int) Fetcher {
	return fetcher{
		ratelimiter: ratelimiter,
		timeout:     timeout,
	}
}

func (f fetcher) Fetch(url string) (urls []string, err error) {
	<-f.ratelimiter
	// Make HTTP request
	var response *http.Response
	http.DefaultClient.Timeout = time.Duration(f.timeout) * time.Second
	response, err = http.Get(url)
	if err != nil {
		err = fmt.Errorf("HTTP request failed: %s ", err)
		return
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		err = fmt.Errorf("loading HTTP response body failed: %s ", err)
		return
	}

	// Find all links and process them
	var link string
	var hidden bool
	document.Find("a").Each(
		func(index int, element *goquery.Selection) {
			link, _ = element.Attr("href")
			_, hidden = element.Attr("hidden")
			if url != "" && !hidden {
				urls = append(urls, link)
			}
		},
	)
	return
}
