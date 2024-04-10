package main

import (
	"log"
	"os/exec"
)

func generatePrivValidatorState(homePath string) {
	// generate priv_validator_state.json with the following content:
	// {
	// 	"height": "0",
	// 	"round": 0,
	// 	"step": 0
	// }

	// Command and arguments
	args := []string{
		"-c",
		"echo",
		"{\"height\": \"0\", \"round\": 0, \"step\": 0}",
		">",
		homePath + "/data/priv_validator_state.json",
	}

	// Execute the command
	if err := exec.Command("sh", args...).Run(); err != nil {
		log.Fatalf(ColorRed+"Failed to generate priv_validator_state.json: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(ColorYellow + "priv_validator_state.json generated successfully")
}
