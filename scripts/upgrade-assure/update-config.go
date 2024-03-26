package main

import (
	"log"
	"os/exec"
	"runtime"
)

func isLinux() bool {
	// Check if the OS is Linux
	return runtime.GOOS == "linux"
}

func updateConfig(homePath string) {
	// Path to config files
	configPath := homePath + "/config/config.toml"
	appPath := homePath + "/config/app.toml"

	// Update config.toml for cors_allowed_origins
	args := []string{"-i", "", "s/^cors_allowed_origins =.*/cors_allowed_origins = [\\\"*\\\"]/", configPath}

	// Execute the sed command
	if err := exec.Command("sed", args...).Run(); err != nil {
		log.Fatalf(Red+"Error updating config.toml: %v\n", err)
	}

	// Update config.toml for timeout_broadcast_tx_commit
	args = []string{"-i", "", "s/^timeout_broadcast_tx_commit =.*/timeout_broadcast_tx_commit = \\\"30s\\\"/", configPath}

	// Execute the sed command
	if err := exec.Command("sed", args...).Run(); err != nil {
		log.Fatalf(Red+"Error updating config.toml: %v\n", err)
	}

	// only apply this change if the os is linux
	if isLinux() {
		args = []string{"-i", "", "s/^db_backend =.*/db_backend = \\\"rocksdb\\\"/", configPath}

		// Execute the sed command
		if err := exec.Command("sed", args...).Run(); err != nil {
			log.Fatalf(Red+"Error updating config.toml: %v\n", err)
		}
	}

	// Update app.toml for enabling the API server
	args = []string{"-i", "", "/^# Enable defines if the API server should be enabled./{n;s/enable = false/enable = true/;}", appPath}

	// Execute the sed command
	if err := exec.Command("sed", args...).Run(); err != nil {
		log.Fatalf(Red+"Error updating app.toml: %s\n", err)
	}

	log.Printf(Yellow + "config files have been updated successfully.")
}
