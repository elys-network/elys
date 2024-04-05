package main

import (
	"log"
	"os/exec"
)

func createValidator(cmdPath, name, selfDelegation, moniker, pubkey, homePath, keyringBackend, chainId, node, broadcastMode string) {
	// Command and arguments
	args := []string{
		"tx",
		"staking",
		"create-validator",
		"--amount", "1000000uelys",
		"--pubkey", pubkey,
		"--moniker", moniker,
		"--commission-rate", "0.05",
		"--commission-max-rate", "0.50",
		"--commission-max-change-rate", "0.01",
		"--min-self-delegation", selfDelegation,
		"--from", name,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", node,
		"--broadcast-mode", broadcastMode,
		"--fees", "100000uelys",
		"--gas", "3000000",
		"--gas-adjustment", "1.5",
		"--home", homePath,
		"--yes",
	}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Failed to create validator: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"Validator %s created successfully", moniker)
}
