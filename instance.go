package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DiskSpec is a definition of an instance disk.
type DiskSpec struct {
	Base string `json:"base"`
	Size int    `json:"size"`
	Bus  string `json:"bus"`
	Type string `json:"type"`
}

// Instance is a definition of an instance.
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
	UserData     string                 `json:"user_data"`
	BlockDevices map[string]interface{} `json:"block_devices"`
	State        string                 `json:"state"`
	StateUpdated float64                `json:"state_updated"`
}

// GetInstances fetches a list of instances.
func (c *Client) GetInstances() ([]Instance, error) {
	instances := []Instance{}
	err := c.doRequestJSON("instances", "GET", bytes.Buffer{}, &instances)

	return instances, err
}

// GetInstance fetches a specific instance by UUID.
func (c *Client) GetInstance(uuid string) (Instance, error) {
	instance := Instance{}
	err := c.doRequestJSON("instances/"+uuid, "GET", bytes.Buffer{}, &instance)

	return instance, err
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

// CreateInstance creates a new instance.
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
	err = c.doRequestJSON("instances", "POST", *bytes.NewBuffer(post), &instance)

	return instance, err
}

// snapshotRequest defines options when making a snapshot of an instance.
type snapshotRequest struct {
	All bool `json:"all"`
}

// SnapshotInstance takes a snapshot of an instance.
func (c *Client) SnapshotInstance(uuid string, all bool) error {
	path := "instances/" + uuid + "/snapshot"

	request := &snapshotRequest{
		All: all,
	}
	post, err := json.Marshal(request)
	if err != nil {
		return err
	}

	err = c.doRequestJSON(path, "POST", *bytes.NewBuffer(post), nil)

	return err
}

// Snapshot defines a snapshot of an instance.
type Snapshot struct {
	UUID    string `json:"uuid"`
	Device  string `json:"device"`
	Created int64  `json:"created"`
}

// GetInstanceSnapshots fetches a list of instance snapshots.
func (c *Client) GetInstanceSnapshots(uuid string) ([]Snapshot, error) {
	snapshots := []Snapshot{}
	path := "instances/" + uuid + "/snapshot"
	err := c.doRequestJSON(path, "GET", bytes.Buffer{}, &snapshots)

	return snapshots, err
}

// RebootInstance reboots an instance.
func (c *Client) RebootInstance(uuid string) error {
	return c.postRequest("instances", uuid, "reboot")
}

// PowerOffInstance powers on an instance.
func (c *Client) PowerOffInstance(uuid string) error {
	return c.postRequest("instances", uuid, "poweroff")
}

// PowerOnInstance powers on an instance.
func (c *Client) PowerOnInstance(uuid string) error {
	return c.postRequest("instances", uuid, "poweron")
}

// PauseInstance will pause an instance.
func (c *Client) PauseInstance(uuid string) error {
	return c.postRequest("instances", uuid, "pause")
}

// UnPauseInstance will unpause an instance.
func (c *Client) UnPauseInstance(uuid string) error {
	return c.postRequest("instances", uuid, "unpause")
}

// DeleteInstance deletes an instance.
func (c *Client) DeleteInstance(uuid string) error {
	err := c.doRequestJSON("instances/"+uuid, "DELETE", bytes.Buffer{}, nil)
	return err
}

// Event defines an event that occurred on an instance.
type Event struct {
	Timestamp float32 `json:"timestamp"`
	FQDN      string  `json:"fqdn"`
	Operation string  `json:"operation"`
	Phase     string  `json:"phase"`
	Duration  float32 `json:"duration"`
	Message   string  `json:"message"`
}

// GetInstanceEvents fetches events that have occurred on a specific instance.
func (c *Client) GetInstanceEvents(uuid string) ([]Event, error) {
	events := []Event{}
	err := c.getRequest("instances", uuid, "events", &events)
	return events, err
}

// ImageRequest defines a link to an image.
type imageRequest struct {
	URL string `json:"url"`
}

// CacheImage will cache an image.
func (c *Client) CacheImage(imageURL string) error {
	request := &imageRequest{
		URL: imageURL,
	}
	post, err := json.Marshal(request)
	if err != nil {
		return err
	}

	err = c.doRequestJSON("images", "POST", *bytes.NewBuffer(post), nil)

	return err
}

type consoleDataReq struct {
	Length int `json:"length"`
}

// GetConsoleData retrieves the last n bytes of console data from an instance.
func (c *Client) GetConsoleData(uuid string, n int) (string, error) {
	path := "instances/" + uuid + "/consoledata"

	req := &consoleDataReq{
		Length: n,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("cannot marshal consoledata request: %v", err)
	}

	resp, err := c.doRequest(path, "GET", *bytes.NewBuffer(reqData))
	if err != nil {
		return "", fmt.Errorf("cannot retrieve console data: %v", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp)
	if err != nil {
		return "", fmt.Errorf("cannot read http response buffer: %v", err)
	}
	d := buf.String()

	return d, nil
}
