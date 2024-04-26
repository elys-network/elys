package types

import (
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
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	transferhookkeeper "github.com/elys-network/elys/x/transferhook/keeper"
)

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(
	moduleQueriers []ModuleQuerier,
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
	perpetual *perpetualkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	parameter *parameterkeeper.Keeper,
	stablestake *stablestakekeeper.Keeper,
	staking *stakingkeeper.Keeper,
	tokenomics *tokenomicskeeper.Keeper,
	transferhook *transferhookkeeper.Keeper,
	masterchef *masterchefkeeper.Keeper,
	estaking *estakingkeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		moduleQueriers:      moduleQueriers,
		accountedpoolKeeper: accountedpool,
		ammKeeper:           amm,
		assetprofileKeeper:  assetprofile,
		authKeeper:          auth,
		bankKeeper:          bank,
		burnerKeeper:        burner,
		clockKeeper:         clock,
		commitmentKeeper:    commitment,
		epochsKeeper:        epochs,
		incentiveKeeper:     incentive,
		leveragelpKeeper:    leveragelp,
		stakingKeeper:       staking,
		perpetualKeeper:     perpetual,
		oracleKeeper:        oracle,
		parameterKeeper:     parameter,
		stablestakeKeeper:   stablestake,
		tokenomicsKeeper:    tokenomics,
		transferhookKeeper:  transferhook,
		masterchefKeeper:    masterchef,
		estakingKeeper:      estaking,
	}
}
