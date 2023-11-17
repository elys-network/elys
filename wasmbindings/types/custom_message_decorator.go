package types

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	accountedpoolkeeper "github.com/elys-network/elys/x/accountedpool/keeper"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	assetprofilekeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	burnerkeeper "github.com/elys-network/elys/x/burner/keeper"
	clockkeeper "github.com/elys-network/elys/x/clock/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	epochskeeper "github.com/elys-network/elys/x/epochs/keeper"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	transferhookkeeper "github.com/elys-network/elys/x/transferhook/keeper"
)

func CustomMessageDecorator(
	moduleMessengers []ModuleMessenger,
	accountedpool *accountedpoolkeeper.Keeper,
	amm *ammkeeper.Keeper,
	assetprofile *assetprofilekeeper.Keeper,
	auth *authkeeper.AccountKeeper,
	bank *bankkeeper.BaseKeeper,
	burner *burnerkeeper.Keeper,
	clock *clockkeeper.Keeper,
	commitment *commitmentkeeper.Keeper,
	epochs *epochskeeper.Keeper,
	incentive *incentivekeeper.Keeper,
	leveragelp *leveragelpkeeper.Keeper,
	margin *marginkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	parameter *parameterkeeper.Keeper,
	stablestake *stablestakekeeper.Keeper,
	staking *stakingkeeper.Keeper,
	tokenomics *tokenomicskeeper.Keeper,
	transferhook *transferhookkeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:          old,
			moduleMessengers: moduleMessengers,
			accountedpool:    accountedpool,
			amm:              amm,
			assetprofile:     assetprofile,
			auth:             auth,
			bank:             bank,
			burner:           burner,
			clock:            clock,
			commitment:       commitment,
			epochs:           epochs,
			incentive:        incentive,
			leveragelp:       leveragelp,
			margin:           margin,
			oracle:           oracle,
			parameter:        parameter,
			stablestake:      stablestake,
			staking:          staking,
			tokenomics:       tokenomics,
			transferhook:     transferhook,
		}
	}
}
