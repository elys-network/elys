package main

import (
	"fmt"
	"log"
	"time"
)

func waitForTxConfirmation(cmdPath, node, txHash string, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout reached while waiting for tx confirmation")
		}
		success, err := checkTxStatus(cmdPath, node, txHash)
		if err != nil {
			log.Printf("error checking tx status: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if success {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
}
