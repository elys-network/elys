package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	leveragelptypes "github.com/elys-network/elys/v7/x/leveragelp/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (k Keeper) OnLeverageLpEnablePool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	params := k.GetParams(ctx)
	pool := types.NewPool(ammPool, params.LeverageMax, params.SafetyFactor)
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

func (h LeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error {
	perpetualPool, found := h.k.GetPool(ctx, ammPool.PoolId)
	if !found {
		return fmt.Errorf("perpetual pool (id: %d) not found", ammPool.PoolId)
	}

	err := h.k.UpdatePoolHealth(ctx, &perpetualPool)
	if err != nil {
		return err
	}

	err = h.k.CheckLowPoolHealthAndMinimumCustody(ctx, ammPool.PoolId, false)
	if err != nil {
		return err
	}
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, _ sdk.AccAddress, ammPool ammtypes.Pool) error {
	perpetualPool, found := h.k.GetPool(ctx, ammPool.PoolId)
	if !found {
		return fmt.Errorf("perpetual pool (id: %d) not found", ammPool.PoolId)
	}

	err := h.k.UpdatePoolHealth(ctx, &perpetualPool)
	if err != nil {
		return err
	}

	err = h.k.CheckLowPoolHealthAndMinimumCustody(ctx, ammPool.PoolId, false)
	if err != nil {
		return err
	}
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, _ sdk.AccAddress, ammPool ammtypes.Pool) error {
	perpetualPool, found := h.k.GetPool(ctx, ammPool.PoolId)
	if !found {
		return fmt.Errorf("perpetual pool (id: %d) not found", ammPool.PoolId)
	}

	err := h.k.UpdatePoolHealth(ctx, &perpetualPool)
	if err != nil {
		return err
	}

	err = h.k.CheckLowPoolHealthAndMinimumCustody(ctx, ammPool.PoolId, false)
	if err != nil {
		return err
	}
	return nil
}
