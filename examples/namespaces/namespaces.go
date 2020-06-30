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
		"http://localhost",
		13000,
		os.Getenv("SHAKENFIST_NAMESPACE"),
		os.Getenv("SHAKENFIST_KEY"),
	)

	fmt.Println("******************************")
	fmt.Println("*** Get list of namespaces ***")
	fmt.Println("******************************")
	namespaces, err := c.GetNameSpaces()
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
	err = c.CreateNameSpace("example-testspace", "testkeyname", "testkey")
	if err != nil {
		fmt.Println("Create new namespace:", err)
		return
	}

	fmt.Println("******************************")
	fmt.Println("*** Get list of namespaces ***")
	fmt.Println("******************************")
	namespaces, err = c.GetNameSpaces()
	if err != nil {
		fmt.Println("GetNamespaces request error:", err)
		return
	}

	for _, n := range namespaces {
		fmt.Println(n)
	}

	fmt.Println("*******************************")
	fmt.Println("*** Delete key in namespace ***")
	fmt.Println("*******************************")
	err = c.DeleteNameSpaceKey("example-testspace", "testkeyname")
	if err != nil {
		fmt.Println("Delete key error: ", err)
		return
	}

	fmt.Println("************************")
	fmt.Println("*** Delete namespace ***")
	fmt.Println("************************")
	if err = c.DeleteNameSpace("example-testspace"); err != nil {
		fmt.Println("Delete namespace error: ", err)
		return
	}

	fmt.Println("******************************")
	fmt.Println("*** Get list of namespaces ***")
	fmt.Println("******************************")
	namespaces, err = c.GetNameSpaces()
	if err != nil {
		fmt.Println("GetNamespaces request error:", err)
		return
	}

	for _, n := range namespaces {
		fmt.Println(n)
	}
}
