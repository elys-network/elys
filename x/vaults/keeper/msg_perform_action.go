package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/vaults/types"
)

func (k msgServer) PerformAction(goCtx context.Context, req *types.MsgPerformAction) (*types.MsgPerformActionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrVaultNotFound, "vault %d not found", req.VaultId)
	}
	if vault.Manager != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "vault %d is not managed by %s", req.VaultId, req.Creator)
	}

	vaultAddress := types.NewVaultAddress(req.VaultId)
	verify := k.AllowedAction(ctx, req.Action, sdk.MustBech32ifyAddressBytes("elys", vaultAddress))
	if !verify {
		return nil, errorsmod.Wrapf(types.ErrInvalidAction, "vault %d does not allow this action", req.VaultId)
	}

	switch perform_action := req.Action.(type) {
	case *types.MsgPerformAction_JoinPool:
		_, _, err := k.amm.JoinPoolNoSwap(ctx, vaultAddress, perform_action.JoinPool.PoolId, perform_action.JoinPool.ShareAmountOut, perform_action.JoinPool.MaxAmountsIn)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
		}
	case *types.MsgPerformAction_ExitPool:
		_, _, _, _, _, err := k.amm.ExitPool(ctx, vaultAddress, perform_action.ExitPool.PoolId, perform_action.ExitPool.ShareAmountIn, perform_action.ExitPool.MinAmountsOut, perform_action.ExitPool.TokenOutDenom, false, true)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
		}
	case *types.MsgPerformAction_SwapByDenom:
		// TODO: check if swap will be executed before end block otherwise we need to check what happened with the coins
		_, err := k.amm.SwapByDenom(ctx, perform_action.SwapByDenom)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
		}
	}

	// get coins after action
	coinsAfter := k.bk.GetAllBalances(ctx, vaultAddress)

	// check if coins after action are only allowed coins
	for _, coin := range coinsAfter {
		found := false
		for _, allowedCoin := range vault.AllowedCoins {
			if coin.Denom == allowedCoin {
				found = true
				break
			}
		}
		if !found {
			return nil, errorsmod.Wrapf(types.ErrInvalidAction, "vault %d does not allow this action", req.VaultId)
		}
	}

	return &types.MsgPerformActionResponse{}, nil
}

func (k Keeper) AllowedAction(ctx sdk.Context, action interface{}, vaultAddress string) bool {
	switch perform_action := action.(type) {
	case *types.MsgPerformAction_JoinPool:
		// Verify join pool fields
		if perform_action.JoinPool == nil {
			return false
		}
		if perform_action.JoinPool.PoolId == 0 {
			return false
		}
		if len(perform_action.JoinPool.MaxAmountsIn) == 0 {
			return false
		}
		if perform_action.JoinPool.ShareAmountOut.IsNil() || perform_action.JoinPool.ShareAmountOut.IsZero() {
			return false
		}
		return true
	case *types.MsgPerformAction_ExitPool:
		// Verify exit pool fields
		if perform_action.ExitPool == nil {
			return false
		}
		if perform_action.ExitPool.PoolId == 0 {
			return false
		}
		if len(perform_action.ExitPool.MinAmountsOut) == 0 {
			return false
		}
		if perform_action.ExitPool.ShareAmountIn.IsNil() || perform_action.ExitPool.ShareAmountIn.IsZero() {
			return false
		}
		return true
	case *types.MsgPerformAction_SwapByDenom:
		// Verify swap by denom fields
		if perform_action.SwapByDenom == nil {
			return false
		}
		if perform_action.SwapByDenom.Amount.IsNil() || perform_action.SwapByDenom.Amount.IsZero() {
			return false
		}
		if perform_action.SwapByDenom.MinAmount.IsNil() || perform_action.SwapByDenom.MinAmount.IsZero() {
			return false
		}
		if perform_action.SwapByDenom.MaxAmount.IsNil() || perform_action.SwapByDenom.MaxAmount.IsZero() {
			return false
		}
		if perform_action.SwapByDenom.DenomIn == "" || perform_action.SwapByDenom.DenomOut == "" {
			return false
		}
		if perform_action.SwapByDenom.Recipient != vaultAddress {
			return false
		}
		return true
	}
	return false
}
