package client

import (
	"bytes"
)

// Node defines a ShakenFist node.
type Node struct {
	Name     string  `json:"name"`
	IP       string  `json:"ip"`
	LastSeen float64 `json:"lastseen"`
}

// GetNodes fetches a list of nodes.
func (c *Client) GetNodes() ([]Node, error) {
	nodes := []Node{}
	err := c.doRequestJSON("nodes", "GET", bytes.Buffer{}, &nodes)
	return nodes, err
}
