package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// PromptWriter wraps an io.Writer and adds color codes to the output
type PromptWriter struct {
	w       *os.File
	color   string
	moniker string
}

// Write adds color codes to the data and writes it to the log
func (cw PromptWriter) Write(data []byte) (int, error) {
	// Add color codes to the data
	coloredData := []byte(cw.color + "[" + cw.moniker + "]" + ColorReset + " " + string(data))
	_, err := cw.w.Write(coloredData)
	if err != nil {
		log.Fatalf("Error writing to log: %v", err)
	}
	return len(data), err
}

func start(cmdPath, homePath, rpc, p2p, pprof, api, moniker, successColor, errorColor string) *exec.Cmd {
	// Set the log level
	logLevel := "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}

	// Command and arguments
	args := []string{"start", "--home", homePath, "--rpc.laddr", rpc, "--p2p.laddr", p2p, "--rpc.pprof_laddr", pprof, "--api.address", api, "--log_level", logLevel}

	// Set up the command
	cmd := exec.Command(cmdPath, args...)

	// Use PromptWriter to handle logging for standard output and error
	cmd.Stdout = PromptWriter{w: os.Stdout, color: successColor, moniker: moniker} // ColorGreen for stdout
	cmd.Stderr = PromptWriter{w: os.Stderr, color: errorColor, moniker: moniker}   // ColorRed for stderr

	// Execute the command and stream the output in a goroutine to avoid blocking
	go func() {
		err := cmd.Run()
		if err != nil {
			// Check if the error is because of the process being killed
			if exitErr, ok := err.(*exec.ExitError); ok {
				// If the process was killed, log it as a non-fatal error
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
					log.Printf("Process was killed: %v", err)
					return
				}
			}
			// For other errors, log them as fatal
			log.Fatalf("Command execution failed: %v", err)
		}
	}()

	return cmd
}
