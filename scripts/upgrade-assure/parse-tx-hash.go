package main

import (
	"encoding/json"
	"fmt"
)

// TxResponse represents the structure of the transaction output
type TxResponse struct {
	TxHash string `json:"txhash"`
}

// parseTxHash takes raw output from a command and extracts the transaction hash.
func parseTxHash(rawOutput []byte) (string, error) {
	var resp TxResponse
	if err := json.Unmarshal(rawOutput, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal transaction response: %w", err)
	}
	if resp.TxHash == "" {
		return "", fmt.Errorf("transaction hash not found in the response")
	}
	return resp.TxHash, nil
}
