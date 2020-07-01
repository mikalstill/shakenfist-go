package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// GetNamespaces fetches a list of all Namespaces.
func (c *Client) GetNameSpaces() ([]string, error) {
	namespaces := []string{}
	err := c.doRequest("auth/namespace", "GET", bytes.Buffer{}, &namespaces)
	return namespaces, err
}

type createNamespaceReq struct {
	Namespace string `json:"namespace"`
	KeyName   string `json:"key_name"`
	Key       string `json:"key"`
}

// CreateNameSpace creates a new Namespace, KeyName and Key.
//
// Use this function to create a new Namespace and also to create
// a new key.
func (c *Client) CreateNameSpace(namespace, keyName, key string) error {
	req := &createNamespaceReq{
		Namespace: namespace,
		KeyName:   keyName,
		Key:       key,
	}

	post, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("cannot marshal namespace req: %v", err)
	}

	err = c.doRequest("auth/namespace", "POST", *bytes.NewBuffer(post), nil)
	if err != nil {
		return fmt.Errorf("cannot create namespace: %v", err)
	}

	return nil
}

// UpdateNameSpaceKey will modify an existing Key within a Namespace.
func (c *Client) UpdateNameSpaceKey(namespace, keyName, key string) error {
	return c.CreateNameSpace(namespace, keyName, key)
}

// DeleteNameSpace attempts to delete the namespace from Shaken Fist.
func (c *Client) DeleteNameSpace(namespace string) error {
	path := "auth/namespace/" + namespace

	err := c.doRequest(path, "DELETE", bytes.Buffer{}, nil)
	if err != nil {
		return fmt.Errorf("unable to delete namespace: %v", err)
	}

	return nil
}

// DeleteNameSpaceKey attempts to delete the key from the specified namespace.
func (c *Client) DeleteNameSpaceKey(namespace, keyName string) error {
	path := "auth/namespace/" + namespace + "/" + keyName

	err := c.doRequest(path, "DELETE", bytes.Buffer{}, nil)
	if err != nil {
		return fmt.Errorf("unable to delete key: %v", err)
	}

	return nil
}
