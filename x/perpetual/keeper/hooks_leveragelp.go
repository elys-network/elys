package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OnLeverageLpEnablePool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	pool := types.NewPool(ammPool)
	k.SetPool(ctx, pool)
	return nil
}

func (k Keeper) OnLeverageLpDisablePool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	// TODO Figure out a way to close positions before deleting it
	// In accounted pool, hooks it checks before if non amm pool token balance is 0 or not. So this won't remove the pool without positions are closed
	// If we add here code to forcefully close all positions then we will have to call this first
	k.RemovePool(ctx, ammPool.PoolId)
	return nil
}

type LeverageLpHooks struct {
	k Keeper
}

var _ leveragelptypes.LeverageLpHooks = LeverageLpHooks{}

// Return the wrapper struct
func (k Keeper) LeverageLpHooks() LeverageLpHooks {
	return LeverageLpHooks{k}
}

func (h LeverageLpHooks) AfterEnablingPool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	return h.k.OnLeverageLpEnablePool(ctx, ammPool)
}

func (h LeverageLpHooks) AfterDisablingPool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	return h.k.OnLeverageLpDisablePool(ctx, ammPool)
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress) error {
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, sender sdk.AccAddress) error {
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress) error {
	return nil
}
