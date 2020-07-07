package client

import (
	"io/ioutil"
	"net/http"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata management functions", func() {
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

	//
	// Namespace
	//
	It("should get namespace metadata", func() {
		reqPath := test_url + "/auth/namespaces/testspace/metadata"

		// JSON data that SF would return
		testMetdata := []byte(`{
			"name":"bob",
			"pet":"dog"
		}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testMetdata))

		// Make client request
		inst, err := client.GetNamespaceMetadata("testspace")
		Expect(err).To(BeNil())
		Expect(inst).To(Equal(Metadata{
			"name": "bob",
			"pet":  "dog",
		}))

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})

	It("should set namespace metadata", func() {
		reqPath := test_url + "/auth/namespaces/testspace/metadata/name"

		// JSON data that SF would return
		testMetdata := []byte(`{"value":"bob"}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("POST", reqPath,
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(testMetdata))

				return httpmock.NewStringResponse(200, ""), nil
			},
		)

		// Make client request
		err := client.SetNamespaceMetadata("testspace", "name", "bob")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should delete namespace metadata", func() {
		reqPath := test_url + "/auth/namespaces/testspace/metadata/name"

		// Prepare mocked HTTP
		httpmock.RegisterResponder("DELETE", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.DeleteNamespaceMetadata("testspace", "name")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["DELETE "+reqPath]).To(Equal(1))
	})

	//
	// Instance
	//
	It("should get instance metadata", func() {
		reqPath := test_url + "/instances/1a2b-6c7d-eeff/metadata"

		// JSON data that SF would return
		testMetdata := []byte(`{
			"name":"bob",
			"pet":"dog"
		}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testMetdata))

		// Make client request
		inst, err := client.GetInstanceMetadata("1a2b-6c7d-eeff")
		Expect(err).To(BeNil())
		Expect(inst).To(Equal(Metadata{
			"name": "bob",
			"pet":  "dog",
		}))

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})

	It("should set instance metadata", func() {
		reqPath := test_url + "/instances/1a2b-6c7d-eeff/metadata/name"

		// JSON data expected to be sent
		testMetdata := []byte(`{"value":"bob"}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("POST", reqPath,
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(testMetdata))

				return httpmock.NewStringResponse(200, ""), nil
			},
		)

		// Make client request
		err := client.SetInstanceMetadata("1a2b-6c7d-eeff", "name", "bob")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should delete instance metadata", func() {
		reqPath := test_url + "/instances/1a2b-6c7d-eeff/metadata/name"

		// Prepare mocked HTTP
		httpmock.RegisterResponder("DELETE", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.DeleteInstanceMetadata("1a2b-6c7d-eeff", "name")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["DELETE "+reqPath]).To(Equal(1))
	})

	//
	// Network metadata
	//
	It("should get network metadata", func() {
		reqPath := test_url + "/networks/1234-5678/metadata"

		// JSON data that SF would return
		testMetdata := []byte(`{
			"name":"bob",
			"pet":"dog"
		}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testMetdata))

		// Make client request
		inst, err := client.GetNetworkMetadata("1234-5678")
		Expect(err).To(BeNil())
		Expect(inst).To(Equal(Metadata{
			"name": "bob",
			"pet":  "dog",
		}))

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})

	It("should set network metadata", func() {
		reqPath := test_url + "/networks/1234-5678/metadata/name"

		// JSON data that SF would return
		testMetdata := []byte(`{"value":"bob"}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("POST", reqPath,
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(testMetdata))

				return httpmock.NewStringResponse(200, ""), nil
			},
		)

		// Make client request
		err := client.SetNetworkMetadata("1234-5678", "name", "bob")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should delete network metadata", func() {
		reqPath := test_url + "/networks/1234-5678/metadata/name"

		// Prepare mocked HTTP
		httpmock.RegisterResponder("DELETE", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.DeleteNetworkMetadata("1234-5678", "name")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["DELETE "+reqPath]).To(Equal(1))
	})
})
