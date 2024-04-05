package main

import (
	"log"
	"time"
)

func waitForServiceToStart(url, name string) {
	timeout := 120 * time.Second
	start := time.Now()

	// Wait for the node to be running with timout
	for !isServiceRunning(url) {
		if time.Since(start) > timeout {
			log.Fatalf(Red + "Service (" + name + ") did not start within the specified timeout")
		}
		log.Println(Yellow + "Waiting for service (" + name + ") to start...")
		time.Sleep(5 * time.Second)
	}
	log.Println(Yellow + "Service (" + name + ") is running.")
}
