package main

import (
	"log"
	"os/exec"
)

func queryUpgradeApplied(cmdPath, node, newVersion string) {
	// Command and arguments
	args := []string{"q", "upgrade", "applied", newVersion, "--node", node}

	// Execute the command
	err := exec.Command(cmdPath, args...).Run()
	if err != nil {
		log.Fatalf("Failed to retrieve applied upgrade: %v", err)
	}
}
