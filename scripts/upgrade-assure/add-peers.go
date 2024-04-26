package main

import (
	"log"
	"strings"
)

func addPeers(homePath, p2p, nodeId string) {
	// Path to config files
	configPath := homePath + "/config/config.toml"

	// update p2p url to remove the `tcp://` or `http://` or `https://` prefix
	p2p = strings.TrimPrefix(p2p, "tcp://")
	p2p = strings.TrimPrefix(p2p, "http://")
	p2p = strings.TrimPrefix(p2p, "https://")

	// escape the `:` character from p2p
	p2p = strings.ReplaceAll(p2p, ":", "\\:")
	// escape the `.` character from p2p
	p2p = strings.ReplaceAll(p2p, ".", "\\.")

	sed("s/^persistent_peers =.*/persistent_peers = \\\""+nodeId+"\\@"+p2p+"\\\"/", configPath)

	log.Printf(ColorYellow + "peers have been added successfully.")
}
