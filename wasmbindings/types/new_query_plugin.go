package types

import (
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
)

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(
	moduleQueriers []ModuleQuerier,
	amm *ammkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	bank *bankkeeper.BaseKeeper,
	staking *stakingkeeper.Keeper,
	commitment *commitmentkeeper.Keeper,
	margin *marginkeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		moduleQueriers:   moduleQueriers,
		ammKeeper:        amm,
		oracleKeeper:     oracle,
		bankKeeper:       bank,
		stakingKeeper:    staking,
		commitmentKeeper: commitment,
		marginKeeper:     margin,
	}
}
