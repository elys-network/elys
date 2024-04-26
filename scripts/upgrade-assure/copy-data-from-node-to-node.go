package main

import (
	"log"
	"os/exec"
)

func copyDataFromNodeToNode(homePath, homePath2 string) {
	// Delete first data folder if exists on homePath2
	args := []string{
		"-rf",
		homePath2 + "/data",
	}

	if err := exec.Command("rm", args...).Run(); err != nil {
		log.Fatalf(ColorRed+"Failed to delete data folder on node 2: %v", err)
	}

	// Copy data from node 1 to node 2
	args = []string{
		"-r",
		homePath + "/data",
		homePath2 + "/data",
	}

	if err := exec.Command("cp", args...).Run(); err != nil {
		log.Fatalf(ColorRed+"Failed to copy data from node 1 to node 2: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(ColorYellow + "Data copied from node 1 to node 2")
}
