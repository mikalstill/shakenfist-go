package main

import (
	"fmt"
	"os"
	"strconv"
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

func main() {
	port, ok := strconv.Atoi(os.Getenv("SHAKENFIST_PORT"))
	if ok != nil {
		port = 13000
	}

	c := client.NewClient(
		os.Getenv("SHAKENFIST_HOSTNAME"),
		port,
		os.Getenv("SHAKENFIST_NAMESPACE"),
		os.Getenv("SHAKENFIST_KEY"),
	)

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
	fmt.Printf("Requesting network: %v\n", createdNetwork.UUID)
	network, err := c.GetNetwork(createdNetwork.UUID)
	if err != nil {
		fmt.Println("GetNetwork request error: ", err)
		return
	}
	printNetwork(network)

	fmt.Print("\n\n")
	fmt.Println("***********************************")
	fmt.Println("*** Set metadata on the network ***")
	fmt.Println("***********************************")
	fmt.Println("Set nets='old-people'")
	err = c.SetNetworkMetadata(createdNetwork.UUID, "nets", "old people")
	if err != nil {
		fmt.Println("Error setting metadata 'nets': ", err)
		return
	}

	fmt.Println("Set action='shakes fist'")
	err = c.SetNetworkMetadata(createdNetwork.UUID, "action", "shakes fist")
	if err != nil {
		fmt.Println("Error setting metadata 'action': ", err)
		return
	}

	fmt.Println("******************************************")
	fmt.Println("*** Retrieve metadata from the network ***")
	fmt.Println("******************************************")
	meta, err := c.GetNetworkMetadata(createdNetwork.UUID)
	if err != nil {
		fmt.Println("Cannot get metadata: ", err)
		return
	}

	fmt.Println("Metadata:")
	for k, v := range meta {
		fmt.Println("   ", k, "=", v)
	}

	fmt.Println("")

	fmt.Println("**************************************")
	fmt.Println("*** Delete metadata on the network ***")
	fmt.Println("**************************************")

	err = c.DeleteNetworkMetadata(createdNetwork.UUID, "action")
	if err != nil {
		fmt.Println("Error deleting metadata 'action': ", err)
		return
	}

	fmt.Println("******************************************")
	fmt.Println("*** Retrieve metadata from the network ***")
	fmt.Println("******************************************")
	meta, err = c.GetNetworkMetadata(createdNetwork.UUID)
	if err != nil {
		fmt.Println("Cannot get metadata: ", err)
		return
	}

	fmt.Println("Metadata:")
	for k, v := range meta {
		fmt.Println("   ", k, "=", v)
	}

	fmt.Println("************************")
	fmt.Println("*** Delete a network ***")
	fmt.Println("************************")
	err = c.DeleteNetwork(createdNetwork.UUID)
	if err != nil {
		fmt.Println("DeleteNetwork request error: ", err)
		return
	}
}
