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

	verify := k.AllowedAction(ctx, req.Action)
	if !verify {
		return nil, errorsmod.Wrapf(types.ErrInvalidAction, "vault %d does not allow this action", req.VaultId)
	}

	// 1. get balance before action
	// 2. perform action
	// 3. get balance after action
	// 4. all coins before and after action should be allowed coins

	return &types.MsgPerformActionResponse{}, nil
}

func (k Keeper) AllowedAction(ctx sdk.Context, action interface{}) bool {
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
