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
				log.Fatalf(Red+"Failed to kill process: %v", err)
			}
			log.Println(Yellow + "Process killed successfully")
		}
	}
}
