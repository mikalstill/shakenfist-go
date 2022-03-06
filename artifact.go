package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//
func (c *Client) CacheArtifact(image_url string) error {
	path := "artifacts"
	r := &struct {
		URL string `json:"url"`
	}{
		URL: image_url,
	}

	req, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("Unable to marshal data: %v", err)
	}

	err = c.doRequestJSON(path, "POST", *bytes.NewBuffer(req), nil)
	return err
}

type Artifact struct {
	UUID            string
	Type            string `json:"artifact_type"`
	State           string
	SourceURL       string `json:"source_url"`
	Version         int
	MaxVersions     int `json:"max_versions"`
	MostRecentIndex int `json:"index"`
	Blobs           map[int]Blob
}

func (c *Client) GetArtifact(uuid string) (Artifact, error) {
	var artifact Artifact

	path := "artifacts/" + uuid
	err := c.doRequestJSON(path, "GET", bytes.Buffer{}, &artifact)
	return artifact, err
}

func (c *Client) GetArtifacts(node string) ([]Artifact, error) {
	var artifacts []Artifact

	path := "artifacts"
	r := &struct {
		Node string `json:"node"`
	}{
		Node: node,
	}
	req, err := json.Marshal(r)
	if err != nil {
		return []Artifact{}, fmt.Errorf("Unable to marshal data: %v", err)
	}

	err = c.doRequestJSON(path, "GET", *bytes.NewBuffer(req), &artifacts)
	return artifacts, err
}

func (c *Client) GetArtifactEvents(uuid string) ([]Event, error) {
	var events []Event

	path := "artifacts/" + uuid + "/events"
	err := c.doRequestJSON(path, "GET", bytes.Buffer{}, &events)
	return events, err
}

func (c *Client) GetArtifactVersions(uuid string) ([]Blob, error) {
	var blobs []Blob

	path := "artifacts/" + uuid + "/versions"
	err := c.doRequestJSON(path, "GET", bytes.Buffer{}, &blobs)
	return blobs, err
}

/***
Missing API calls

def set_artifact_max_versions(self, artifact_uuid, max_versions):
r = self._request_url('POST',
					  '/artifacts/' + artifact_uuid + '/versions',
					  data={'max_versions': max_versions})
return r.json()

def delete_artifact(self, artifact_uuid):
r = self._request_url('DELETE', '/artifacts/' + artifact_uuid)
return r.json()

def delete_artifact_version(self, artifact_uuid, version_id):
r = self._request_url('DELETE', '/artifacts/' + artifact_uuid +
					  '/versions/' + str(version_id))
return r.json()

***/
