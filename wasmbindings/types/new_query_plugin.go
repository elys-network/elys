package types

import (
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ammkeeper "github.com/elys-network/elys/v5/x/amm/keeper"
)

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(
	moduleQueriers []ModuleQuerier,
	amm *ammkeeper.Keeper,
	auth *authkeeper.AccountKeeper,
	bank *bankkeeper.BaseKeeper,
) *QueryPlugin {
	return &QueryPlugin{
		moduleQueriers: moduleQueriers,
		ammKeeper:      amm,
		authKeeper:     auth,
		bankKeeper:     bank,
	}
}
