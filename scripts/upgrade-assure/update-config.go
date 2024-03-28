package main

import (
	"log"
	"os/exec"
	"runtime"
)

// is linux?
func isLinux() bool {
	return runtime.GOOS == "linux"
}

func sed(pattern, file string) {
	// Update config.toml for cors_allowed_origins
	var args []string

	if isLinux() {
		args = []string{"-i", pattern, file}
	} else {
		args = []string{"-i", "", pattern, file}
	}

	// Execute the sed command
	if err := exec.Command("sed", args...).Run(); err != nil {
		log.Fatalf(Red+"Error updating config.toml: %v\n", err)
	}
}

func updateConfig(homePath string) {
	// Path to config files
	configPath := homePath + "/config/config.toml"
	appPath := homePath + "/config/app.toml"

	// Update config.toml for cors_allowed_origins
	sed("s/^cors_allowed_origins =.*/cors_allowed_origins = [\\\"*\\\"]/", configPath)

	// Update config.toml for timeout_broadcast_tx_commit
	sed("s/^timeout_broadcast_tx_commit =.*/timeout_broadcast_tx_commit = \\\"30s\\\"/", configPath)

	// Update config.toml for db_backend
	sed("s/^db_backend =.*/db_backend = \\\"pebbledb\\\"/", configPath)

	// Update app.toml for enabling the API server
	sed("/^# Enable defines if the API server should be enabled./{n;s/enable = false/enable = true/;}", appPath)

	log.Printf(Yellow + "config files have been updated successfully.")
}
