package main

import (
	"log"
	"strconv"
)

func queryAndCalcUpgradeBlockHeight(cmdPath, node string) string {
	// query block height
	blockHeight, err := queryBlockHeight(cmdPath, node)
	if err != nil {
		log.Fatalf(ColorRed+"Failed to query block height: %v", err)
	}

	// Convert blockHeight from string to int
	blockHeightInt, err := strconv.Atoi(blockHeight)
	if err != nil {
		log.Fatalf(ColorRed+"Failed to convert blockHeight to integer: %v", err)
	}

	// set upgrade block height
	upgradeBlockHeight := blockHeightInt + 10

	// return upgrade block height as a string
	return strconv.Itoa(upgradeBlockHeight)
}
