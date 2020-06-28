package main

import (
	"fmt"
	"os"
	"time"

	client "github.com/shakenfist/client-go"
)

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
	c := client.NewClient(
		"http://localhost",
		13000,
		os.Getenv("SHAKENFIST_NAMESPACE"),
		os.Getenv("SHAKENFIST_KEYNAME"),
		os.Getenv("SHAKENFIST_KEY"),
	)

	fmt.Println("*******************************")
	fmt.Println("*** Get a list of instances ***")
	fmt.Println("*******************************")
	instances, err := c.GetInstances()
	if err != nil {
		fmt.Println("GetInstances request error: ", err)
		return
	}

	for _, instance := range instances {
		printInstance(instance)
	}

	fmt.Println("**********************")
	fmt.Println("*** Make a network ***")
	fmt.Println("**********************")
	createdNetwork, err := c.CreateNetwork("192.168.50.0/24", true, true, "golang")
	if err != nil {
		fmt.Println("CreateNetwork request error: ", err)
		return
	}
	networkUUID := createdNetwork.UUID

	fmt.Println("**************************")
	fmt.Println("*** Create an instance ***")
	fmt.Println("**************************")
	instance, err := c.CreateInstance("golang", 1, 1,
		[]client.NetworkSpec{{NetworkUUID: networkUUID}},
		[]client.DiskSpec{{Base: "cirros", Size: 8, Type: "disk", Bus: ""}},
		"", "")
	if err != nil {
		fmt.Println("CreateInstance request error: ", err)
		return
	}
	printInstance(instance)

	fmt.Println("*******************************")
	fmt.Println("*** Get a specific instance ***")
	fmt.Println("*******************************")
	instance, err = c.GetInstance(instance.UUID)
	if err != nil {
		fmt.Println("GetInstance request error: ", err)
		return
	}
	printInstance(instance)

	fmt.Println("**********************************************")
	fmt.Println("*** Get interfaces for a specific instance ***")
	fmt.Println("**********************************************")
	interfaces, err := c.GetInstanceInterfaces(instance.UUID)
	if err != nil {
		fmt.Println("GetInstanceInterfaces request error: ", err)
		return
	}
	printInterfaces(interfaces)

	fmt.Println("**************************")
	fmt.Println("*** Float the instance ***")
	fmt.Println("**************************")
	err = c.FloatInterface(interfaces[0].UUID)
	if err != nil {
		fmt.Println("FloatInterface request error: ", err)
		return
	}

	interfaces, err = c.GetInstanceInterfaces(instance.UUID)
	if err != nil {
		fmt.Println("GetInstanceInterfaces request error: ", err)
		return
	}
	fmt.Println("Interfaces:")
	printInterfaces(interfaces)

	fmt.Println("****************************")
	fmt.Println("*** Defloat the instance ***")
	fmt.Println("****************************")
	err = c.DefloatInterface(interfaces[0].UUID)
	if err != nil {
		fmt.Println("DefloatInterface request error: ", err)
		return
	}

	interfaces, err = c.GetInstanceInterfaces(instance.UUID)
	if err != nil {
		fmt.Println("GetInstanceInterfaces request error: ", err)
		return
	}
	fmt.Println("Interfaces:")
	printInterfaces(interfaces)

	fmt.Println("***************************")
	fmt.Println("*** Delete the instance ***")
	fmt.Println("***************************")
	err = c.DeleteInstance(instance.UUID)
	if err != nil {
		fmt.Println("DeleteInstance request error: ", err)
		return
	}

	fmt.Println("**************************")
	fmt.Println("*** Delete the network ***")
	fmt.Println("**************************")
	err = c.DeleteNetwork(networkUUID)
	if err != nil {
		fmt.Println("DeleteNetwork request error: ", err)
		return
	}
}
