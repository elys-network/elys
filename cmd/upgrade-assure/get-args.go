package main

import (
	"log"
)

func getArgs(args []string) (snapshotUrl, oldBinaryUrl, newBinaryUrl string) {
	snapshotUrl = args[0] // https://snapshots.polkachu.com/testnet-snapshots/elys/elys_4392223.tar.lz4
	if snapshotUrl == "" {
		log.Fatalf(ColorRed + "snapshot url is required")
	}

	oldBinaryUrl = args[1] // https://github.com/elys-network/elys/releases/download/v0.19.0/elysd-v0.19.0-darwin-arm64
	if oldBinaryUrl == "" {
		log.Fatalf(ColorRed + "old binary url is required")
	}

	newBinaryUrl = args[2] // https://github.com/elys-network/elys/releases/download/v0.20.0/elysd-v0.20.0-darwin-arm64
	if newBinaryUrl == "" {
		log.Fatalf(ColorRed + "new binary url is required")
	}

	return
}
