package main

import (
	"log"
	"os/exec"
)

func stop(cmds ...*exec.Cmd) {
	for _, cmd := range cmds {
		// Stop the process
		if cmd != nil && cmd.Process != nil {
			err := cmd.Process.Kill()
			if err != nil {
				log.Fatalf(ColorRed+"Failed to kill process: %v", err)
			}
			log.Println(ColorYellow + "Process killed successfully")
		}
	}
}
