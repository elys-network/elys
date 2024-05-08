package main

import (
	"log"
	"os/exec"
	"time"
)

func unbondValidator(cmdPath, validatorKeyName, operatorAddress, validatorSelfDelegation, keyringBackend, chainId, rpc, broadcastMode, homePath string) {
	// Command and arguments
	args := []string{
		"tx",
		"staking",
		"unbond",
		operatorAddress,
		validatorSelfDelegation + "uelys",
		"--from", validatorKeyName,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", rpc,
		"--broadcast-mode", broadcastMode,
		"--fees", "100000uelys",
		"--gas", "1000000",
		"--home", homePath,
		"--output", "json",
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
	log.Printf(ColorYellow+"Unbonded validator: %s, self-delegation: %s", operatorAddress, validatorSelfDelegation)

	waitForTxConfirmation(cmdPath, rpc, txHash, 5*time.Minute)
}
