package main

import (
	"log"
	"os"
	"time"

	elyscmd "github.com/elys-network/elys/cmd/elysd/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "upgrade-assure [snapshot_url] [old_binary_url] [new_binary_url] [flags]",
		Short: "Upgrade Assure is a tool for running a chain from a snapshot and to test out the upgrade process.",
		Long:  `A tool for running a chain from a snapshot.`,
		Args:  cobra.ExactArgs(3), // Expect exactly 1 argument
		Run: func(cmd *cobra.Command, args []string) {
			snapshotUrl, oldBinaryUrl, newBinaryUrl := getArgs(args)
			// global flags
			skipSnapshot, skipChainInit, skipNodeStart, skipProposal, skipBinary, chainId, keyringBackend, genesisFilePath, broadcastMode, dbEngine,
				// node 1 flags
				homePath, moniker, validatorKeyName, validatorBalance, validatorSelfDelegation, validatorMnemonic, rpc, p2p,
				// node 2 flags
				homePath2, moniker2, validatorKeyName2, validatorBalance2, validatorSelfDelegation2, validatorMnemonic2, rpc2, p2p2 := getFlags(cmd)

			// set address prefix
			elyscmd.InitSDKConfig()

			// download and run old binary
			oldBinaryPath, oldVersion, err := downloadAndRunVersion(oldBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(ColorRed+"Error downloading and running old binary: %v", err)
			}

			// print old binary path and version
			log.Printf(ColorGreen+"Old binary path: %v and version: %v", oldBinaryPath, oldVersion)

			// download and run new binary
			newBinaryPath, newVersion, err := downloadAndRunVersion(newBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(ColorRed+"Error downloading and running new binary: %v", err)
			}

			// print new binary path and version
			log.Printf(ColorGreen+"New binary path: %v and version: %v", newBinaryPath, newVersion)

			if !skipSnapshot {
				// remove home path
				removeHome(homePath)

				// init chain
				initNode(oldBinaryPath, moniker, chainId, homePath)

				// update config files
				updateConfig(homePath, dbEngine)

				// retrieve the snapshot
				retrieveSnapshot(snapshotUrl, homePath)

				// export genesis file
				export(oldBinaryPath, homePath, genesisFilePath)
			}

			if !skipChainInit {
				// remove home paths
				removeHome(homePath)
				removeHome(homePath2)

				// init nodes
				initNode(oldBinaryPath, moniker, chainId, homePath)
				initNode(oldBinaryPath, moniker2, chainId, homePath2)

				// update config files to enable api and cors
				updateConfig(homePath, dbEngine)
				updateConfig(homePath2, dbEngine)

				// query node 1 id
				node1Id := queryNodeId(oldBinaryPath, homePath)

				// add peers
				addPeers(homePath2, p2p, node1Id)

				// add validator keys to node 1
				validatorAddress := addKey(oldBinaryPath, validatorKeyName, validatorMnemonic, homePath, keyringBackend)
				validatorAddress2 := addKey(oldBinaryPath, validatorKeyName2, validatorMnemonic2, homePath, keyringBackend)

				// add validator keys to node 2
				_ = addKey(oldBinaryPath, validatorKeyName, validatorMnemonic, homePath2, keyringBackend)
				_ = addKey(oldBinaryPath, validatorKeyName2, validatorMnemonic2, homePath2, keyringBackend)

				// add genesis accounts
				addGenesisAccount(oldBinaryPath, validatorAddress, validatorBalance, homePath)
				addGenesisAccount(oldBinaryPath, validatorAddress2, validatorBalance2, homePath)

				// generate genesis tx
				genTx(oldBinaryPath, validatorKeyName, validatorSelfDelegation, chainId, homePath, keyringBackend)

				// collect genesis txs
				collectGentxs(oldBinaryPath, homePath)

				// validate genesis
				validateGenesis(oldBinaryPath, homePath)

				// update genesis
				updateGenesis(validatorBalance, homePath, genesisFilePath)
			}

			if !skipNodeStart {
				// start node 1
				oldBinaryCmd := start(oldBinaryPath, homePath, rpc, p2p, moniker, ColorGreen, ColorRed)

				// wait for rpc to start
				waitForServiceToStart(rpc, moniker)

				// wait for next block
				waitForNextBlock(oldBinaryPath, rpc, moniker)

				if skipProposal {
					// listen for signals
					listenForSignals(oldBinaryCmd)
					return
				}

				// query validator pubkey
				validatorPubkey2 := queryValidatorPubkey(oldBinaryPath, homePath2)

				// create validator node 2
				createValidator(oldBinaryPath, validatorKeyName2, validatorSelfDelegation2, moniker2, validatorPubkey2, homePath, keyringBackend, chainId, rpc, broadcastMode)

				// wait for next block
				waitForNextBlock(oldBinaryPath, rpc, moniker)

				// stop old binary
				stop(oldBinaryCmd)

				// copy data from node 1 to node 2
				copyDataFromNodeToNode(homePath, homePath2)

				// generate priv_validator_state.json file for node 2
				generatePrivValidatorState(homePath2)

				// start node 1 and 2
				oldBinaryCmd = start(oldBinaryPath, homePath, rpc, p2p, moniker, ColorGreen, ColorRed)
				oldBinaryCmd2 := start(oldBinaryPath, homePath2, rpc2, p2p2, moniker2, ColorGreen, ColorRed)

				// wait for rpc 1 and 2 to start
				waitForServiceToStart(rpc, moniker)
				waitForServiceToStart(rpc2, moniker2)

				// query and calculate upgrade block height
				upgradeBlockHeight := queryAndCalcUpgradeBlockHeight(oldBinaryPath, rpc)

				// query next proposal id
				proposalId, err := queryNextProposalId(oldBinaryPath, rpc)
				if err != nil {
					log.Printf(ColorYellow+"Error querying next proposal id: %v", err)
					log.Printf(ColorYellow + "Setting proposal id to 1")
					proposalId = "1"
				}

				// submit upgrade proposal
				txHash := submitUpgradeProposal(oldBinaryPath, validatorKeyName, newVersion, upgradeBlockHeight, homePath, keyringBackend, chainId, rpc, broadcastMode)

				err = waitForTxConfirmation(oldBinaryPath, rpc, txHash, 5*time.Minute)
				if err != nil {
					log.Fatalf("upgrade proposal not confirmed: %v", err)
				}

				// vote on upgrade proposal
				txHash = voteOnUpgradeProposal(oldBinaryPath, validatorKeyName, proposalId, homePath, keyringBackend, chainId, rpc, broadcastMode)

				err = waitForTxConfirmation(oldBinaryPath, rpc, txHash, 5*time.Minute)
				if err != nil {
					log.Fatalf("voting on upgrade proposal not confirmed: %v", err)
				}

				// wait for upgrade block height
				waitForBlockHeight(oldBinaryPath, rpc, upgradeBlockHeight)

				// wait 5 seconds
				time.Sleep(5 * time.Second)

				// stop old binaries
				stop(oldBinaryCmd, oldBinaryCmd2)

				// wait 5 seconds
				time.Sleep(5 * time.Second)

				// start new binary
				newBinaryCmd := start(newBinaryPath, homePath, rpc, p2p, moniker, "\033[32m", "\033[31m")
				newBinaryCmd2 := start(newBinaryPath, homePath2, rpc2, p2p2, moniker2, "\033[32m", "\033[31m")

				// wait for node to start
				waitForServiceToStart(rpc, moniker)
				waitForServiceToStart(rpc2, moniker2)

				// wait for next block
				waitForNextBlock(newBinaryPath, rpc, moniker)
				waitForNextBlock(newBinaryPath, rpc2, moniker2)

				// check if the upgrade was successful
				queryUpgradeApplied(newBinaryPath, rpc, newVersion)
				queryUpgradeApplied(newBinaryPath, rpc2, newVersion)

				// stop new binaries
				stop(newBinaryCmd, newBinaryCmd2)
			}
		},
	}

	// get HOME environment variable
	homeEnv, _ := os.LookupEnv("HOME")

	// global flags
	rootCmd.PersistentFlags().Bool(flagSkipSnapshot, false, "skip snapshot retrieval")
	rootCmd.PersistentFlags().Bool(flagSkipChainInit, false, "skip chain init")
	rootCmd.PersistentFlags().Bool(flagSkipNodeStart, false, "skip node start")
	rootCmd.PersistentFlags().Bool(flagSkipProposal, false, "skip proposal")
	rootCmd.PersistentFlags().Bool(flagSkipBinary, false, "skip binary download")
	rootCmd.PersistentFlags().String(flagChainId, "elystestnet-1", "chain id")
	rootCmd.PersistentFlags().String(flagKeyringBackend, "test", "keyring backend")
	rootCmd.PersistentFlags().String(flagGenesisFilePath, "/tmp/genesis.json", "genesis file path")
	rootCmd.PersistentFlags().String(flagBroadcastMode, "sync", "broadcast mode")
	rootCmd.PersistentFlags().String(flagDbEngine, "pebbledb", "database engine to use")

	// node 1 flags
	rootCmd.PersistentFlags().String(flagHome, homeEnv+"/.elys", "home directory")
	rootCmd.PersistentFlags().String(flagMoniker, "alice", "moniker")
	rootCmd.PersistentFlags().String(flagValidatorKeyName, "validator", "validator key name")
	rootCmd.PersistentFlags().String(flagValidatorBalance, "200000000000000", "validator balance")
	rootCmd.PersistentFlags().String(flagValidatorSelfDelegation, "50000000000000", "validator self delegation")
	rootCmd.PersistentFlags().String(flagValidatorMnemonic, "bargain toss help way dash forget bar casual boat drill execute ordinary human lecture leopard enroll joy rural shed express kite sample brick void", "validator mnemonic")
	rootCmd.PersistentFlags().String(flagRpc, "tcp://0.0.0.0:26657", "rpc")
	rootCmd.PersistentFlags().String(flagP2p, "tcp://0.0.0.0:26656", "p2p")

	// node 2 flags
	rootCmd.PersistentFlags().String(flagHome2, homeEnv+"/.elys2", "home directory 2")
	rootCmd.PersistentFlags().String(flagMoniker2, "bob", "moniker 2")
	rootCmd.PersistentFlags().String(flagValidatorKeyName2, "validator-2", "validator key name 2")
	rootCmd.PersistentFlags().String(flagValidatorBalance2, "200000000000000", "validator balance 2")
	rootCmd.PersistentFlags().String(flagValidatorSelfDelegation2, "1000000", "validator self delegation 2")
	rootCmd.PersistentFlags().String(flagValidatorMnemonic2, "kidney seat stay demand panel garlic uncle flock plunge logic link owner laugh sponsor desk scare pipe derive trick smart coffee goat arrange cause", "validator mnemonic 2")
	rootCmd.PersistentFlags().String(flagRpc2, "tcp://0.0.0.0:26667", "rpc")
	rootCmd.PersistentFlags().String(flagP2p2, "tcp://0.0.0.0:26666", "p2p")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(ColorRed+"Error executing command: %v", err)
	}
}
