package main

import (
	"log"
	"os/exec"
)

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
