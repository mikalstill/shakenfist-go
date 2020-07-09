package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// GetNamespaces fetches a list of all Namespaces.
func (c *Client) GetNamespaces() ([]string, error) {
	namespaces := []string{}
	err := c.doRequestJSON("auth/namespaces", "GET", bytes.Buffer{}, &namespaces)
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

	err = c.doRequestJSON("auth/namespaces", "POST", *bytes.NewBuffer(post), nil)
	if err != nil {
		return fmt.Errorf("cannot create namespace: %v", err)
	}

	return nil
}

type createNamespaceKeyReq struct {
	KeyName string `json:"key_name"`
	Key     string `json:"key"`
}

// CreateNameSpaceKey creates a key within a namespace.
func (c *Client) CreateNamespaceKey(namespace, keyName, key string) error {
	req := &createNamespaceKeyReq{
		KeyName: keyName,
		Key:     key,
	}

	post, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("cannot marshal key req: %v", err)
	}

	path := "auth/namespaces/" + namespace + "/keys"
	err = c.doRequestJSON(path, "POST", *bytes.NewBuffer(post), nil)
	if err != nil {
		return fmt.Errorf("cannot create namespace: %v", err)
	}

	return nil
}

type updateNamespaceKeyReq struct {
	Key string `json:"key"`
}

// UpdateNameSpaceKey will modify an existing Key within a Namespace.
func (c *Client) UpdateNamespaceKey(namespace, keyName, key string) error {
	req := &updateNamespaceKeyReq{
		Key: key,
	}

	put, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("cannot marshal key req: %v", err)
	}

	path := "auth/namespaces/" + namespace + "/keys/" + keyName
	err = c.doRequestJSON(path, "PUT", *bytes.NewBuffer(put), nil)
	if err != nil {
		return fmt.Errorf("cannot create namespace: %v", err)
	}

	return nil
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

	err := c.doRequestJSON(path, "DELETE", bytes.Buffer{}, nil)
	if err != nil {
		return fmt.Errorf("unable to delete namespace: %v", err)
	}

	return nil
}

// DeleteNameSpaceKey attempts to delete the key from the specified namespace.
func (c *Client) DeleteNamespaceKey(namespace, keyName string) error {
	path := "auth/namespaces/" + namespace + "/keys/" + keyName

	err := c.doRequestJSON(path, "DELETE", bytes.Buffer{}, nil)
	if err != nil {
		return fmt.Errorf("unable to delete key: %v", err)
	}

	return nil
}
