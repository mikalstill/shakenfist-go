package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// GetNamespaces fetches a list of all Namespaces.
func (c *Client) GetNamespaces() ([]string, error) {
	namespaces := []string{}
	err := c.doRequest("auth/namespaces", "GET", bytes.Buffer{}, &namespaces)
	return namespaces, err
}

type createNamespaceReq struct {
	Namespace string `json:"namespace"`
}

// CreateNameSpace creates a new Namespace.
func (c *Client) CreateNamespace(namespace string) error {
	req := &createNamespaceReq{
		Namespace: namespace,
	}

	post, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("cannot marshal namespace req: %v", err)
	}

	err = c.doRequest("auth/namespaces", "POST", *bytes.NewBuffer(post), nil)
	if err != nil {
		return fmt.Errorf("cannot create namespace: %v", err)
	}

	return nil
}

type createNamespaceKeyReq struct {
	Key string `json:"key"`
}

// CreateNameSpaceKey creates a key within a namespace.
func (c *Client) CreateNamespaceKey(namespace, keyName, key string) error {
	req := &createNamespaceKeyReq{
		Key: key,
	}

	post, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("cannot marshal key req: %v", err)
	}

	path := "auth/namespaces/" + namespace + "/keys/" + keyName
	err = c.doRequest(path, "POST", *bytes.NewBuffer(post), nil)
	if err != nil {
		return fmt.Errorf("cannot create namespace: %v", err)
	}

	return nil
}

// UpdateNameSpaceKey will modify an existing Key within a Namespace.
func (c *Client) UpdateNamespaceKey(namespace, keyName, key string) error {
	return c.CreateNamespaceKey(namespace, keyName, key)
}

// GetNameSpaceKeys retrieves a list of keys within the namespace
func (c *Client) GetNamespaceKeys(namespace string) ([]string, error) {
	keyNames := []string{}
	err := c.getRequest("auth/namespaces", namespace, "keys", &keyNames)
	return keyNames, err
}

// DeleteNameSpace attempts to delete the namespace from Shaken Fist.
func (c *Client) DeleteNamespace(namespace string) error {
	path := "auth/namespaces/" + namespace

	err := c.doRequest(path, "DELETE", bytes.Buffer{}, nil)
	if err != nil {
		return fmt.Errorf("unable to delete namespace: %v", err)
	}

	return nil
}

// DeleteNameSpaceKey attempts to delete the key from the specified namespace.
func (c *Client) DeleteNamespaceKey(namespace, keyName string) error {
	path := "auth/namespaces/" + namespace + "/keys/" + keyName

	err := c.doRequest(path, "DELETE", bytes.Buffer{}, nil)
	if err != nil {
		return fmt.Errorf("unable to delete key: %v", err)
	}

	return nil
}
