package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mo3m3n/webcrawler/crawler"
	"github.com/mo3m3n/webcrawler/logger"
	"github.com/mo3m3n/webcrawler/site"
)

// webcrawler options
const addr = "127.0.0.1:8080"
const timeout = 300
const maxconn = 5
const ratelimit = 1
const maxdepth = -1
const loglevel logger.LogLevel = logger.Info

type options struct {
	address   string
	timeout   int
	maxconn   int
	ratelimit int
	loglevel  int
}

var args options

// connections limiter
var conn chan struct{}

func getParams() {
	addrPtr := flag.String("address", addr, "the TCP network address the webcrawler is going to listen to")
	timeoutPtr := flag.Int("timeout", timeout, "the number of seconds the webcrawler is going to wait for a crawl operation before interrupting it")
	maxconnPtr := flag.Int("maxconn", maxconn, "the maximum number of concurrent requests the webcrawler can accept")
	rateLimitPtr := flag.Int("ratelimit", ratelimit, "the maximum number of requests/second the webcrawler is allowed to send to a given website")
	logPtr := flag.Int("log", int(loglevel), "the webcrawler logging level: 1=error, 2=warning, 3=info, 4=debug")
	flag.Parse()
	args.address = *addrPtr
	args.timeout = *timeoutPtr
	args.maxconn = *maxconnPtr
	args.ratelimit = *rateLimitPtr
	args.loglevel = *logPtr
}

func requestHandler(writer http.ResponseWriter, request *http.Request) {
	var sitemap site.SiteMap
	var param []string
	var url string
	var depth = maxdepth
	var ok bool
	var err error
	log := logger.New(logger.LogLevel(args.loglevel), request.RemoteAddr)
	log.Infof("crawling request received")
	select {
	default:
		http.Error(writer, "Service busy, get back later", http.StatusServiceUnavailable)
		log.Warningf("request rejected: maxconn '%d' reached", args.maxconn)
		return
	case conn <- struct{}{}:
		defer func() { <-conn }()
		// Get params
		param, ok = request.URL.Query()["url"]
		if !ok {
			http.Error(writer, "url parameter required", http.StatusNotAcceptable)
			log.Errorf("url parameter required")
			return
		}
		url = param[0]
		param = request.URL.Query()["depth"]
		if len(param) != 0 {
			depth, err = strconv.Atoi(param[0])
			if err != nil {
				http.Error(writer, fmt.Sprintf("incorrect depth parameter %s", err), http.StatusNotAcceptable)
				log.Errorf("incorrect depth '%s' parameter %s", param[0], err)
				return
			}
		}
		// Crawl
		log.Infof("starting crawler for url '%s' and depth '%d'", url, depth)
		sitemap, err = crawler.Crawl(request.Context(), url, args.timeout, args.ratelimit, depth, log)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotAcceptable)
			log.Errorf("crawling request failed: %s", err)
		}
		bytes, err := sitemap.Marshal()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotAcceptable)
			log.Errorf("response marshalling failed: %s", err)
			return
		}
		// Write HTTP Response
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(bytes)
		if err != nil {
			log.Errorf("unable to write response back to client: %s", err)
		} else {
			log.Infof("crawling request for '%s' finished", url)
		}
	}
}

func main() {
	getParams()
	conn = make(chan struct{}, args.maxconn)
	fmt.Printf("Starting web crawler at '%s'\n", args.address)
	fmt.Printf("Using options: maxconn %d, ratelmit %d, timeout %d\n", args.maxconn, args.ratelimit, args.timeout)
	server := http.Server{
		Addr:         args.address,
		Handler:      http.TimeoutHandler(http.HandlerFunc(requestHandler), time.Duration(args.timeout)*time.Second, "Timeout!\n"),
		WriteTimeout: time.Duration(args.timeout) * time.Second * 2, // Should be higher than handler timeout to close connection properly
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: unable to start http listener: %s\n", err)
		os.Exit(1)
	}
}
