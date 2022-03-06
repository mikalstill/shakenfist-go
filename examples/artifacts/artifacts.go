package main

import (
	"fmt"
	"os"

	client "github.com/shakenfist/client-go"
)

func printBlob(blob client.Blob) {
	fmt.Printf("       UUID: %s\n", blob.UUID)
	fmt.Printf("       Instances: %s\n", blob.Instances)
	fmt.Printf("       Size: %d\n", blob.Size)
	fmt.Printf("       ReferenceCount: %d\n", blob.ReferenceCount)
	fmt.Printf("       DependsOn: %s\n", blob.DependsOn)
}

func printArtifact(artifact client.Artifact) {
	fmt.Printf("UUID: %s\n", artifact.UUID)
	fmt.Printf("Type: %s\n", artifact.Type)
	fmt.Printf("State: %s\n", artifact.State)
	fmt.Printf("SourceURL: %s\n", artifact.SourceURL)
	fmt.Printf("Version: %d\n", artifact.Version)
	fmt.Printf("MaxVersions: %d\n", artifact.MaxVersions)
	fmt.Printf("MostRecentIndex: %d\n", artifact.MostRecentIndex)
	fmt.Println("Blobs:")
	for _, b := range artifact.Blobs {
		printBlob(b)
		println()
	}
}

func printEvent(event client.Event) {
	fmt.Printf("Timestamp: %f\n", event.Timestamp)
	fmt.Printf("FQDN: %s\n", event.FQDN)
	fmt.Printf("Operation: %s\n", event.Operation)
	fmt.Printf("Phase: %s\n", event.Phase)
	fmt.Printf("Duration: %f\n", event.Duration)
	fmt.Printf("Message: %s\n", event.Message)
}

func main() {
	c := client.NewClient(
		os.Getenv("SHAKENFIST_API_URL"),
		os.Getenv("SHAKENFIST_NAMESPACE"),
		os.Getenv("SHAKENFIST_KEY"),
	)

	fmt.Println("*********************")
	fmt.Println("*** Get Artifacts ***")
	fmt.Println("*********************")
	artifacts, err := c.GetArtifacts("")
	if err != nil {
		fmt.Println("GetArtifacts request error: ", err)
		return
	}
	for _, a := range artifacts {
		printArtifact(a)
		println()
	}

	fmt.Println("*****************")
	fmt.Println("*** Get Blobs ***")
	fmt.Println("*****************")
	blobs, err := c.GetBlobs("")
	if err != nil {
		fmt.Println("GetBlobs request error: ", err)
		return
	}
	for _, b := range blobs {
		printBlob(b)
		println()
	}

	fmt.Println("******************")
	fmt.Println("*** Get Events ***")
	fmt.Println("******************")
	events, err := c.GetArtifactEvents(artifacts[0].UUID)
	if err != nil {
		fmt.Println("GetArtifactEvents request error: ", err)
		return
	}
	for _, e := range events {
		printEvent(e)
		println()
	}

	fmt.Println("*****************************")
	fmt.Println("*** Get Artifact Versions ***")
	fmt.Println("*****************************")
	blobs, err = c.GetArtifactVersions(artifacts[0].UUID)
	if err != nil {
		fmt.Println("GetArtifactVersions request error: ", err)
		return
	}
	for _, b := range blobs {
		printBlob(b)
		println()
	}
}
