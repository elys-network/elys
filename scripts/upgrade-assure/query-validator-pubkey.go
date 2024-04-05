package main

import (
	"log"
	"os/exec"
	"strings"
)

func queryValidatorPubkey(cmdPath, home string) string {
	// Command and arguments
	args := []string{"tendermint", "show-validator", "--home", home}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf(Red+"Failed to query validator pubkey: %v", err)
	}

	// trim the output
	outputStr := strings.TrimSpace(string(output))

	return outputStr
}
