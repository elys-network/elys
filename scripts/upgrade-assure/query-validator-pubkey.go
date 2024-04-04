package main

import (
	"log"
	"os/exec"
)

func queryValidatorPubkey(cmdPath, home string) string {
	// Command and arguments
	args := []string{"tendermint", "show-validator", "--home", home}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf(Red+"Failed to query validator pubkey: %v", err)
	}

	// convert output to string
	outputStr := string(output)

	return outputStr
}
