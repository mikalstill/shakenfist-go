// With many thanks to the example code from
// https://github.com/spaceapegames/terraform-provider-example
package client

// Note that the following API calls are not yet implemented as
// they are not needed for the terraform provider, which is the
// primary user of this client:
//

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client holds all of the information required to connect to
// the server
type Client struct {
	hostname   string
	port       int
	httpClient *http.Client
}

// NewClient returns a client ready for use
func NewClient(hostname string, port int) *Client {
	return &Client{
		hostname:   hostname,
		port:       port,
		httpClient: &http.Client{},
	}
}

// Network is a definition of a network
type Network struct {
	UUID            string  `json:"uuid"`
	Name            string  `json:"name"`
	VXId            int     `json:"vxid"`
	NetBlock        string  `json:"netblock"`
	ProvideDHCP     bool    `json:"provide_dhcp"`
	ProvideNAT      bool    `json:"provide_nat"`
	Owner           string  `json:"owner"`
	FloatingGateway string  `json:"floating_gateway"`
	State           string  `json:"state"`
	StateUpdated    float64 `json:"state_updated"`
}

// GetNetworks fetches a list of networks
func (c *Client) GetNetworks() ([]Network, error) {
	networks := []Network{}
	err := c.doRequest("networks", "GET", bytes.Buffer{}, &networks)
	return networks, err
}

// GetNetwork fetches a specific instance by UUID
func (c *Client) GetNetwork(uuid string) (Network, error) {
	network := Network{}
	err := c.doRequest("networks/"+uuid, "GET", bytes.Buffer{}, &network)
	return network, err
}

type createNetworkRequest struct {
	Name        string `json:"name"`
	Netblock    string `json:"netblock"`
	ProvideDHCP bool   `json:"provide_dhcp"`
	ProvideNAT  bool   `json:"provide_nat"`
}

// CreateNetwork creates a new network
func (c *Client) CreateNetwork(netblock string, provideDHCP bool, provideNAT bool,
	name string) (Network, error) {
	request := &createNetworkRequest{
		Netblock:    netblock,
		ProvideDHCP: provideDHCP,
		ProvideNAT:  provideNAT,
		Name:        name,
	}
	post, err := json.Marshal(request)
	if err != nil {
		return Network{}, err
	}

	network := Network{}
	err = c.doRequest("networks", "POST", *bytes.NewBuffer(post), &network)
	return network, err
}

// DeleteNetwork removes a network with a specified UUID
func (c *Client) DeleteNetwork(uuid string) error {
	path := "networks/" + uuid
	_, err := c.httpRequest(path, "DELETE", bytes.Buffer{})
	return err
}

// DiskSpec is a definition of an instance disk
type DiskSpec struct {
	Base string `json:"base"`
	Size int    `json:"size"`
	Bus  string `json:"bus"`
	Type string `json:"type"`
}

// NetworkSpec is a definition of an instance network connect
type NetworkSpec struct {
	NetworkUUID string `json:"network_uuid"`
	Address     string `json:"address"`
	MACAddress  string `json:"macaddress"`
	Model       string `json:"model"`
}

// Instance is a definition of an instance
type Instance struct {
	UUID         string                 `json:"uuid"`
	Name         string                 `json:"name"`
	CPUs         int                    `json:"cpus"`
	Memory       int                    `json:"memory"`
	DiskSpecs    []DiskSpec             `json:"disk_spec"`
	SSHKey       string                 `json:"ssh_key"`
	Node         string                 `json:"node"`
	ConsolePort  int                    `json:"console_port"`
	VDIPort      int                    `json:"vdi_port"`
	UserData     string                 `json:"User_data"`
	BlockDevices map[string]interface{} `json:"block_devices"`
	State        string                 `json:"state"`
	StateUpdated float64                `json:"state_updated"`
}

// GetInstances fetches a list of instances
func (c *Client) GetInstances() ([]Instance, error) {
	instances := []Instance{}
	err := c.doRequest("instances", "GET", bytes.Buffer{}, &instances)

	return instances, err
}

// GetInstance fetches a specific instance by UUID
func (c *Client) GetInstance(uuid string) (Instance, error) {
	instance := Instance{}
	err := c.doRequest("instances/"+uuid, "GET", bytes.Buffer{}, &instance)

	return instance, err
}

// NetworkInterface is a definition of an network interface for an instance
type NetworkInterface struct {
	UUID         string  `json:"uuid"`
	NetworkUUID  string  `json:"network_uuid"`
	InstanceUUID string  `json:"instance_uuid"`
	MACAddress   string  `json:"macaddr"`
	IPv4         string  `json:"ipv4"`
	Order        int     `json:"order"`
	Floating     string  `json:"floating"`
	State        string  `json:"state"`
	StateUpdated float64 `json:"state_updated"`
	Model        string  `json:"model"`
}

// GetInstanceInterfaces fetches a list of network interfaces for an instance
func (c *Client) GetInstanceInterfaces(uuid string) ([]NetworkInterface, error) {
	path := "instances/" + uuid + "/interfaces"
	interfaces := []NetworkInterface{}
	err := c.doRequest(path, "GET", bytes.Buffer{}, &interfaces)

	return interfaces, err
}

type createInstanceRequest struct {
	Name     string        `json:"name"`
	CPUs     int           `json:"cpus"`
	Memory   int           `json:"memory"`
	Network  []NetworkSpec `json:"network"`
	Disk     []DiskSpec    `json:"disk"`
	SSHKey   string        `json:"ssh_key"`
	UserData string        `json:"user_data"`
}

