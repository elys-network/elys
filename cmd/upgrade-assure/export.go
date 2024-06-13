package main

import (
	"log"
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
		// "bank", // FIXME: causes the init genesis state from genesis.json to OOM (kill process)
		"burner",
		"capability",
		"clock",
		// "commitment", // FIXME: causes the balance mismatch error
		"consensus",
		"crisis",
		"distribution",
		"epochs",
		"estaking",
		"evidence",
		"feegrant",
		"genutil",
		// "gov", // FIXME: should be re-enabled after gov proposal fixes
		"group",
		"ibc",
		"incentive",
		"interchainaccounts",
		"leveragelp",
		"masterchef",
		"perpetual",
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
	args := []string{"export", "--home", homePath, "--output-document", genesisFilePath, "--modules-to-export", modulesStr}

	// Execute the command and capture the output
	cmd := exec.Command(cmdPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command execution failed: %v\nOutput: %s", err, out)
	}

	log.Printf("Output successfully written to %s", genesisFilePath)
}
