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
