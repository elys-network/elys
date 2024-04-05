package main

import (
	"log"
	"os/exec"
	"runtime"
	"strings"
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
		log.Fatalf(Red+"Error updating "+file+": %v\n", err)
	}
}

func updateConfig(homePath, p2p, nodeId, dbEngine string) {
	// Path to config files
	configPath := homePath + "/config/config.toml"
	appPath := homePath + "/config/app.toml"
	clientPath := homePath + "/config/client.toml"

	// Update config.toml for cors_allowed_origins
	sed("s/^cors_allowed_origins =.*/cors_allowed_origins = [\\\"*\\\"]/", configPath)

	// Update config.toml for timeout_broadcast_tx_commit
	sed("s/^timeout_broadcast_tx_commit =.*/timeout_broadcast_tx_commit = \\\"30s\\\"/", configPath)

	// Update config.toml for db_backend
	sed("s/^db_backend =.*/db_backend = \\\""+dbEngine+"\\\"/", configPath)

	// update p2p url to remove the `tcp://` or `http://` or `https://` prefix
	p2p = strings.TrimPrefix(p2p, "tcp://")
	p2p = strings.TrimPrefix(p2p, "http://")
	p2p = strings.TrimPrefix(p2p, "https://")

	// escape the `:` character from p2p
	p2p = strings.ReplaceAll(p2p, ":", "\\:")
	// escape the `.` character from p2p
	p2p = strings.ReplaceAll(p2p, ".", "\\.")

	sed("s/^persistent_peers =.*/persistent_peers = \\\""+nodeId+"\\@"+p2p+"\\\"/", configPath)

	// Update app.toml for enabling the API server
	sed("/^# Enable defines if the API server should be enabled./{n;s/enable = false/enable = true/;}", appPath)

	// Update app.toml for app-db-backend
	sed("s/^app\\-db\\-backend =.*/app\\-db\\-backend = \\\""+dbEngine+"\\\"/", appPath)

	// Update client.toml for keyring-backend
	sed("s/^keyring\\-backend =.*/keyring\\-backend = \\\"test\\\"/", clientPath)

	log.Printf(Yellow + "config files have been updated successfully.")
}
