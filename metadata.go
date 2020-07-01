package client

// Metadata key-value set and retrieval on a specific instance.

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Metadata is a map of key value pairs storing metadata on an instance.
type Metadata map[string]string

// GetMetadata retrieves the metadata attached to an instance.
func (c *Client) GetMetadata(res ResourceType, uuid string) (Metadata, error) {
	path := res.String() + "/" + uuid + "/metadata"

	meta := Metadata{}
	if err := c.doRequest(path, "GET", bytes.Buffer{}, &meta); err != nil {
		return meta, fmt.Errorf("unable to retrieve metadata: %v", err)
	}

	return meta, nil
}

type reqMeta struct {
	Value string `json:"value"`
}

// SetMetadata sets key-value metadata on an instance.
func (c *Client) SetMetadata(res ResourceType, uuid, key, value string) error {
	path := res.String() + "/" + uuid + "/metadata/" + key

	req := &reqMeta{
		Value: value,
	}

	post, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("cannot marshal data into JSON: %v", err)
	}

	if err := c.doRequest(path, "POST", *bytes.NewBuffer(post), nil); err != nil {
		return fmt.Errorf("unable to set metadata: %v", err)
	}

	return nil
}
