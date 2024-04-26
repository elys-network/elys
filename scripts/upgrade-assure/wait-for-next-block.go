package main

import (
	"log"
	"strconv"
	"time"
)

func waitForNextBlock(cmdPath, node, moniker string, timeout time.Duration) {
	var currentBlockHeight, newBlockHeight int
	var err error

	start := time.Now()

	// First, get the current block height
	for {
		if time.Since(start) > timeout {
			log.Fatalf(ColorRed + "[" + moniker + "] Failed to get current block height within the specified timeout")
		}
		var blockHeightStr string
		blockHeightStr, err = queryBlockHeight(cmdPath, node)
		if err == nil {
			currentBlockHeight, err = strconv.Atoi(blockHeightStr)
			if err == nil && currentBlockHeight > 0 {
				break
			}
		}
		log.Println(ColorYellow + "[" + moniker + "] Waiting for current block height...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	log.Printf(ColorYellow+"["+moniker+"] Current Block Height: %d", currentBlockHeight)

	start = time.Now()

	// Now, wait for the block height to increase
	for {
		if time.Since(start) > timeout {
			log.Fatalf(ColorRed + "[" + moniker + "] Failed to get new block height within the specified timeout")
		}
		var blockHeightStr string
		blockHeightStr, err = queryBlockHeight(cmdPath, node)
		if err == nil {
			newBlockHeight, err = strconv.Atoi(blockHeightStr)
			if err == nil && newBlockHeight > currentBlockHeight {
				break
			}
		}
		log.Println(ColorYellow + "[" + moniker + "] Waiting for next block height...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	log.Printf(ColorYellow+"["+moniker+"] New Block Height: %d", newBlockHeight)
}
