package main

import (
	"fmt"
	"os/exec"
)

func checkTxStatus(cmdPath, node, txHash string) (bool, error) {
	args := []string{"q", "tx", txHash, "--node", node, "--output", "json"}
	_, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to query tx status: %w", err)
	}
	return true, nil
}
