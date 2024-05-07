package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func listenForSignals(cmds ...*exec.Cmd) {
	// Set up channel to listen for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan

	for _, cmd := range cmds {
		// Stop the process when a signal is received
		stop(cmd)
	}
}
