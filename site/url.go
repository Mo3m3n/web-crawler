package site

import (
	"encoding/json"
	"net/url"
)

type URLNode interface {
	// GetURL returns the URL struct of the node
	GetURL() *url.URL
	// GetHostName returns the url's hostname
	GetHostName() string
	// GetPath returns the url's path
	GetPath() string
	// GetDepth returns the url's depth
	GetDepth() int
	// MarshalJSON implements json.Marshal interface
	MarshalJSON() ([]byte, error)
}

type node struct {
	url   *url.URL
	depth int
	path  string
	urls  map[string]URLNode
}

func (n *node) GetURL() *url.URL {
	return n.url
}
func (n *node) GetHostName() string {
	return n.url.Hostname()
}
func (n *node) GetPath() string {
	return n.url.Path
}
func (n *node) GetDepth() int {
	return n.depth
}
func (n *node) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(n.urls, "", "  ")
}

// NewURLNode takes a url (as URL pointer from net/url)
// and its depths in the sitemap and creates a new url node
func NewURLNode(u *url.URL, depth int) URLNode {
	return &node{
		url:   u,
		depth: depth,
		path:  u.Path,
		urls:  map[string]URLNode{},
	}
}
