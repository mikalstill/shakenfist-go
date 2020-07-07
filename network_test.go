package client

import (
	"io/ioutil"
	"net/http"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Network management functions", func() {
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

	It("should get a list of networks", func() {
		reqPath := test_url + "/networks"

		// JSON data that SF would return
		testNetJSON := []byte(`[
			{
			"uuid":"1234-5678",
			"name":"bobnet",
			"vxid":112222,
			"netblock":"10.0.1.0/24",
			"provide_dhcp":true,
			"provide_nat":true,
			"owner":"1234",
			"floating_gateway":"10.0.1.250",
			"state":"created",
			"state_updated":4564.5
			},
			{
				"uuid":"abab-cdcd",
				"name":"bobnet2",
				"vxid":112222,
				"netblock":"10.0.1.0/24",
				"provide_dhcp":true,
				"provide_nat":true,
				"owner":"1234",
				"floating_gateway":"10.0.1.250",
				"state":"created",
				"state_updated":4564.5
			}
		]`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testNetJSON))

		// Make client request
		net, err := client.GetNetworks()
		Expect(err).To(BeNil())
		Expect(net).To(Equal([]Network{
			{
				UUID:            "1234-5678",
				Name:            "bobnet",
				VXId:            112222,
				NetBlock:        "10.0.1.0/24",
				ProvideDHCP:     true,
				ProvideNAT:      true,
				Owner:           "1234",
				FloatingGateway: "10.0.1.250",
				State:           "created",
				StateUpdated:    4564.5,
			},
			{
				UUID:            "abab-cdcd",
				Name:            "bobnet2",
				VXId:            112222,
				NetBlock:        "10.0.1.0/24",
				ProvideDHCP:     true,
				ProvideNAT:      true,
				Owner:           "1234",
				FloatingGateway: "10.0.1.250",
				State:           "created",
				StateUpdated:    4564.5,
			},
		}))
	})

	It("should get a network", func() {
		reqPath := test_url + "/networks/1234-5678"

		// JSON data that SF would return
		testNetJSON := []byte(`{
			"uuid":"1234-5678",
			"name":"bobnet",
			"vxid":112222,
			"netblock":"10.0.1.0/24",
			"provide_dhcp":true,
			"provide_nat":true,
			"owner":"1234",
			"floating_gateway":"10.0.1.250",
			"state":"created",
			"state_updated":4564.5
			}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testNetJSON))

		// Make client request
		net, err := client.GetNetwork("1234-5678")
		Expect(err).To(BeNil())
		Expect(net).To(Equal(Network{
			UUID:            "1234-5678",
			Name:            "bobnet",
			VXId:            112222,
			NetBlock:        "10.0.1.0/24",
			ProvideDHCP:     true,
			ProvideNAT:      true,
			Owner:           "1234",
			FloatingGateway: "10.0.1.250",
			State:           "created",
			StateUpdated:    4564.5,
		}))
	})

	It("should create a network", func() {
		reqPath := test_url + "/networks"

		// JSON data expected to be sent
		sentData := []byte(`{"name":"nowherenet","netblock":"10.0.1.0/24","provide_dhcp":true,"provide_nat":true}`)

		// JSON data that SF would return
		testNetJSON := []byte(`{
			"uuid":"1234-5678",
			"name":"nowherenet",
			"vxid":112222,
			"netblock":"10.0.1.0/24",
			"provide_dhcp":true,
			"provide_nat":true,
			"owner":"1234",
			"floating_gateway":"10.0.1.250",
			"state":"created",
			"state_updated":4564.5
			}`)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("POST", reqPath,
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(sentData))

				return httpmock.NewBytesResponse(200, testNetJSON), nil
			},
		)

		// Make client request
		net, err := client.CreateNetwork(
			"10.0.1.0/24",
			true,
			true,
			"nowherenet",
		)
		Expect(err).To(BeNil())
		Expect(net).To(Equal(Network{
			UUID:            "1234-5678",
			Name:            "nowherenet",
			VXId:            112222,
			NetBlock:        "10.0.1.0/24",
			ProvideDHCP:     true,
			ProvideNAT:      true,
			Owner:           "1234",
			FloatingGateway: "10.0.1.250",
			State:           "created",
			StateUpdated:    4564.5,
		}))
	})

	It("should delete a network", func() {
		reqPath := test_url + "/networks/1234-5678"

		// Prepare mocked HTTP
		httpmock.RegisterResponder("DELETE", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.DeleteNetwork("1234-5678")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["DELETE "+reqPath]).To(Equal(1))
	})

	It("should fail deleting a non-existent network", func() {
		reqPath := test_url + "/networks/1234-5678"

		// Prepare mocked HTTP
		httpmock.RegisterResponder("DELETE", reqPath,
			httpmock.NewBytesResponder(404, nil))

		// Make client request
		err := client.DeleteNetwork("1234-5678")
		Expect(err).ToNot(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["DELETE "+reqPath]).To(Equal(1))
	})
})
