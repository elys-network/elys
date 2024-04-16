package main

import (
	"log"
	"os/exec"
)

func voteOnUpgradeProposal(cmdPath, name, proposalId, homePath, keyringBackend, chainId, node, broadcastMode string) string {
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
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf(ColorRed+"Command execution failed: %v", err)
	}

	// Parse output to find the transaction hash
	txHash, err := parseTxHash(output)
	if err != nil {
		log.Fatalf(ColorRed+"Failed to parse transaction hash: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(ColorYellow+"Voted on upgrade proposal: %s", proposalId)

	// Return the transaction hash
	return txHash
}
