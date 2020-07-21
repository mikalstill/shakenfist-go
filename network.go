package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Generic structure to send a UUID.
type reqUUID struct {
	UUID string `json:"uuid"`
}

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
	err := c.doRequestJSON("networks", "GET", bytes.Buffer{}, &networks)
	return networks, err
}

// GetNetwork fetches a specific instance by UUID.
func (c *Client) GetNetwork(uuid string) (Network, error) {
	network := Network{}
	err := c.doRequestJSON("networks/"+uuid, "GET", bytes.Buffer{}, &network)
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
	err = c.doRequestJSON("networks", "POST", *bytes.NewBuffer(post), &network)
	return network, err
}

// DeleteNetwork removes a network with a specified UUID.
func (c *Client) DeleteNetwork(uuid string) error {
	path := "networks/" + uuid
	err := c.doRequestJSON(path, "DELETE", bytes.Buffer{}, nil)
	return err
}

// DeleteAllNetworks deletes all networks within a namespace. Specifying
// namespace of "system" will attempt to delete all networks in a cluster.
func (c *Client) DeleteAllNetworks(namespace string) ([]string, error) {
	networks := []string{}

	n := deleteAllRequest{
		Namespace: namespace,
		Confirm:   true,
	}
	req, err := json.Marshal(n)
	if err != nil {
		return networks, fmt.Errorf("Unable to marshal data: %v", err)
	}

	err = c.doRequestJSON("networks", "DELETE", *bytes.NewBuffer(req), &networks)

	return networks, err
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
	err := c.doRequestJSON(path, "GET", bytes.Buffer{}, &interfaces)

	return interfaces, err
}

// GetInterface fetches a specific network interface.
func (c *Client) GetInterface(uuid string) (NetworkInterface, error) {
	path := "interfaces/" + uuid
	iface := NetworkInterface{}
	err := c.doRequestJSON(path, "GET", bytes.Buffer{}, &iface)

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
	err := c.getRequest("networks", uuid, "events", &events)
	return events, err
}

// DeployNetworkNode sends a DeployNetworkNode command to Shaken Fist.
func (c *Client) DeployNetworkNode(uuid string) error {
	req := &reqUUID{
		UUID: uuid,
	}

	put, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}

	err = c.doRequestJSON("deploy_network_node", "PUT", *bytes.NewBuffer(put), nil)
	return err
}

// UpdateDHCP sends an UpdateDHCP command to Shaken Fist.
func (c *Client) UpdateDHCP(uuid string) error {
	req := &reqUUID{
		UUID: uuid,
	}

	put, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}

	err = c.doRequestJSON("update_dhcp", "PUT", *bytes.NewBuffer(put), nil)
	return err
}

// RemoveDHCP sends an RemoveDHCP command to Shaken Fist.
func (c *Client) RemoveDHCP(uuid string) error {
	req := &reqUUID{
		UUID: uuid,
	}

	put, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}

	err = c.doRequestJSON("remove_dhcp", "PUT", *bytes.NewBuffer(put), nil)
	return err
}
