package main

import (
	"fmt"
	"os"

	client "github.com/shakenfist/client-go"
)

func main() {
	c := client.NewClient(
		os.Getenv("SHAKENFIST_API_URL"),
		os.Getenv("SHAKENFIST_NAMESPACE"),
		os.Getenv("SHAKENFIST_KEY"),
	)

	fmt.Println("**************************")
	fmt.Println("*** Create an instance ***")
	fmt.Println("**************************")
	instance, err := c.CreateInstance("golang", 1, 1,
		[]client.NetworkSpec{},
		[]client.DiskSpec{{Base: "cirros", Size: 8, Type: "disk", Bus: ""}},
		client.VideoSpec{Model: "cirrus", Memory: 16384},
		"", "", "", "", false, false, "")
	if err != nil {
		fmt.Println("CreateInstance request error: ", err)
		return
	}

	fmt.Println("************************************")
	fmt.Println("*** Set metadata on the instance ***")
	fmt.Println("************************************")
	fmt.Println("Set person='old-man'")
	err = c.SetInstanceMetadata(instance.UUID, "person", "old-man")
	if err != nil {
		fmt.Println("Error setting metadata 'person': ", err)
		return
	}

	fmt.Println("Set action='shakes fist'")
	err = c.SetInstanceMetadata(instance.UUID, "action", "shakes fist")
	if err != nil {
		fmt.Println("Error setting metadata 'action': ", err)
		return
	}

	fmt.Println("*******************************************")
	fmt.Println("*** Retrieve metadata from the instance ***")
	fmt.Println("*******************************************")
	meta, err := c.GetInstanceMetadata(instance.UUID)
	if err != nil {
		fmt.Println("Cannot get metadata: ", err)
		return
	}

	fmt.Println("Metadata:")
	for k, v := range meta {
		fmt.Println("   ", k, "=", v)
	}

	fmt.Println("")

	fmt.Println("***************************************")
	fmt.Println("*** Delete metadata on the instance ***")
	fmt.Println("***************************************")

	err = c.DeleteInstanceMetadata(instance.UUID, "action")
	if err != nil {
		fmt.Println("Error deleting metadata 'action': ", err)
		return
	}

	fmt.Println("*******************************************")
	fmt.Println("*** Retrieve metadata from the instance ***")
	fmt.Println("*******************************************")
	meta, err = c.GetInstanceMetadata(instance.UUID)
	if err != nil {
		fmt.Println("Cannot get metadata: ", err)
		return
	}

	fmt.Println("Metadata:")
	for k, v := range meta {
		fmt.Println("   ", k, "=", v)
	}

	fmt.Println("")
	fmt.Println("***************************")
	fmt.Println("*** Delete the instance ***")
	fmt.Println("***************************")
	err = c.DeleteInstance(instance.UUID, "")
	if err != nil {
		fmt.Println("DeleteInstance request error: ", err)
		return
	}
}
