package main

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/mo3m3n/webcrawler/fetcher"
	"github.com/mo3m3n/webcrawler/logger"
	"github.com/mo3m3n/webcrawler/ratelimiter"
	"github.com/mo3m3n/webcrawler/site"
)

func parseURL(p *url.URL, urlString string) (*url.URL, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	// handle relative url
	if u.Host == "" && p != nil {
		u.Host = p.Host
	}
	if u.Scheme == "" && p != nil {
		u.Scheme = p.Scheme
	}
	if u.Path[0] == byte('.') {
		u.Path = path.Join(p.Path, u.Path)
	}
	return u, nil
}

// crawls a website starting from a rool url and using a Breadth First Search algorithm via a local queue
func crawl(ctx context.Context, rootURL string, ratelimit, maxDepth int, log logger.Logger) (site.SiteMap, error) {
	var queue []site.URLNode
	var parent, child site.URLNode
	var result []string
	// Create root URLNode
	u, err := parseURL(nil, rootURL)
	if err != nil {
		return nil, fmt.Errorf("parsing url '%s': %s", rootURL, err)
	}
	rootNode := site.NewURLNode(u, 0)
	sitemap := site.NewSiteMap(rootNode, maxDepth)
	// Get RateLimiter and create Fetcher
	hostname := rootNode.GetHostName()
	limiter := ratelimiter.Get(hostname, ratelimit)
	defer ratelimiter.Stop(hostname)
	fetcher := fetcher.New(limiter, args.timeout)
	// Crawl
	queue = []site.URLNode{rootNode}
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("aborted, %s", ctx.Err())
		default:
			if len(queue) == 0 {
				return sitemap, nil
			}
			parent = queue[0]
			queue = queue[1:]
			result, err = fetcher.Fetch(parent.GetURL().String())
			if err != nil {
				log.Errorf("fetching url '%s': %s", parent.GetURL().String(), err)
			}
			for _, urlString := range result {
				if u, err = parseURL(parent.GetURL(), urlString); err != nil {
					log.Errorf("parsing url '%s': %s", urlString, err)
					continue
				}
				child = site.NewURLNode(u, parent.GetDepth()+1)
				err = sitemap.AddChild(parent, child)
				switch err.(type) {
				case nil:
					queue = append(queue, child)
				case site.ErrInvalidNode:
					log.Debugf("url '%s' skipped due to:%v", urlString, err)
				default:
					log.Errorf("url '%s' skipped due to:%v", urlString, err)
				}
			}
		}
	}
}
