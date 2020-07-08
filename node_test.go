package client

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Node management functions", func() {
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

	It("should get a list of nodes", func() {
		reqPath := test_url + "/nodes"

		// JSON data that SF would return
		testJSON := []byte(`[
			{
				"name":"sf-1",
				"ip":"10.0.1.1",
				"lastseen":"1594251513.6553159"
			},
			{
				"name":"sf-2",
				"ip":"10.0.1.2",
				"lastseen":"1594251513.7450194"
			}
		]`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testJSON))

		// Make client request
		nodes, err := client.GetNodes()
		Expect(err).To(BeNil())
		Expect(nodes).To(Equal([]Node{
			{
				Name:     "sf-1",
				IP:       "10.0.1.1",
				LastSeen: "1594251513.6553159",
			},
			{
				Name:     "sf-2",
				IP:       "10.0.1.2",
				LastSeen: "1594251513.7450194",
			},
		}))
	})
})
