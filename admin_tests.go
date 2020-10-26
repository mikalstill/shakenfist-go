package client

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Image management functions", func() {
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

	It("should get a list of locks", func() {
		testInstance := []byte(`{
			"/sflocks/sf/queue/sf-2": {
				"node": "sf-2",
				"operation": "Restart",
				"pid": 1142
			},
			"/sflocks/sf/network/8a3b3f4c-9eed-4fa6-9136-49d03a0859ea": {
				"node": "sf-1",
				"operation": null,
				"pid": 10083
			},
			"/sflocks/sf/queue/networknode": {
				"node": "sf-3",
				"operation": "Enqueue",
				"pid": 7083
			}
		}`)

		// Prepare mocked HTTP
		reqPath := test_url + "/admin/locks"
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testInstance))

		// Make client request
		inst, err := client.GetLocks()
		Expect(err).To(BeNil())
		Expect(inst).To(Equal(Locks{
			"/sflocks/sf/queue/sf-2": {
				Node:      "sf-2",
				Operation: "Restart",
				PID:       1142,
			},
			"/sflocks/sf/network/8a3b3f4c-9eed-4fa6-9136-49d03a0859ea": {
				Node:      "sf-1",
				Operation: "",
				PID:       10083,
			},
			"/sflocks/sf/queue/networknode": {
				Node:      "sf-3",
				Operation: "Enqueue",
				PID:       7083,
			},
		},
		))

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})
})
