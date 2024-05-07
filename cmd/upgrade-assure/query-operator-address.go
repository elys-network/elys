package main

import (
	"log"
	"os/exec"
	"strings"
)

func queryOperatorAddress(cmdPath, homePath, keyringBackend, validatorName string) string {
	// Command and arguments
	args := []string{"keys", "show", validatorName, "--bech", "val", "--home", homePath, "--keyring-backend", keyringBackend, "--address"}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf(ColorRed+"Failed to query validator pubkey: %v", err)
	}

	// trim the output
	outputStr := strings.TrimSpace(string(output))

	return outputStr
}
