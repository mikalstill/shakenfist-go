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

	It("should make a delete all networks request", func() {
		// Prepare mocked HTTP
		reqPath := test_url + "/networks"
		expReqData := []byte(`{"namespace":"bobspace","confirm":true}`)
		jsonResp := `["1234-im-a-uuid", "1234abcd-uuid-really"]`

		httpmock.RegisterResponder("DELETE", reqPath,
			func(req *http.Request) (*http.Response, error) {
				buf, err := ioutil.ReadAll(req.Body)

				Expect(err).To(BeNil())
				Expect(buf).To(Equal(expReqData))

				return httpmock.NewStringResponse(200, jsonResp), nil
			},
		)

		// Make client request
		networks, err := client.DeleteAllNetworks("bobspace")
		Expect(err).To(BeNil())
		Expect(networks).To(Equal(
			[]string{"1234-im-a-uuid", "1234abcd-uuid-really"}))

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

	It("should get a list of interfaces on the network", func() {
		reqPath := test_url + "/networks/1234-5678/interfaces"

		// JSON data that SF would return
		testNetJSON := []byte(`[
		{
		  "floating": null,
		  "instance_uuid": "5f01aa37-7df5-4bbb-867b-63606d964802",
		  "ipv4": "10.0.0.234",
		  "macaddr": "00:00:00:a8:58:f9",
		  "model": "virtio",
		  "network_uuid": "0e766fda-b5fc-40b1-9f96-e69bbb5cf590",
		  "order": 0,
		  "state": "created",
		  "state_updated": 1596169109.229394,
		  "uuid": "372e4f5c-7d5f-4c7b-a1b2-512ddb7da82a"
		},
		{
		  "floating": "10.10.0.175",
		  "instance_uuid": "7594c492-a27f-4da3-ba2a-19fdbe7e14a2",
		  "ipv4": "10.0.0.34",
		  "macaddr": "00:00:00:af:14:16",
		  "model": "virtio",
		  "network_uuid": "0e766fda-b5fc-40b1-9f96-e69bbb5cf590",
		  "order": 0,
		  "state": "created",
		  "state_updated": 1596166987.5956044,
		  "uuid": "77e3b00f-89a1-4fc0-b595-8f6ed092f8e7"
		}
	  ]
	  `)

		// Prepare mocked HTTP
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testNetJSON))

		// Make client request
		net, err := client.GetNetworkInterfaces("1234-5678")
		Expect(err).To(BeNil())
		Expect(net).To(Equal([]NetworkInterface{
			{
				Floating:     "",
				InstanceUUID: "5f01aa37-7df5-4bbb-867b-63606d964802",
				IPv4:         "10.0.0.234",
				MACAddress:   "00:00:00:a8:58:f9",
				Model:        "virtio",
				NetworkUUID:  "0e766fda-b5fc-40b1-9f96-e69bbb5cf590",
				Order:        0,
				State:        "created",
				StateUpdated: 1596169109.229394,
				UUID:         "372e4f5c-7d5f-4c7b-a1b2-512ddb7da82a",
			},
			{
				Floating:     "10.10.0.175",
				InstanceUUID: "7594c492-a27f-4da3-ba2a-19fdbe7e14a2",
				IPv4:         "10.0.0.34",
				MACAddress:   "00:00:00:af:14:16",
				Model:        "virtio",
				NetworkUUID:  "0e766fda-b5fc-40b1-9f96-e69bbb5cf590",
				Order:        0,
				State:        "created",
				StateUpdated: 1596166987.5956044,
				UUID:         "77e3b00f-89a1-4fc0-b595-8f6ed092f8e7",
			},
		}))
	})
})
