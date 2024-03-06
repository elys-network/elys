package main

import (
	"log"
	"os/exec"
)

func export(cmdPath, homePath, genesisFilePath string) {
	// Command and arguments
	args := []string{"export", "--home", homePath, "--output-document", genesisFilePath, "--modules-to-export", "amm,assetprofile,auth,authz,bank,burner,capability,clock,commitment,consensus,crisis,epochs,evidence,feegrant,genutil,gov,group,ibc,incentive,interchainaccounts,leveragelp,perpetual,mint,oracle,parameter,params,poolaccounted,stablestake,staking,tokenomics,transfer,transferhook,upgrade,vesting"}

	// Execute the command and capture the output
	cmd := exec.Command(cmdPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command execution failed: %v\nOutput: %s", err, out)
	}

	log.Printf("Output successfully written to %s", genesisFilePath)
}