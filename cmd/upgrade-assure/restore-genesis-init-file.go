package main

import (
	"log"
	"os/exec"
)

func restoreGenesisInitFile(homePath string) {
	// Copy genesis_init.json to genesis.json
	args := []string{
		homePath + "/config/genesis_init.json",
		homePath + "/config/genesis.json",
	}

	if err := exec.Command("cp", args...).Run(); err != nil {
		log.Fatalf(ColorRed+"Failed to copy genesis_init.json to genesis.json: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(ColorYellow + "Genesis file copied to genesis.json")
}
