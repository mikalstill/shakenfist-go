package client

import (
	"bytes"
)

// Node defines a ShakenFist node
type Node struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	LastSeen string `json:"lastseen"`
}

// GetNodes fetches a list of nodes
func (c *Client) GetNodes() ([]Node, error) {
	nodes := []Node{}
	err := c.doRequest("nodes", "GET", bytes.Buffer{}, &nodes)
	return nodes, err
}
