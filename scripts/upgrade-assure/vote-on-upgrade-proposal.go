package main

import (
	"log"
	"os/exec"
)

func voteOnUpgradeProposal(cmdPath, name, proposalId, homePath, keyringBackend, chainId, node, broadcastMode string) {
	// Command and arguments
	args := []string{
		"tx", "gov", "vote", proposalId, "yes",
		"--from", name,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", node,
		"--broadcast-mode", broadcastMode,
		"--fees", "100000uelys",
		"--gas", "1000000",
		"--home", homePath,
		"--yes",
	}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(ColorRed+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(ColorYellow+"Voted on upgrade proposal: %s", proposalId)
}
