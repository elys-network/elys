package main

import (
	"log"
	"os/exec"
)

func updateConfig(homePath string) {
	// Path to config files
	configPath := homePath + "/config/config.toml"
	appPath := homePath + "/config/app.toml"

	// Update config.toml
	args := []string{"-i", "", "s/^cors_allowed_origins =.*/cors_allowed_origins = [\\\"*\\\"]/", configPath}

	// Execute the sed command
	if err := exec.Command("sed", args...).Run(); err != nil {
		log.Fatalf(Red+"Error updating config.toml: %v\n", err)
	}

	// Update app.toml
	args = []string{"-i", "", "/^# Enable defines if the API server should be enabled./{n;s/enable = false/enable = true/;}", appPath}

	// Execute the sed command
	if err := exec.Command("sed", args...).Run(); err != nil {
		log.Fatalf(Red+"Error updating app.toml: %s\n", err)
	}

	log.Printf(Yellow + "config files have been updated successfully.")
}
