package client

import (
	"bytes"
	"encoding/json"
)

// Network is a definition of a network.
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

// GetNetworks fetches a list of networks.
func (c *Client) GetNetworks() ([]Network, error) {
	networks := []Network{}
	err := c.doRequest("networks", "GET", bytes.Buffer{}, &networks)
	return networks, err
}

// GetNetwork fetches a specific instance by UUID.
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

// CreateNetwork creates a new network.
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

// DeleteNetwork removes a network with a specified UUID.
func (c *Client) DeleteNetwork(uuid string) error {
	path := "networks/" + uuid
	err := c.doRequest(path, "DELETE", bytes.Buffer{}, nil)
	return err
}

// NetworkSpec is a definition of an instance network connect.
type NetworkSpec struct {
	NetworkUUID string `json:"network_uuid"`
	Address     string `json:"address"`
	MACAddress  string `json:"macaddress"`
	Model       string `json:"model"`
}

// NetworkInterface is a definition of an network interface for an instance.
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

// GetInstanceInterfaces fetches a list of network interfaces for an instance.
func (c *Client) GetInstanceInterfaces(uuid string) ([]NetworkInterface, error) {
	path := "instances/" + uuid + "/interfaces"
	interfaces := []NetworkInterface{}
	err := c.doRequest(path, "GET", bytes.Buffer{}, &interfaces)

	return interfaces, err
}

// GetInterface fetches a specific network interface.
func (c *Client) GetInterface(uuid string) (NetworkInterface, error) {
	path := "interfaces/" + uuid
	iface := NetworkInterface{}
	err := c.doRequest(path, "GET", bytes.Buffer{}, &iface)

	return iface, err
}

// FloatInterface adds a floating IP to an interface.
func (c *Client) FloatInterface(interfaceUUID string) error {
	return c.postRequest("interfaces", interfaceUUID, "float")
}

// DefloatInterface removes a floating IP from an interface.
func (c *Client) DefloatInterface(interfaceUUID string) error {
	return c.postRequest("interfaces", interfaceUUID, "defloat")
}

// GetNetworkEvents fetches events that have occurred on a specific network.
func (c *Client) GetNetworkEvents(uuid string) ([]Event, error) {
	events := []Event{}
	err := c.doRequest("network/"+uuid+"/events", "GET", bytes.Buffer{}, &events)
	return events, err
}
