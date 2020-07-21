package client

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Instance management functions", func() {
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

	It("should get an instance", func() {
		testInstance := []byte(`{
			"uuid":"123-456",
			"name":"test",
			"cpus":1,
			"memory":1024,
			"disk_spec":[{
				"base":"DiskBase",
				"size":5,
				"bus":"DiskBus",
				"type":"DiskType"
			}],
			"ssh_key":"longSSHKey",
			"node":"somenode",
			"console_port":1234,
			"vdi_port":678,
			"user_data":"long story",
			"block_devices":{},
			"state":"nice",
			"state_updated":1.2,
			"power_state":"created"
		}`)

		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456"
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, testInstance))

		// Make client request
		inst, err := client.GetInstance("123-456")
		Expect(err).To(BeNil())
		Expect(inst).To(Equal(Instance{
			UUID:   "123-456",
			Name:   "test",
			CPUs:   1,
			Memory: 1024,
			DiskSpecs: []DiskSpec{
				{
					Base: "DiskBase",
					Size: 5,
					Bus:  "DiskBus",
					Type: "DiskType",
				},
			},
			SSHKey:       "longSSHKey",
			Node:         "somenode",
			ConsolePort:  1234,
			VDIPort:      678,
			UserData:     "long story",
			BlockDevices: map[string]interface{}{},
			State:        "nice",
			StateUpdated: 1.2,
			PowerState:   "created",
		}))

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})

	It("should create an instance", func() {

		testInstance := []byte(`{
			"uuid":"123-456",
			"name":"test",
			"cpus":1,
			"memory":1024,
			"disk_spec":[{
				"base":"DiskBase",
				"size":5,
				"bus":"DiskBus",
				"type":"DiskType"
			}],
			"ssh_key":"longSSHKey",
			"node":"somenode",
			"console_port":1234,
			"vdi_port":678,
			"user_data":"long story",
			"block_devices":{},
			"state":"nice",
			"state_updated":1.2
		}`)

		// Prepare mocked HTTP
		reqPath := test_url + "/instances"
		httpmock.RegisterResponder("POST", reqPath,
			httpmock.NewBytesResponder(200, testInstance))

		// Make client request
		inst, err := client.CreateInstance(
			"test", 1, 1024,
			[]NetworkSpec{
				{
					NetworkUUID: "123",
				},
			},
			[]DiskSpec{
				{
					Base: "DiskBase",
					Size: 5,
					Bus:  "DiskBus",
					Type: "DiskType",
				},
			},
			"longSSHKey",
			"long story")

		Expect(err).To(BeNil())
		Expect(inst).To(Equal(Instance{
			UUID:   "123-456",
			Name:   "test",
			CPUs:   1,
			Memory: 1024,
			DiskSpecs: []DiskSpec{
				{
					Base: "DiskBase",
					Size: 5,
					Bus:  "DiskBus",
					Type: "DiskType",
				},
			},
			SSHKey:       "longSSHKey",
			Node:         "somenode",
			ConsolePort:  1234,
			VDIPort:      678,
			UserData:     "long story",
			BlockDevices: map[string]interface{}{},
			State:        "nice",
			StateUpdated: 1.2,
		}))
	})

	It("should make a reboot request", func() {
		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456/reboot"
		httpmock.RegisterResponder("POST", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.RebootInstance("123-456")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should make a poweroff request", func() {
		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456/poweroff"
		httpmock.RegisterResponder("POST", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.PowerOffInstance("123-456")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should make a poweron request", func() {
		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456/poweron"
		httpmock.RegisterResponder("POST", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.PowerOnInstance("123-456")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should make a pause instance request", func() {
		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456/pause"
		httpmock.RegisterResponder("POST", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.PauseInstance("123-456")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["POST "+reqPath]).To(Equal(1))
	})

	It("should delete an instance request", func() {
		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456"
		httpmock.RegisterResponder("DELETE", reqPath,
			httpmock.NewBytesResponder(200, nil))

		// Make client request
		err := client.DeleteInstance("123-456")
		Expect(err).To(BeNil())

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["DELETE "+reqPath]).To(Equal(1))
	})

	It("should retrieve console data from the instance", func() {
		jsonResp := []byte(`=== cirros: current=0.5.1 latest=0.5.1 uptime=4.20 ===
		____               ____  ____
	   / __/ __ ____ ____ / __ \/ __/
	  / /__ / // __// __// /_/ /\ \
	  \___//_//_/  /_/   \____/___/
		 http://cirros-cloud.net`)

		// Prepare mocked HTTP
		reqPath := test_url + "/instances/123-456/consoledata"
		httpmock.RegisterResponder("GET", reqPath,
			httpmock.NewBytesResponder(200, jsonResp))

		// Make client request
		consoledata, err := client.GetConsoleData("123-456", 100)
		Expect(err).To(BeNil())
		Expect(consoledata).To(Equal(`=== cirros: current=0.5.1 latest=0.5.1 uptime=4.20 ===
		____               ____  ____
	   / __/ __ ____ ____ / __ \/ __/
	  / /__ / // __// __// /_/ /\ \
	  \___//_//_/  /_/   \____/___/
		 http://cirros-cloud.net`))

		// Check correct URL requested
		info := httpmock.GetCallCountInfo()
		Expect(info["GET "+reqPath]).To(Equal(1))
	})
})
