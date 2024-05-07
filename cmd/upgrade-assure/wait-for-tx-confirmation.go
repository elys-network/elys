package main

import (
	"log"
	"time"
)

func waitForTxConfirmation(cmdPath, node, txHash string, timeout time.Duration) {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			log.Fatalf(ColorRed + "timeout reached while waiting for tx confirmation")
		}
		success, err := checkTxStatus(cmdPath, node, txHash)
		if err != nil {
			log.Printf(ColorRed+"error checking tx status, retrying in 5 seconds: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if success {
			break
		}
		log.Printf(ColorYellow+"waiting for tx confirmation %s", txHash)
		time.Sleep(5 * time.Second)
	}
	log.Printf(ColorGreen+"tx %s confirmed", txHash)
}
