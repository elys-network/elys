package main

import (
	"log"
	"strconv"
	"time"
)

func waitForNextBlock(cmdPath, node, moniker string) {
	var currentBlockHeight, newBlockHeight int
	var err error

	timeout := 120 * time.Second
	start := time.Now()

	// First, get the current block height
	for {
		if time.Since(start) > timeout {
			log.Fatalf(Red + "[" + moniker + "] Failed to get current block height within the specified timeout")
		}
		var blockHeightStr string
		blockHeightStr, err = queryBlockHeight(cmdPath, node)
		if err == nil {
			currentBlockHeight, err = strconv.Atoi(blockHeightStr)
			if err == nil && currentBlockHeight > 0 {
				break
			}
		}
		log.Println(Yellow + "[" + moniker + "] Waiting for current block height...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	log.Printf(Yellow+"["+moniker+"] Current Block Height: %d", currentBlockHeight)

	start = time.Now()

	// Now, wait for the block height to increase
	for {
		if time.Since(start) > timeout {
			log.Fatalf(Red + "[" + moniker + "] Failed to get new block height within the specified timeout")
		}
		var blockHeightStr string
		blockHeightStr, err = queryBlockHeight(cmdPath, node)
		if err == nil {
			newBlockHeight, err = strconv.Atoi(blockHeightStr)
			if err == nil && newBlockHeight > currentBlockHeight {
				break
			}
		}
		log.Println(Yellow + "[" + moniker + "] Waiting for next block height...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	log.Printf(Yellow+"["+moniker+"] New Block Height: %d", newBlockHeight)
}
