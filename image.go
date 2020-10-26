package client

import (
	"bytes"
	"encoding/json"
)

// ImageRequest defines a link to an image.
type imageRequest struct {
	URL string `json:"url"`
}

// CacheImage will cache an image.
func (c *Client) CacheImage(imageURL string) error {
	request := &imageRequest{
		URL: imageURL,
	}
	post, err := json.Marshal(request)
	if err != nil {
		return err
	}

	err = c.doRequestJSON("images", "POST", *bytes.NewBuffer(post), nil)

	return err
}

// ImageMeta contains the metadata for a cached Image
type ImageMeta struct {
	Checksum    string `json:"checksum"`
	Fetched     string `json:"fetched"`
	FileVersion int    `json:"file_version"`
	Modified    string `json:"modified"`
	Node        string `json:"node"`
	Ref         string `json:"ref"`
	Size        string `json:"size"`
	URL         string `json:"url"`
}

// GetImageMeta retrieves a list of Image metadata
func (c *Client) GetImageMeta() ([]ImageMeta, error) {
	images := []ImageMeta{}
	err := c.doRequestJSON("images", "GET", bytes.Buffer{}, &images)

	return images, err
}
