// nolint: nakedret
package main

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	// global
	flagOnlyStartWithNewBinary   = "only-start-with-new-binary"
	flagSkipSnapshot             = "skip-snapshot"
	flagSkipChainInit            = "skip-chain-init"
	flagSkipNodeStart            = "skip-node-start"
	flagSkipProposal             = "skip-proposal"
	flagSkipBinary               = "skip-binary"
	flagSkipCreateValidator      = "skip-create-validator"
	flagSkipPrepareValidatorData = "skip-prepare-validator-data"
	flagSkipSubmitProposal       = "skip-submit-proposal"
	flagSkipUpgradeToNewBinary   = "skip-upgrade-to-new-binary"
	flagSkipUnbondValidator      = "skip-unbond-validator"
	flagChainId                  = "chain-id"
	flagKeyringBackend           = "keyring-backend"
	flagGenesisFilePath          = "genesis-file-path"
	flagBroadcastMode            = "broadcast-mode"
	flagDbEngine                 = "db-engine"

	// timeout
	flagTimeOutToWaitForService = "timeout-wait-for-node"
	flagTimeOutNextBlock        = "timeout-next-block"

	// node 1
	flagHome                    = "home"
	flagMoniker                 = "moniker"
	flagValidatorKeyName        = "validator-key-name"
	flagValidatorBalance        = "validator-balance"
	flagValidatorSelfDelegation = "validator-self-delegation"
	flagValidatorMnemonic       = "validator-mnemonic"
	flagRpc                     = "rpc"
	flagP2p                     = "p2p"
	flagPprof                   = "pprof"
	flagApi                     = "api"

	// node 2
	flagHome2                    = "home-2"
	flagMoniker2                 = "moniker-2"
	flagValidatorKeyName2        = "validator-key-name-2"
	flagValidatorBalance2        = "validator-balance-2"
	flagValidatorSelfDelegation2 = "validator-self-delegation-2"
	flagValidatorMnemonic2       = "validator-mnemonic-2"
	flagRpc2                     = "rpc-2"
	flagP2p2                     = "p2p-2"
	flagPprof2                   = "pprof-2"
	flagApi2                     = "api-2"
)

