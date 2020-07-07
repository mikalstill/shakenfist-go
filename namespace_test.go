package client

import (
	"io/ioutil"
	"net/http"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Namespace management functions", func() {
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

	It("should get a list of namespaces", func() {
		reqPath := test_url + "/auth/namespaces"

		// JSON data that SF would return
		testNames := []byte(`["deepspace","bobspace","meatspace"]`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testNames))

		// Make client request
		inst, err := client.GetNamespaces()
		Expect(err).To(BeNil())
		Expect(inst).To(Equal([]string{
			"deepspace",
			"bobspace",
			"meatspace",
		}))

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})

	It("should create a namespace", func() {
		reqPath := test_url + "/auth/namespaces"

		// JSON data expected to be sent
		testSpace := []byte(`{"namespace":"bobspace"}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("POST", reqPath,
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(testSpace))

				return httpmock.NewStringResponse(200, ""), nil
			},
		)

		// Make client request
		err := client.CreateNamespace("bobspace")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should delete a namespace", func() {
		reqPath := test_url + "/auth/namespaces/bobspace"

		// Prepare mocked HTTP
		httpmock.RegisterResponder("DELETE", reqPath,
			httpmock.NewStringResponder(200, ""))

		// Make client request
		err := client.DeleteNamespace("bobspace")
		Expect(err).To(BeNil())
	})
})

var _ = Describe("Namespace Key management functions", func() {
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

	It("should create a namespace key", func() {
		reqPath := test_url + "/auth/namespaces/bobspace/keys/bobskey"

		// JSON data expected to be sent
		testSpace := []byte(`{"key":"supersecret"}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("POST", reqPath,
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(testSpace))

				return httpmock.NewStringResponse(200, ""), nil
			},
		)

		// Make client request
		err := client.CreateNamespaceKey("bobspace", "bobskey", "supersecret")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should update a namespace key", func() {
		reqPath := test_url + "/auth/namespaces/bobspace/keys/bobskey"

		// JSON data expected to be sent
		testSpace := []byte(`{"key":"supersecret"}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("POST", reqPath,
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(testSpace))

				return httpmock.NewStringResponse(200, ""), nil
			},
		)

		// Make client request
		err := client.UpdateNamespaceKey("bobspace", "bobskey", "supersecret")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should delete a namespace key", func() {
		reqPath := test_url + "/auth/namespaces/bobspace/keys/bobskey"

		// Prepare mocked HTTP
		httpmock.RegisterResponder("DELETE", reqPath,
			httpmock.NewStringResponder(200, ""))

		// Make client request
		err := client.DeleteNamespaceKey("bobspace", "bobskey")
		Expect(err).To(BeNil())
	})
})
