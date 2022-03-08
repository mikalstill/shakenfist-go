package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Blob struct {
	UUID           string
	Instances      []string
	Size           int
	ReferenceCount int    `json:"reference_count"`
	DependsOn      string `json:"depends_on"`
}

func (c *Client) GetBlobs(node string) ([]Blob, error) {
	var blobs []Blob

	path := "blobs"
	r := &struct {
		Node string `json:"node"`
	}{
		Node: node,
	}
	req, err := json.Marshal(r)
	if err != nil {
		return []Blob{}, fmt.Errorf("Unable to marshal data: %v", err)
	}

	err = c.doRequestJSON(path, "GET", *bytes.NewBuffer(req), &blobs)
	return blobs, err
}
