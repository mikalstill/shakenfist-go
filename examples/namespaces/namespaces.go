package main

import (
	"fmt"
	"os"

	client "github.com/shakenfist/client-go"
)

func main() {
	if os.Getenv("SHAKENFIST_NAMESPACE") != "system" {
		fmt.Println("ERROR: Only the system namespace can access other namespaces")
		return
	}

	c := client.NewClient(
		os.Getenv("SHAKENFIST_URL"),
		os.Getenv("SHAKENFIST_NAMESPACE"),
		os.Getenv("SHAKENFIST_KEY"),
	)

	fmt.Println("******************************")
	fmt.Println("*** Get list of namespaces ***")
	fmt.Println("******************************")
	namespaces, err := c.GetNamespaces()
	if err != nil {
		fmt.Println("GetNamespaces request error:", err)
		return
	}

	for _, n := range namespaces {
		fmt.Println(n)
	}

	fmt.Println("****************************")
	fmt.Println("*** Create new namespace ***")
	fmt.Println("****************************")
	err = c.CreateNamespace("example-testspace")
	if err != nil {
		fmt.Println("Create new namespace:", err)
		return
	}

	fmt.Println("******************************")
	fmt.Println("*** Get list of namespaces ***")
	fmt.Println("******************************")
	namespaces, err = c.GetNamespaces()
	if err != nil {
		fmt.Println("GetNamespaces request error:", err)
		return
	}

	for _, n := range namespaces {
		fmt.Println(n)
	}

	fmt.Println("**************************")
	fmt.Println("*** Create access keys ***")
	fmt.Println("**************************")
	err = c.CreateNamespaceKey("example-testspace", "testkeyname", "testkey")
	if err != nil {
		fmt.Println("Create first key:", err)
		return
	}

	err = c.CreateNamespaceKey("example-testspace", "key2", "key2long")
	if err != nil {
		fmt.Println("Create second key:", err)
		return
	}

	fmt.Println("*******************************")
	fmt.Println("*** Get list of access keys ***")
	fmt.Println("*******************************")
	keys, err := c.GetNamespaceKeys("example-testspace")
	if err != nil {
		fmt.Println("GetNamespaceKeys request error:", err)
		return
	}

	fmt.Println("Keys:")
	for _, n := range keys {
		fmt.Println("   ", n)
	}

	fmt.Println("****************************************")
	fmt.Println("*** Delete access key from namespace ***")
	fmt.Println("****************************************")
	err = c.DeleteNamespaceKey("example-testspace", "testkeyname")
	if err != nil {
		fmt.Println("Delete key error: ", err)
		return
	}

	fmt.Println("*******************************")
	fmt.Println("*** Get list of access keys ***")
	fmt.Println("*******************************")
	keys, err = c.GetNamespaceKeys("example-testspace")
	if err != nil {
		fmt.Println("GetNamespaceKeys request error:", err)
		return
	}

	fmt.Println("Keys:")
	for _, n := range keys {
		fmt.Println("   ", n)
	}

	fmt.Print("\n\n")
	fmt.Println("*************************************")
	fmt.Println("*** Set metadata on the namespace ***")
	fmt.Println("*************************************")
	fmt.Println("Set home='old-people'")
	err = c.SetNamespaceMetadata("example-testspace", "home", "old-people")
	if err != nil {
		fmt.Println("Error setting metadata 'person': ", err)
		return
	}

	fmt.Println("Set exercise='shakes fist'")
	err = c.SetNamespaceMetadata("example-testspace", "exercise", "shakes fist")
	if err != nil {
		fmt.Println("Error setting metadata 'exercise': ", err)
		return
	}

	fmt.Println("********************************************")
	fmt.Println("*** Retrieve metadata from the namespace ***")
	fmt.Println("********************************************")
	meta, err := c.GetNamespaceMetadata("example-testspace")
	if err != nil {
		fmt.Println("Cannot get metadata: ", err)
		return
	}

	fmt.Println("Metadata:")
	for k, v := range meta {
		fmt.Println("   ", k, "=", v)
	}

	fmt.Println("")

	fmt.Println("****************************************")
	fmt.Println("*** Delete metadata on the namespace ***")
	fmt.Println("****************************************")

	err = c.DeleteNamespaceMetadata("example-testspace", "exercise")
	if err != nil {
		fmt.Println("Error deleting metadata 'exercise': ", err)
		return
	}

	fmt.Println("********************************************")
	fmt.Println("*** Retrieve metadata from the namespace ***")
	fmt.Println("********************************************")
	meta, err = c.GetNamespaceMetadata("example-testspace")
	if err != nil {
		fmt.Println("Cannot get metadata: ", err)
		return
	}

	fmt.Println("Metadata:")
	for k, v := range meta {
		fmt.Println("   ", k, "=", v)
	}

	fmt.Println("************************")
	fmt.Println("*** Delete namespace ***")
	fmt.Println("************************")
	if err = c.DeleteNamespace("example-testspace"); err != nil {
		fmt.Println("Delete namespace error: ", err)
		return
	}

	fmt.Println("******************************")
	fmt.Println("*** Get list of namespaces ***")
	fmt.Println("******************************")
	namespaces, err = c.GetNamespaces()
	if err != nil {
		fmt.Println("GetNamespaces request error:", err)
		return
	}

	for _, n := range namespaces {
		fmt.Println(n)
	}
}
