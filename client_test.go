package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client timeout", func() {
	const (
		test_url       string = "http://server:13000"
		test_namespace string = "testspace"
		test_key       string = "testkey"
	)

	var (
		client *Client
	)

	BeforeEach(func() {
		// Configure client
		client = NewClient(test_url, test_namespace, test_key)
	})

	It("should have the default timeout", func() {
		Expect(client.httpClient.Timeout).To(Equal(time.Duration(0)))
	})

	It("should set the timeout", func() {
		client.SetTimeout(123)
		Expect(client.httpClient.Timeout).To(Equal(123 * time.Second))
	})
})

var _ = Describe("Auth request", func() {
	const (
		test_url       string = "http://server:13000"
		test_namespace string = "testspace"
		test_key       string = "testkey"
	)

	var (
		client *Client
	)

	BeforeEach(func() {
		// Configure client
		client = NewClient(test_url, test_namespace, test_key)
	})

	It("should send the correct auth data", func() {

		testAuth := []byte(`{"namespace":"testspace","key":"testkey"}`)

		httpmock.RegisterResponder("POST", test_url+"/auth",
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(testAuth))

				token := []byte(`{"access_token":"ABC123"}`)
				return httpmock.NewBytesResponse(200, token), nil
			})

		// Make auth request
		err := client.requestAuth()
		Expect(err).To(BeNil())

		Expect(client.cachedAuth).To(Equal("Bearer ABC123"))

		// Check auth request was made only once
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+test_url+"/auth"]).To(Equal(1))
	})
})

var _ = Describe("Auth caching", func() {
	const (
		test_url       string = "http://server:13000"
		test_namespace string = "testspace"
		test_key       string = "testkey"
	)

	var (
		client *Client
	)

	BeforeEach(func() {
		// Configure client
		client = NewClient(test_url, test_namespace, test_key)

		httpmock.RegisterResponder("POST", test_url+"/auth",
			httpmock.NewBytesResponder(200, []byte(`{"access_token":"ABC123"}`)))
	})

	It("should get and store auth token", func() {

		// Make auth request
		err := client.requestAuth()
		Expect(err).To(BeNil())

		Expect(client.cachedAuth).To(Equal("Bearer ABC123"))

		// Check auth request was made only once
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+test_url+"/auth"]).To(Equal(1))
	})

	It("should get and store auth token", func() {

		// Make auth request
		err := client.requestAuth()
		Expect(err).To(BeNil())

		Expect(client.cachedAuth).To(Equal("Bearer ABC123"))

		// Check auth request was made only once
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+test_url+"/auth"]).To(Equal(1))
	})

	It("should get auth token and use it for subsequent requests", func() {
		instances := []Instance{}

		// Prepare mocked HTTP
		jsonResp, err := json.Marshal([]Instance{})
		Expect(err).ShouldNot(HaveOccurred())

		// Make one request
		httpmock.RegisterResponder("GET", test_url+"/instances",
			httpmock.NewBytesResponder(200, jsonResp))

		err = client.doRequestJSON("instances", "GET", bytes.Buffer{}, &instances)
		Expect(err).To(BeNil())

		// Make second request, expecting auth token to be cached
		httpmock.RegisterResponder("GET", test_url+"/instances",
			httpmock.NewBytesResponder(200, jsonResp))

		err = client.doRequestJSON("instances", "GET", bytes.Buffer{}, &instances)
		Expect(err).To(BeNil())

		// Check auth request was made only once
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+test_url+"/auth"]).To(Equal(1))
	})
})

var _ = Describe("Helper functions", func() {
	const (
		test_url       string = "http://server:13000"
		test_namespace string = "testspace"
		test_key       string = "testkey"
	)

	var (
		client *Client
	)

	BeforeEach(func() {
		// Configure client
		client = NewClient(test_url, test_namespace, test_key)

		httpmock.RegisterResponder("POST", test_url+"/auth",
			httpmock.NewBytesResponder(200, []byte(`{"access_token":"ABC123"}`)))
	})

	It("should make a GET request", func() {
		instances := []Instance{}

		// Prepare mocked HTTP
		jsonResp, err := json.Marshal([]Instance{})
		Expect(err).To(BeNil())

		reqPath := test_url + "/instances/123-456/cmd"
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, jsonResp))

		// Make client request
		err = client.getRequest("instances", "123-456", "cmd", &instances)
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})

	It("should make a POST request", func() {
		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456/cmd"
		httpmock.RegisterResponder("POST", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.postRequest("instances", "123-456", "cmd")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should fail when wrong path requested", func() {
		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456/cmd"
		httpmock.RegisterResponder("POST", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.postRequest("wrong", "123-456", "cmd")
		Expect(err).ToNot(BeNil())
	})
})
