package main

import (
	"log"
	"os/exec"
)

func addGenesisAccount(cmdPath, address, balance, homePath string) {
	// Command and arguments
	args := []string{"add-genesis-account", address, balance + "uelys," + balance + "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65," + balance + "ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4," + balance + "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", "--home", homePath}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err) // nolint: goconst
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"add genesis account with address %s, balance: %s and home path %s successfully", address, balance, homePath)
}
