package main

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func updateGenesis(validatorBalance, homePath, genesisFilePath string) {
	genesis, err := readGenesisFile(genesisFilePath)
	if err != nil {
		log.Fatalf(Red+"Error reading genesis file: %v", err)
	}

	genesisInitFilePath := homePath + "/config/genesis.json"
	genesisInit, err := readGenesisFile(genesisInitFilePath)
	if err != nil {
		log.Fatalf(Red+"Error reading initial genesis file: %v", err)
	}

	filterAccountAddresses := []string{
		"elys1gpv36nyuw5a92hehea3jqaadss9smsqscr3lrp", // remove existing account 0
	}
	filterBalanceAddresses := []string{
		"elys1gpv36nyuw5a92hehea3jqaadss9smsqscr3lrp", // remove existing account 0
		authtypes.NewModuleAddress("distribution").String(),
		authtypes.NewModuleAddress("bonded_tokens_pool").String(),
		authtypes.NewModuleAddress("not_bonded_tokens_pool").String(),
	}

	var coinsToRemove sdk.Coins

	genesis.AppState.Auth.Accounts = filterAccounts(genesis.AppState.Auth.Accounts, filterAccountAddresses)
	genesis.AppState.Bank.Balances, coinsToRemove = filterBalances(genesis.AppState.Bank.Balances, filterBalanceAddresses)

	newValidatorBalance, ok := sdk.NewIntFromString(validatorBalance)
	if !ok {
		panic(Red + "invalid number")
	}
	newValidatorBalanceCoin := sdk.NewCoin("uelys", newValidatorBalance)

	// update supply
	genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.Sub(coinsToRemove...).Add(newValidatorBalanceCoin)

	// Add new validator account and balance
	genesis.AppState.Auth.Accounts = append(genesis.AppState.Auth.Accounts, genesisInit.AppState.Auth.Accounts[0])
	genesis.AppState.Bank.Balances = append(genesis.AppState.Bank.Balances, genesisInit.AppState.Bank.Balances[0])

	// reset staking data
	stakingParams := genesis.AppState.Staking.Params
	genesis.AppState.Staking = genesisInit.AppState.Staking
	genesis.AppState.Staking.Params = stakingParams

	// reset slashing data
	genesis.AppState.Slashing = genesisInit.AppState.Slashing

	// reset distribution data
	genesis.AppState.Distribution = genesisInit.AppState.Distribution

	// set genutil from genesisInit
	genesis.AppState.Genutil = genesisInit.AppState.Genutil

	// reset gov as there are broken proposoals
	genesis.AppState.Gov = genesisInit.AppState.Gov

	// add localhost as allowed client
	genesis.AppState.Ibc.ClientGenesis.Params.AllowedClients = append(genesis.AppState.Ibc.ClientGenesis.Params.AllowedClients, "09-localhost")

	// update voting period
	genesis.AppState.Gov.Params.VotingPeriod = "10s"
	genesis.AppState.Gov.Params.MaxDepositPeriod = "10s"
	genesis.AppState.Gov.Params.MinDeposit = sdk.Coins{sdk.NewInt64Coin("uelys", 10000000)}
	// set deprecated settings
	genesis.AppState.Gov.VotingParams.VotingPeriod = "10s"
	genesis.AppState.Gov.DepositParams.MaxDepositPeriod = "10s"
	genesis.AppState.Gov.DepositParams.MinDeposit = sdk.Coins{sdk.NewInt64Coin("uelys", 10000000)}

	// update commitment params
	genesis.AppState.Commitment.Params.VestingInfos[0].NumMaxVestings = "100000"

	// update wasm params
	// genesis.AppState.Wasm.Params = wasmtypes.DefaultParams()
	genesis.AppState.Wasm = genesisInit.AppState.Wasm

	// update clock params
	genesis.AppState.Clock.Params.ContractAddresses = []string{
		"elys1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqau4f4q",
		"elys17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgs98tvuy",
	}
	genesis.AppState.Clock.Params.ContractGasLimit = "100000000"

	// update broker address
	genesis.AppState.Parameter.Params.BrokerAddress = "elys1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqau4f4q"

	outputFilePath := homePath + "/config/genesis.json"
	if err := writeGenesisFile(outputFilePath, genesis); err != nil {
		log.Fatalf(Red+"Error writing genesis file: %v", err)
	}
}
