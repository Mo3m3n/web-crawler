package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`This is a webcrawler client that sends crawl requests to the webcrawler server.

webcralwer <server-url> <url> [depth]

	server-url: the url of the webcrawler. Example: 'http://127.0.0.1:8080/'
	url: the starting url to crawl from. Example: 'https://example.com/foo'
  depth: the extent/level to which the webcrawler fetchs links. -1 means no limit.
`)
	}
	flag.Parse()
	address := flag.Arg(0)
	url := flag.Arg(1)
	depth := flag.Arg(2)
	if address == "" {
		fmt.Println("Error empty address")
		return
	}
	if url == "" {
		fmt.Println("Error empty url")
		return
	}
	// Create HTTP request
	var err error
	var req *http.Request
	req, err = http.NewRequest("GET", address, nil)
	if err != nil {
		fmt.Printf("Error unable to create http request: %s\n", err)
		return
	}
	q := req.URL.Query()
	q.Add("url", url)
	if depth != "" {
		if _, err = strconv.Atoi(depth); err != nil {
			fmt.Printf("Incorrect depth arg '%s': should be an integer\n", depth)
			return
		}
		q.Add("depth", depth)
	}
	req.URL.RawQuery = q.Encode()
	// Get HTTP Resonse
	var resp *http.Response
	var body []byte
	resp, err = http.Get(req.URL.String())
	if err != nil {
		fmt.Printf("Error unable to get http response: %s\n", err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error unable to read http response: %s\n", err)
		return
	}
	fmt.Print(string(body))
}
