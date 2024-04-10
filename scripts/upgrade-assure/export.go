package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func export(cmdPath, homePath, genesisFilePath string) {
	// Define modules in a slice
	modules := []string{
		"amm",
		"assetprofile",
		"auth",
		"authz",
		"bank",
		"burner",
		"capability",
		"clock",
		"commitment",
		"consensus",
		"crisis",
		"epochs",
		"evidence",
		"feegrant",
		"genutil",
		"gov",
		"group",
		"ibc",
		"incentive",
		"interchainaccounts",
		"leveragelp",
		"perpetual",
		"mint",
		"oracle",
		"parameter",
		"params",
		"poolaccounted",
		"stablestake",
		"staking",
		"tokenomics",
		"transfer",
		"transferhook",
		"upgrade",
		"vesting",
	}

	// Combine the modules into a comma-separated string
	modulesStr := strings.Join(modules, ",")

	// Command and arguments
	args := []string{"export", "--home", homePath, "--modules-to-export", modulesStr}

	// Create the command
	cmd := exec.Command(cmdPath, args...)

	// Create the output file
	outFile, err := os.Create(genesisFilePath)
	if err != nil {
		// Handle error
		panic(err)
	}
	defer outFile.Close()

	// Redirect the output to a file
	cmd.Stdout = outFile
	// cmd.Stderr = outFile // You can also direct stderr to a different file or os.Stderr

	// Execute the command
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command execution failed: %v\nOutput: %s", err, out)
	}

	log.Printf("Output successfully written to %s", genesisFilePath)
}
