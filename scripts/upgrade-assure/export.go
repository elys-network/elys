package main

import (
	"log"
	"os/exec"
)

func export(cmdPath, homePath, genesisFilePath string) {
	// Command and arguments
	args := []string{"export", "--home", homePath, "--output-document", genesisFilePath, "--modules-to-export", "amm,assetprofile,auth,authz,bank,burner,capability,clock,commitment,consensus,crisis,epochs,evidence,feegrant,genutil,group,ibc,incentive,interchainaccounts,leveragelp,perpetual,mint,oracle,parameter,params,poolaccounted,stablestake,staking,tokenomics,transfer,transferhook,upgrade,vesting"}

	// Execute the command and capture the output
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	log.Printf(Yellow+"Output successfully written to %s", genesisFilePath)
}
