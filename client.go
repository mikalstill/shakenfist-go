// Client is an API wrapper to access the Shaken Fist HTTP API.
package client

// With many thanks to the example code from
// https://github.com/spaceapegames/terraform-provider-example

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResourceType is a type of Shaken Fist resource
type ResourceType int

const (
	TypeNamespace ResourceType = iota
	TypeInstance
	TypeNetwork
)

func (r ResourceType) String() string {
	return [...]string{"auth/namespaces", "instances", "networks"}[r]
}

// Client holds all of the information required to connect to
// the server
type Client struct {
	hostname   string
	port       int
	httpClient *http.Client
	namespace  string
	apiKey     string
	cachedAuth string
}

// NewClient returns a client ready for use
func NewClient(hostname string, port int,
	namespace, apiKey string) *Client {

	return &Client{
		hostname:   hostname,
		port:       port,
		httpClient: &http.Client{},
		namespace:  namespace,
		apiKey:     apiKey,
	}
}

//
// Internal helper functions
//

func (c *Client) getRequest(
	object, uuid string, cmd string, resp interface{}) error {

	err := c.doRequest(object+"/"+uuid+"/"+cmd, "GET", bytes.Buffer{}, resp)
	return err
}

func (c *Client) postRequest(object string, uuid string, cmd string) error {
	err := c.doRequest(object+"/"+uuid+"/"+cmd, "POST", bytes.Buffer{}, nil)
	return err
}

func (c *Client) doRequest(
	path, method string, data bytes.Buffer, resp interface{}) error {

	if c.cachedAuth == "" {
		err := c.requestAuth()
		if err != nil {
			return fmt.Errorf("unable to get auth token: %v", err)
		}
	}

	body, statusCode, err := c.httpRequest(path, method, data)

	// If auth token has expired, then get a new token
	if statusCode == http.StatusUnauthorized {
		if c.requestAuth() != nil {
			return fmt.Errorf("unable to refresh auth token: %v", err)
		}

		// Try with new token, if second error occurs it is returned
		body, _, err = c.httpRequest(path, method, data)
	}

	// Return on error or if JSON decoding not required
	if err != nil || resp == nil {
		return err
	}

	return json.NewDecoder(body).Decode(resp)
}

func (c *Client) httpRequest(
	path, method string, body bytes.Buffer) (io.ReadCloser, int, error) {

	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", c.cachedAuth)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to connect to server: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, resp.StatusCode,
				fmt.Errorf("received non 200 status code: %v",
					resp.StatusCode)
		}
		return nil, resp.StatusCode,
			fmt.Errorf("received non 200 status code: %v - %s",
				resp.StatusCode, respBody.String())
	}
	return resp.Body, 0, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("http://%s:%v/%s", c.hostname, c.port, path)
}

type authRequest struct {
	Namespace string `json:"namespace"`
	APIKey    string `json:"key"`
}

type authResponse struct {
	Token string `json:"access_token"`
}

func (c *Client) requestAuth() error {
	req := &authRequest{
		Namespace: c.namespace,
		APIKey:    c.apiKey,
	}
	post, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("unable to marshal auth request: %v", err)
	}

	body, _, err := c.httpRequest("/auth", "POST", *bytes.NewBuffer(post))
	if err != nil {
		return fmt.Errorf("auth request failed: %v", err)
	}

	resp := authResponse{}
	err = json.NewDecoder(body).Decode(&resp)
	if err != nil {
		return fmt.Errorf("unable to decode response body: %v", err)
	}

	c.cachedAuth = fmt.Sprintf("Bearer %s", resp.Token)

	return nil
}