// CreateInstance creates a new instance
func (c *Client) CreateInstance(Name string, CPUs int, Memory int,
	Networks []NetworkSpec, Disks []DiskSpec, SSHKey string,
	UserData string) (Instance, error) {
	request := &createInstanceRequest{
		Name:     Name,
		CPUs:     CPUs,
		Memory:   Memory,
		Network:  Networks,
		Disk:     Disks,
		SSHKey:   SSHKey,
		UserData: UserData,
	}
	post, err := json.Marshal(request)
	if err != nil {
		return Instance{}, err
	}

	instance := Instance{}
	err = c.doRequest("instances", "POST", *bytes.NewBuffer(post), &instance)

	return instance, nil
}

// snapshotRequest defines options when making a snapshot of an instance
type snapshotRequest struct {
	All bool `json:"all"`
}

// SnapshotInstance takes a snapshot of an instance
func (c *Client) SnapshotInstance(uuid string, all bool) error {
	path := "instances/" + uuid + "/snapshot"

	request := &snapshotRequest{
		All: all,
	}
	post, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.httpRequest(path, "POST", *bytes.NewBuffer(post))

	return err
}

// Snapshot defines a snapshot of an instance
type Snapshot struct {
	UUID    string `json:"uuid"`
	Device  string `json:"device"`
	Created int64  `json:"created"`
}

// GetInstanceSnapshots fetches a list of instance snapshots
func (c *Client) GetInstanceSnapshots(uuid string) ([]Snapshot, error) {
	snapshots := []Snapshot{}
	path := "instances/" + uuid + "/snapshot"
	err := c.doRequest(path, "GET", bytes.Buffer{}, &snapshots)

	return snapshots, err
}

// RebootInstance reboots an instance
func (c *Client) RebootInstance(uuid string) error {
	return c.postRequest("instances", uuid, "reboot")
}

// PowerOffInstance powers on an instance
func (c *Client) PowerOffInstance(uuid string) error {
	return c.postRequest("instances", uuid, "poweroff")
}

// PowerOnInstance powers on an instance
func (c *Client) PowerOnInstance(uuid string) error {
	return c.postRequest("instances", uuid, "poweron")
}

// PauseInstance will pause an instance
func (c *Client) PauseInstance(uuid string) error {
	return c.postRequest("instances", uuid, "pause")
}

// UnPauseInstance will unpause an instance
func (c *Client) UnPauseInstance(uuid string) error {
	return c.postRequest("instances", uuid, "unpause")
}

// DeleteInstance deletes an instance
func (c *Client) DeleteInstance(uuid string) error {
	_, err := c.httpRequest("instances/"+uuid, "DELETE", bytes.Buffer{})
	return err
}

// Event defines an event that occurred on an instance
type Event struct {
	Timestamp float32 `json:"timestamp"`
	FQDN      string  `json:"fqdn"`
	Operation string  `json:"operation"`
	Phase     string  `json:"phase"`
	Duration  int     `json:"duration"`
	Message   string  `json:"message"`
}

// GetInstanceEvents fetches events that have occurred on a specific instance
func (c *Client) GetInstanceEvents(uuid string) ([]Event, error) {
	events := []Event{}
	err := c.doRequest("instances/"+uuid+"/events", "GET", bytes.Buffer{}, &events)
	return events, err
}

// ImageRequest defines a link to an image
type imageRequest struct {
	URL string `json:"url"`
}

// CacheImage will cache an image
func (c *Client) CacheImage(imageURL string) error {
	request := &imageRequest{
		URL: imageURL,
	}
	post, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.httpRequest("images", "POST", *bytes.NewBuffer(post))

	return err
}

// GetNetworkEvents fetches events that have occurred on a specific network
func (c *Client) GetNetworkEvents(uuid string) ([]Event, error) {
	events := []Event{}
	err := c.doRequest("network/"+uuid+"/events", "GET", bytes.Buffer{}, &events)
	return events, err
}

// Node defines a ShakenFist node
type Node struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	LastSeen string `json:"lastseen"`
}

// GetNodes fetches a list of nodes
func (c *Client) GetNodes() ([]Node, error) {
	nodes := []Node{}
	err := c.doRequest("nodes", "GET", bytes.Buffer{}, &nodes)
	return nodes, err
}

// FloatInterface adds a floating IP to an interface
func (c *Client) FloatInterface(interfaceUUID string) error {
	return c.postRequest("interfaces", interfaceUUID, "float")
}

// DefloatInterface removes a floating IP from an interface
func (c *Client) DefloatInterface(interfaceUUID string) error {
	return c.postRequest("interfaces", interfaceUUID, "defloat")
}

//
// Internal helper functions
//

func (c *Client) getRequest(
	object, uuid string, cmd string, resp interface{}) error {

	err := c.doRequest("GET", object+"/"+uuid+"/"+cmd, bytes.Buffer{}, resp)
	return err
}

func (c *Client) postRequest(object string, uuid string, cmd string) error {
	_, err := c.httpRequest(object+"/"+uuid+"/"+cmd, "POST", bytes.Buffer{})
	return err
}

func (c *Client) doRequest(
	path, method string, data bytes.Buffer, resp interface{}) error {

	body, err := c.httpRequest(path, method, data)

	// Return on error or if JSON decoding not required
	if err != nil || resp == nil {
		return err
	}

	return json.NewDecoder(body).Decode(resp)
}

func (c *Client) httpRequest(
	path, method string, body bytes.Buffer) (io.ReadCloser, error) {

	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Got a non 200 status code: %v",
				resp.StatusCode)
		}
		return nil, fmt.Errorf("Got a non 200 status code: %v - %s",
			resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%s:%v/%s", c.hostname, c.port, path)
}
