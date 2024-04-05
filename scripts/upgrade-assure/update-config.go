package main

import (
	"log"
)

func updateConfig(homePath, dbEngine string) {
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

	// Update app.toml for enabling the API server
	sed("/^# Enable defines if the API server should be enabled./{n;s/enable = false/enable = true/;}", appPath)

	// Update app.toml for app-db-backend
	sed("s/^app\\-db\\-backend =.*/app\\-db\\-backend = \\\""+dbEngine+"\\\"/", appPath)

	// Update client.toml for keyring-backend
	sed("s/^keyring\\-backend =.*/keyring\\-backend = \\\"test\\\"/", clientPath)

	log.Printf(Yellow + "config files have been updated successfully.")
}
