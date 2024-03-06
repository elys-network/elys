package main

import (
	"log"
	"os/exec"
)

func retrieveSnapshot(snapshotUrl, homePath string) {
	var cmdString string

	// if snapshot url ends with `.tar.lz4`
	if snapshotUrl[len(snapshotUrl)-8:] == ".tar.lz4" {
		// Construct the command string
		cmdString = "curl -o - -L " + snapshotUrl + " | lz4 -c -d - | tar -x -C " + homePath
	} else if snapshotUrl[len(snapshotUrl)-7:] == ".tar.gz" {
		// if snapshot url ends with `.tar.gz`
		// Construct the command string
		cmdString = "curl -o - -L " + snapshotUrl + " | tar -xz -C " + homePath
	} else if snapshotUrl[len(snapshotUrl)-4:] == ".tar" {
		// if snapshot url ends with `.tar`
		// Construct the command string
		cmdString = "curl -o - -L " + snapshotUrl + " | tar -x -C " + homePath
	} else {
		// otherwise, the snapshot url is invalid
		log.Fatalf(Red+"Invalid snapshot url: %s", snapshotUrl)
	}

	// print cmdString
	log.Printf(Green+"Retrieving snapshot using command: %s",
		cmdString)

	// Execute the command using /bin/sh
	cmd := exec.Command("/bin/sh", "-c", cmdString)
	if err := cmd.Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"Snapshot retrieved and extracted to path: %s", homePath)
}
