package main

import (
	"log"
	"time"
)

func waitForServiceToStart(url, name string, timeoutFlag int) {
	timeout := time.Duration(timeoutFlag) * time.Second
	start := time.Now()

	// Wait for the node to be running with timout
	for !isServiceRunning(url) {
		if time.Since(start) > timeout {
			log.Fatalf(ColorRed + "[" + name + "] Service did not start within the specified timeout")
		}
		log.Println(ColorYellow + "[" + name + "] Waiting for service to start...")
		time.Sleep(5 * time.Second)
	}
	log.Println(ColorYellow + "[" + name + "] Service is running.")
}