func getFlags(cmd *cobra.Command) (
	// global
	onlyStartWithNewBinaries bool,
	skipSnapshot bool,
	skipChainInit bool,
	skipNodeStart bool,
	skipProposal bool,
	skipBinary bool,
	skipCreateValidator bool,
	skipPrepareValidatorData bool,
	skipSubmitProposal bool,
	skipUpgradeToNewBinary bool,
	skipUnbondValidator bool,
	chainId string,
	keyringBackend string,
	genesisFilePath string,
	broadcastMode string,
	dbEngine string,

	//timeouts
	timeOutWaitForNode int,
	timeOutNextBlock int,

	// node 1
	homePath string,
	moniker string,
	validatorKeyName string,
	validatorBalance string,
	validatorSelfDelegation string,
	validatorMnemonic string,
	rpc string,
	p2p string,
	pprof string,
	api string,

	// node 2
	homePath2 string,
	moniker2 string,
	validatorKeyName2 string,
	validatorBalance2 string,
	validatorSelfDelegation2 string,
	validatorMnemonic2 string,
	rpc2 string,
	p2p2 string,
	pprof2 string,
	api2 string,
) {
	// global
	onlyStartWithNewBinaries, _ = cmd.Flags().GetBool(flagOnlyStartWithNewBinary)
	if onlyStartWithNewBinaries {
		log.Printf(ColorYellow + "only starting with new binaries, skipping all other steps")
	}

	skipSnapshot, _ = cmd.Flags().GetBool(flagSkipSnapshot)
	if skipSnapshot {
		log.Printf(ColorYellow + "skipping snapshot retrieval")
	}

	skipChainInit, _ = cmd.Flags().GetBool(flagSkipChainInit)
	if skipChainInit {
		log.Printf(ColorYellow + "skipping chain init")
	}

	skipNodeStart, _ = cmd.Flags().GetBool(flagSkipNodeStart)
	if skipNodeStart {
		log.Printf(ColorYellow + "skipping node start")
	}

	skipProposal, _ = cmd.Flags().GetBool(flagSkipProposal)
	if skipProposal {
		log.Printf(ColorYellow + "skipping proposal")
	}

	skipBinary, _ = cmd.Flags().GetBool(flagSkipBinary)
	if skipBinary {
		log.Printf(ColorYellow + "skipping binary download")
	}

	skipCreateValidator, _ = cmd.Flags().GetBool(flagSkipCreateValidator)
	if skipCreateValidator {
		log.Printf(ColorYellow + "skipping create validator")
	}

	skipPrepareValidatorData, _ = cmd.Flags().GetBool(flagSkipPrepareValidatorData)
	if skipPrepareValidatorData {
		log.Printf(ColorYellow + "skipping prepare validator data")
	}

	skipSubmitProposal, _ = cmd.Flags().GetBool(flagSkipSubmitProposal)
	if skipSubmitProposal {
		log.Printf(ColorYellow + "skipping submit proposal")
	}

	skipUpgradeToNewBinary, _ = cmd.Flags().GetBool(flagSkipUpgradeToNewBinary)
	if skipUpgradeToNewBinary {
		log.Printf(ColorYellow + "skipping upgrade to new binary")
	}

	skipUnbondValidator, _ = cmd.Flags().GetBool(flagSkipUnbondValidator)
	if skipUnbondValidator {
		log.Printf(ColorYellow + "skipping unbond validator")
	}

	chainId, _ = cmd.Flags().GetString(flagChainId)
	if chainId == "" {
		log.Fatalf(ColorRed + "chain id is required")
	}

	keyringBackend, _ = cmd.Flags().GetString(flagKeyringBackend)
	if keyringBackend == "" {
		log.Fatalf(ColorRed + "keyring backend is required")
	}

	genesisFilePath, _ = cmd.Flags().GetString(flagGenesisFilePath)
	if genesisFilePath == "" {
		log.Fatalf(ColorRed + "genesis file path is required")
	}

	broadcastMode, _ = cmd.Flags().GetString(flagBroadcastMode)
	if broadcastMode == "" {
		log.Fatalf(ColorRed + "broadcast mode is required")
	}

	dbEngine, _ = cmd.Flags().GetString(flagDbEngine)
	if dbEngine == "" {
		log.Fatalf(ColorRed + "database engine is required")
	}

	timeOutWaitForNode, err := cmd.Flags().GetInt(flagTimeOutToWaitForService)

	if err != nil {
		log.Fatalf(ColorRed + err.Error())
	}

	if timeOutWaitForNode == 0 {
		log.Fatalf(ColorRed + "time out to wait for service is required")
	}

	timeOutNextBlock, err = cmd.Flags().GetInt(flagTimeOutNextBlock)

	if err != nil {
		log.Fatalf(ColorRed + err.Error())
	}

	if timeOutNextBlock == 0 {
		log.Fatalf(ColorRed + "time out next block is required")
	}

	// node 1
	homePath, _ = cmd.Flags().GetString(flagHome)
	if homePath == "" {
		log.Fatalf(ColorRed + "home path is required")
	}

	moniker, _ = cmd.Flags().GetString(flagMoniker)
	if moniker == "" {
		log.Fatalf(ColorRed + "moniker is required")
	}

	validatorKeyName, _ = cmd.Flags().GetString(flagValidatorKeyName)
	if validatorKeyName == "" {
		log.Fatalf(ColorRed + "validator key name is required")
	}

	validatorBalance, _ = cmd.Flags().GetString(flagValidatorBalance)
	if validatorBalance == "" {
		log.Fatalf(ColorRed + "validator balance is required")
	}

	validatorSelfDelegation, _ = cmd.Flags().GetString(flagValidatorSelfDelegation)
	if validatorSelfDelegation == "" {
		log.Fatalf(ColorRed + "validator self delegation is required")
	}

	validatorMnemonic, _ = cmd.Flags().GetString(flagValidatorMnemonic)
	if validatorMnemonic == "" {
		log.Fatalf(ColorRed + "validator mnemonic is required")
	}

	rpc, _ = cmd.Flags().GetString(flagRpc)
	if rpc == "" {
		log.Fatalf(ColorRed + "rpc is required")
	}

	p2p, _ = cmd.Flags().GetString(flagP2p)
	if p2p == "" {
		log.Fatalf(ColorRed + "p2p is required")
	}

	pprof, _ = cmd.Flags().GetString(flagPprof)
	if pprof == "" {
		log.Fatalf(ColorRed + "pprof is required")
	}

	api, _ = cmd.Flags().GetString(flagApi)
	if api == "" {
		log.Fatalf(ColorRed + "api is required")
	}

	// node 2
	homePath2, _ = cmd.Flags().GetString(flagHome2)
	if homePath2 == "" {
		log.Fatalf(ColorRed + "home path 2 is required")
	}

	moniker2, _ = cmd.Flags().GetString(flagMoniker2)
	if moniker2 == "" {
		log.Fatalf(ColorRed + "moniker 2 is required")
	}

	validatorKeyName2, _ = cmd.Flags().GetString(flagValidatorKeyName2)
	if validatorKeyName2 == "" {
		log.Fatalf(ColorRed + "validator key name 2 is required")
	}

	validatorBalance2, _ = cmd.Flags().GetString(flagValidatorBalance2)
	if validatorBalance2 == "" {
		log.Fatalf(ColorRed + "validator balance 2 is required")
	}

	validatorSelfDelegation2, _ = cmd.Flags().GetString(flagValidatorSelfDelegation2)
	if validatorSelfDelegation2 == "" {
		log.Fatalf(ColorRed + "validator self delegation 2 is required")
	}

	validatorMnemonic2, _ = cmd.Flags().GetString(flagValidatorMnemonic2)
	if validatorMnemonic2 == "" {
		log.Fatalf(ColorRed + "validator mnemonic 2 is required")
	}

	rpc2, _ = cmd.Flags().GetString(flagRpc2)
	if rpc2 == "" {
		log.Fatalf(ColorRed + "rpc 2 is required")
	}

	p2p2, _ = cmd.Flags().GetString(flagP2p2)
	if p2p2 == "" {
		log.Fatalf(ColorRed + "p2p 2 is required")
	}

	pprof2, _ = cmd.Flags().GetString(flagPprof2)
	if pprof2 == "" {
		log.Fatalf(ColorRed + "pprof 2 is required")
	}

	api2, _ = cmd.Flags().GetString(flagApi2)
	if api2 == "" {
		log.Fatalf(ColorRed + "api 2 is required")
	}

	return
}
