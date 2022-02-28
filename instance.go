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

// VideoSpec defines the type of video card in an instance.
type VideoSpec struct {
	Model  string `json:"model"`
	Memory int    `json:"memory"` // Memory size in KB
}

// Instance is a definition of an instance.
type Instance struct {
	BlockDevices      map[string]interface{} `json:"block_devices"`
	ConsolePort       int                    `json:"console_port"`
	CPUs              int                    `json:"cpus"`
	DiskSpecs         []DiskSpec             `json:"disk_spec"`
	Metadata          string                 `json:"metadata"`
	Memory            int                    `json:"memory"`
	Name              string                 `json:"name"`
	Namespace         string                 `json:"namespace"`
	NetworkInterfaces []NetworkInterface     `json:"network_interfaces"`
	Node              string                 `json:"node"`
	PowerState        string                 `json:"power_state"`
	SSHKey            string                 `json:"ssh_key"`
	State             string                 `json:"state"`
	StateUpdated      float64                `json:"state_updated"`
	SecureBoot        bool                   `json:"secure_boot"`
	UEFI              bool                   `json:"uefi"`
	UserData          string                 `json:"user_data"`
	UUID              string                 `json:"uuid"`
	VDIPort           int                    `json:"vdi_port"`
	Video             VideoSpec              `json:"video"`
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
	Name          string        `json:"name"`
	CPUs          int           `json:"cpus"`
	Memory        int           `json:"memory"`
	Metadata      string        `json:"metadata"`
	NameSpace     string        `json:"namespace"`
	Network       []NetworkSpec `json:"network"`
	NVRAMTemplate string        `json:"nvram_template"`
	Disk          []DiskSpec    `json:"disk"`
	Video         VideoSpec     `json:"video"`
	SecureBoot    bool          `json:"secure_boot"`
	SSHKey        string        `json:"ssh_key"`
	UEFI          bool          `json:"uefi"`
	UserData      string        `json:"user_data"`
}

// CreateInstance creates a new instance.
func (c *Client) CreateInstance(name string, cpus int, memory int,
	networks []NetworkSpec, disks []DiskSpec, video VideoSpec, sshKey string,
	userData string, nameSpace string, metadata string, secureBoot bool,
	uefi bool, nvramTemplate string) (Instance, error) {

	request := &createInstanceRequest{
		Name:          name,
		CPUs:          cpus,
		Memory:        memory,
		Metadata:      metadata,
		NameSpace:     nameSpace,
		Network:       networks,
		NVRAMTemplate: nvramTemplate,
		Disk:          disks,
		Video:         video,
		SecureBoot:    secureBoot,
		SSHKey:        sshKey,
		UEFI:          uefi,
		UserData:      userData,
	}
	post, err := json.Marshal(request)
	if err != nil {
		return Instance{}, err
	}

	instance := Instance{}
	err = c.doRequestJSON("instances", "POST", *bytes.NewBuffer(post), &instance)

	return instance, err
}

// SnapshotInstance takes a snapshot of an instance.
func (c *Client) SnapshotInstance(uuid string, all bool, device string) error {

	path := "instances/" + uuid + "/snapshot"

	request := &struct {
		All    bool   `json:"all"`
		Device string `json:"device"`
	}{
		All:    all,
		Device: device,
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
func (c *Client) DeleteInstance(uuid string, namespace string) error {
	var err error
	var req []byte

	if namespace != "" {
		n := &struct {
			Namespace string `json:"namespace"`
		}{
			Namespace: namespace,
		}
		req, err = json.Marshal(n)
		if err != nil {
			return fmt.Errorf("Unable to marshal data: %v", err)
		}
	}
	err = c.doRequestJSON("instances/"+uuid, "DELETE", *bytes.NewBuffer(req), nil)
	return err
}

type deleteAllRequest struct {
	Namespace string `json:"namespace"`
	Confirm   bool   `json:"confirm"`
}

// DeleteAllInstances deletes all instances within a namespace. Specifying
// namespace "system" will delete all instances in a cluster.
func (c *Client) DeleteAllInstances(namespace string) ([]string, error) {
	instances := []string{}

	n := deleteAllRequest{
		Namespace: namespace,
		Confirm:   true,
	}
	req, err := json.Marshal(n)
	if err != nil {
		return instances, fmt.Errorf("Unable to marshal data: %v", err)
	}

	err = c.doRequestJSON("instances",
		"DELETE", *bytes.NewBuffer(req), &instances)

	return instances, err
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

// SetInstanceMetadataItem sets a metadata key on an instance to value.
func (c *Client) SetInstanceMetadataItem(uuid string, key string, value string) error {
	path := "instances/" + uuid + "/metadata/" + key

	request := &struct {
		Value string `json:"value"`
	}{
		Value: value,
	}

	put, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return c.doRequestJSON(path, "PUT", *bytes.NewBuffer(put), nil)
}

// DeleteInstanceMetadataItem deletes an individual metadata key on an instance.
func (c *Client) DeleteInstanceMetadataItem(uuid string, key string) error {
	path := "instances/" + uuid + "/metadata/" + key
	return c.doRequestJSON(path, "DELETE", bytes.Buffer{}, nil)
}

// UpdateLabel changes the name of a blob label
func (c *Client) UpdateLabel(labelName string, blobUUID string) error {
	path := "label/" + labelName

	request := &struct {
		BlobUUID string `json:"blob_uuid"`
	}{
		BlobUUID: blobUUID,
	}

	put, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return c.doRequestJSON(path, "PUT", *bytes.NewBuffer(put), nil)
}
