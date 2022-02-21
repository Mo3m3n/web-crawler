package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
)

type options struct {
	depth    int
	insecure bool
	username string
	pass     string
}

var args options

func getParams() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`This is a webcrawler client that sends crawl requests to the webcrawler server.

webcralwer [options] <server-url> <url>

  server-url: the url of the webcrawler. Example: 'http://127.0.0.1:8080/'
  url: the starting url to crawl from. Example: 'https://example.com/foo'

  options:
    -depth
          the extent/level to which the webcrawler fetchs links. -1 means no limit.
    -insecure
          ignore server certificate verification when connecting over TLS
    -pass string
          password to be used for basic http authentication
    -username string
          username to be used for basic http authentication
`)
	}
	depthPtr := flag.Int("depth", -1, "the extent/level to which the webcrawler fetchs links. -1 means no limit")
	insecurePtr := flag.Bool("insecure", false, "ignore server certificate verification when connecting over TLS")
	usernamePtr := flag.String("username", "", "username to be used for basic http authentication")
	passPtr := flag.String("pass", "", "password to be used for basic http authentication")
	flag.Parse()
	args.depth = *depthPtr
	args.insecure = *insecurePtr
	args.username = *usernamePtr
	args.pass = *passPtr
}

func main() {
	getParams()
	address := flag.Arg(0)
	url := flag.Arg(1)
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
	if args.depth != -1 {
		q.Add("depth", strconv.Itoa(args.depth))
	}
	req.URL.RawQuery = q.Encode()
	if args.username != "" {
		req.SetBasicAuth(args.username, args.pass)
	}
	// Get HTTP Resonse
	var resp *http.Response
	var body []byte
	client := &http.Client{}
	if args.insecure {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		client.Transport = &http.Transport{
			DialTLS: func(network, addr string) (net.Conn, error) {
				conn, err := tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		}
	}
	resp, err = client.Do(req)
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
