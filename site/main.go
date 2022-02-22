package site

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SiteMap interface {
	// AddChild adds a valid url child node to the parent.
	// Otherwise returns an error.
	// a valid url child should:
	// - have the same Hostname as the parent
	// - have its parent path as a prefix of its own path
	// - have its depth less than max depth
	// - not have been visited already
	AddChild(parent, child URLNode) error
	// Marshal returns the json encoding of the sitemap
	Marshal() ([]byte, error)
}

type ErrInvalidNode error
type ErrCachedNode error

type siteMap struct {
	visited  map[string]bool
	root     URLNode
	maxDepth int
}

// NewSiteMap takes the root url of the site and the maximum
// depth of the sitemap and returns a new sitemap
func NewSiteMap(root URLNode, maxDepth int) siteMap {
	return siteMap{
		visited:  make(map[string]bool),
		root:     root,
		maxDepth: maxDepth,
	}
}

func (s siteMap) AddChild(parent, child URLNode) error {
	var err error
	nodeDepth := child.GetDepth()
	nodeHostName := child.GetHostName()
	parentHostName := parent.GetHostName()
	nodePath := child.GetPath()
	parentPath := parent.GetPath()
	if s.maxDepth >= 0 && nodeDepth > s.maxDepth {
		return ErrInvalidNode(fmt.Errorf("url depth '%d' exceeds max depth '%d'", nodeDepth, s.maxDepth))
	}
	if s.visited[nodePath] {
		err = ErrCachedNode(fmt.Errorf("path '%s' was already visited", child.GetURL().String()))
	}
	if nodeHostName != parentHostName {
		err = ErrInvalidNode(fmt.Errorf("hostname '%s' is different fromt parent one '%s'", nodeHostName, parentHostName))
	}
	if !strings.HasPrefix(nodePath, parentPath) {
		err = ErrInvalidNode(fmt.Errorf("path '%s' is different fromt parent one '%s'", nodePath, parentPath))
	}
	node := parent.(*node)
	if node == nil {
		return fmt.Errorf("internal error with url node '%v'", parent)
	}
	if _, ok := node.urls[nodePath]; !ok {
		node.urls[nodePath] = child
	}
	if err == nil {
		s.visited[nodePath] = true
	}
	return err
}

func (s siteMap) Marshal() ([]byte, error) {
	return json.MarshalIndent(s.root, "", "  ")
}
