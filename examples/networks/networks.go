package main

import (
	"fmt"
	"time"

	client "github.com/shakenfist/client-go"
)

func printNetwork(network client.Network) {
	fmt.Printf("UUID: %s\n", network.UUID)
	fmt.Printf("Name: %s\n", network.Name)
	fmt.Printf("Net Block: %s\n", network.NetBlock)
	fmt.Printf("Provide DHCP: %t\n", network.ProvideDHCP)
	fmt.Printf("Provide NAT: %t\n", network.ProvideNAT)
	fmt.Printf("Owner: %s\n", network.Owner)
	fmt.Printf("Floating Gateway: %s\n", network.FloatingGateway)
	fmt.Printf("State: %s\n", network.State)
	fmt.Printf("StateUpdated: %s\n", time.Unix(int64(network.StateUpdated), 0))
	fmt.Println("")
}

func printInstance(instance client.Instance) {
	fmt.Printf("UUID: %s\n", instance.UUID)
	fmt.Printf("Name: %s\n", instance.Name)
	fmt.Printf("CPUs: %d\n", instance.CPUs)
	fmt.Printf("Memory (MB): %d\n", instance.Memory)
	fmt.Println("Disks:")
	for _, disk := range instance.DiskSpecs {
		fmt.Printf("  - Base: %s\n", disk.Base)
		fmt.Printf("    Size: %d\n", disk.Size)
		fmt.Printf("    Bus:  %s\n", disk.Bus)
		fmt.Printf("    Type: %s\n", disk.Type)
	}
	fmt.Printf("SSHKey: %s\n", instance.SSHKey)
	fmt.Printf("Node: %s\n", instance.Node)
	fmt.Printf("ConsolePort: %d\n", instance.ConsolePort)
	fmt.Printf("VDIPort: %d\n", instance.VDIPort)
	fmt.Printf("UserData: %s\n", instance.UserData)
	fmt.Printf("State: %s\n", instance.State)
	fmt.Printf("StateUpdated: %s\n", time.Unix(int64(instance.StateUpdated), 0))
	fmt.Println("")
}

func printInterfaces(interfaces []client.NetworkInterface) {
	for _, iface := range interfaces {
		fmt.Printf("  - UUID: %s\n", iface.UUID)
		fmt.Printf("    Network UUID: %s\n", iface.NetworkUUID)
		fmt.Printf("    Instance UUID: %s\n", iface.InstanceUUID)
		fmt.Printf("    MAC Address: %s\n", iface.MACAddress)
		fmt.Printf("    IPv4 Address: %s\n", iface.IPv4)
		fmt.Printf("    Order: %d\n", iface.Order)
		fmt.Printf("    Floating Address: %s\n", iface.Floating)
		fmt.Printf("    State: %s\n", iface.State)
		fmt.Printf("    StateUpdated: %s\n", time.Unix(int64(iface.StateUpdated), 0))
		fmt.Printf("    Model: %s\n", iface.Model)
	}

	fmt.Println("")
}

func main() {
	c := client.NewClient("http://localhost", 13000)

	fmt.Println("******************************")
	fmt.Println("*** Get a list of networks ***")
	fmt.Println("******************************")
	networks, err := c.GetNetworks()
	if err != nil {
		fmt.Println("GetNetworks request error: ", err)
		return
	}

	for _, network := range networks {
		printNetwork(network)
	}

	fmt.Println("**********************")
	fmt.Println("*** Make a network ***")
	fmt.Println("**********************")
	createdNetwork, err := c.CreateNetwork("192.168.50.0/24", true, true, "golang")
	if err != nil {
		fmt.Println("CreateNetwork request error: ", err)
		return
	}
	printNetwork(createdNetwork)

	fmt.Println("******************************")
	fmt.Println("*** Get a list of networks ***")
	fmt.Println("******************************")
	networks, err = c.GetNetworks()
	if err != nil {
		fmt.Println("GetNetworks request error: ", err)
		return
	}

	for _, network := range networks {
		printNetwork(network)
	}

	fmt.Println("******************************")
	fmt.Println("*** Get a specific network ***")
	fmt.Println("******************************")
	fmt.Printf("Requesting network: %v\n", networks[0].UUID)
	network, err := c.GetNetwork(networks[0].UUID)
	if err != nil {
		fmt.Println("GetNetwork request error: ", err)
		return
	}
	printNetwork(network)

	fmt.Println("************************")
	fmt.Println("*** Delete a network ***")
	fmt.Println("************************")
	err = c.DeleteNetwork(createdNetwork.UUID)
	if err != nil {
		fmt.Println("DeleteNetwork request error: ", err)
		return
	}
}
