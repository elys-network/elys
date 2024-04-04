// nolint: nakedret
package main

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	// global
	flagSkipSnapshot    = "skip-snapshot"
	flagSkipChainInit   = "skip-chain-init"
	flagSkipNodeStart   = "skip-node-start"
	flagSkipProposal    = "skip-proposal"
	flagSkipBinary      = "skip-binary"
	flagChainId         = "chain-id"
	flagKeyringBackend  = "keyring-backend"
	flagGenesisFilePath = "genesis-file-path"
	flagBroadcastMode   = "broadcast-mode"

	// node 1
	flagHome                    = "home"
	flagMoniker                 = "moniker"
	flagValidatorKeyName        = "validator-key-name"
	flagValidatorBalance        = "validator-balance"
	flagValidatorSelfDelegation = "validator-self-delegation"
	flagValidatorMnemonic       = "validator-mnemonic"
	flagRpc                     = "rpc"
	flagP2p                     = "p2p"

	// node 2
	flagHome2                    = "home-2"
	flagMoniker2                 = "moniker-2"
	flagValidatorKeyName2        = "validator-key-name-2"
	flagValidatorBalance2        = "validator-balance-2"
	flagValidatorSelfDelegation2 = "validator-self-delegation-2"
	flagValidatorMnemonic2       = "validator-mnemonic-2"
	flagRpc2                     = "rpc-2"
	flagP2p2                     = "p2p-2"
)

func getFlags(cmd *cobra.Command) (
	// global
	skipSnapshot bool,
	skipChainInit bool,
	skipNodeStart bool,
	skipProposal bool,
	skipBinary bool,
	chainId string,
	keyringBackend string,
	genesisFilePath string,
	broadcastMode string,

	// node 1
	homePath string,
	moniker string,
	validatorKeyName string,
	validatorBalance string,
	validatorSelfDelegation string,
	validatorMnemonic string,
	rpc string,
	p2p string,

	// node 2
	homePath2 string,
	moniker2 string,
	validatorKeyName2 string,
	validatorBalance2 string,
	validatorSelfDelegation2 string,
	validatorMnemonic2 string,
	rpc2 string,
	p2p2 string,
) {
	// global
	skipSnapshot, _ = cmd.Flags().GetBool(flagSkipSnapshot)
	if skipSnapshot {
		log.Printf(Yellow + "skipping snapshot retrieval")
	}

	skipChainInit, _ = cmd.Flags().GetBool(flagSkipChainInit)
	if skipChainInit {
		log.Printf(Yellow + "skipping chain init")
	}

	skipNodeStart, _ = cmd.Flags().GetBool(flagSkipNodeStart)
	if skipNodeStart {
		log.Printf(Yellow + "skipping node start")
	}

	skipProposal, _ = cmd.Flags().GetBool(flagSkipProposal)
	if skipProposal {
		log.Printf(Yellow + "skipping proposal")
	}

	skipBinary, _ = cmd.Flags().GetBool(flagSkipBinary)
	if skipBinary {
		log.Printf(Yellow + "skipping binary download")
	}

	chainId, _ = cmd.Flags().GetString(flagChainId)
	if chainId == "" {
		log.Fatalf(Red + "chain id is required")
	}

	keyringBackend, _ = cmd.Flags().GetString(flagKeyringBackend)
	if keyringBackend == "" {
		log.Fatalf(Red + "keyring backend is required")
	}

	genesisFilePath, _ = cmd.Flags().GetString(flagGenesisFilePath)
	if genesisFilePath == "" {
		log.Fatalf(Red + "genesis file path is required")
	}

	broadcastMode, _ = cmd.Flags().GetString(flagBroadcastMode)
	if broadcastMode == "" {
		log.Fatalf(Red + "broadcast mode is required")
	}

	// node 1
	homePath, _ = cmd.Flags().GetString(flagHome)
	if homePath == "" {
		log.Fatalf(Red + "home path is required")
	}

	moniker, _ = cmd.Flags().GetString(flagMoniker)
	if moniker == "" {
		log.Fatalf(Red + "moniker is required")
	}

	validatorKeyName, _ = cmd.Flags().GetString(flagValidatorKeyName)
	if validatorKeyName == "" {
		log.Fatalf(Red + "validator key name is required")
	}

	validatorBalance, _ = cmd.Flags().GetString(flagValidatorBalance)
	if validatorBalance == "" {
		log.Fatalf(Red + "validator balance is required")
	}

	validatorSelfDelegation, _ = cmd.Flags().GetString(flagValidatorSelfDelegation)
	if validatorSelfDelegation == "" {
		log.Fatalf(Red + "validator self delegation is required")
	}

	validatorMnemonic, _ = cmd.Flags().GetString(flagValidatorMnemonic)
	if validatorMnemonic == "" {
		log.Fatalf(Red + "validator mnemonic is required")
	}

	rpc, _ = cmd.Flags().GetString(flagRpc)
	if rpc == "" {
		log.Fatalf(Red + "rpc is required")
	}

	p2p, _ = cmd.Flags().GetString(flagP2p)
	if p2p == "" {
		log.Fatalf(Red + "p2p is required")
	}

	// node 2
	homePath2, _ = cmd.Flags().GetString(flagHome2)
	if homePath2 == "" {
		log.Fatalf(Red + "home path 2 is required")
	}

	moniker2, _ = cmd.Flags().GetString(flagMoniker2)
	if moniker2 == "" {
		log.Fatalf(Red + "moniker 2 is required")
	}

	validatorKeyName2, _ = cmd.Flags().GetString(flagValidatorKeyName2)
	if validatorKeyName2 == "" {
		log.Fatalf(Red + "validator key name 2 is required")
	}

	validatorBalance2, _ = cmd.Flags().GetString(flagValidatorBalance2)
	if validatorBalance2 == "" {
		log.Fatalf(Red + "validator balance 2 is required")
	}

	validatorSelfDelegation2, _ = cmd.Flags().GetString(flagValidatorSelfDelegation2)
	if validatorSelfDelegation2 == "" {
		log.Fatalf(Red + "validator self delegation 2 is required")
	}

	validatorMnemonic2, _ = cmd.Flags().GetString(flagValidatorMnemonic2)
	if validatorMnemonic2 == "" {
		log.Fatalf(Red + "validator mnemonic 2 is required")
	}

	rpc2, _ = cmd.Flags().GetString(flagRpc2)
	if rpc2 == "" {
		log.Fatalf(Red + "rpc 2 is required")
	}

	p2p2, _ = cmd.Flags().GetString(flagP2p2)
	if p2p2 == "" {
		log.Fatalf(Red + "p2p 2 is required")
	}

	return
}
