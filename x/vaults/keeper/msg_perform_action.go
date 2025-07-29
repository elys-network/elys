package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/v7/x/amm/types"

	"github.com/elys-network/elys/v7/x/vaults/types"
)

func (k msgServer) PerformActionJoinPool(goCtx context.Context, req *types.MsgPerformActionJoinPool) (*types.MsgPerformActionJoinPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrVaultNotFound, "vault %d not found", req.VaultId)
	}
	if vault.Manager != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "vault %d is not managed by %s", req.VaultId, req.Creator)
	}
	vaultAddress := types.NewVaultAddress(req.VaultId)

	_, sharesOut, err := k.amm.JoinPoolNoSwap(ctx, vaultAddress, req.PoolId, req.ShareAmountOut, req.MaxAmountsIn)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
	}

	// get coins after action
	coinsAfter := k.bk.GetAllBalances(ctx, vaultAddress)
	allCommitments := k.commitment.GetCommitments(ctx, vaultAddress)
	for _, commitment := range allCommitments.CommittedTokens {
		coinsAfter = coinsAfter.Add(sdk.NewCoin(commitment.Denom, commitment.Amount))
	}

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

	return &types.MsgPerformActionJoinPoolResponse{
		ShareAmountOut: sharesOut,
	}, nil
}

func (k msgServer) PerformActionExitPool(goCtx context.Context, req *types.MsgPerformActionExitPool) (*types.MsgPerformActionExitPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrVaultNotFound, "vault %d not found", req.VaultId)
	}
	if vault.Manager != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "vault %d is not managed by %s", req.VaultId, req.Creator)
	}
	vaultAddress := types.NewVaultAddress(req.VaultId)

	exitCoins, weightBalanceBonus, slippage, _, _, err := k.amm.ExitPool(ctx, vaultAddress, req.PoolId, req.ShareAmountIn, req.MinAmountsOut, req.TokenOutDenom, false, true)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
	}

	// get coins after action
	coinsAfter := k.bk.GetAllBalances(ctx, vaultAddress)
	allCommitments := k.commitment.GetCommitments(ctx, vaultAddress)
	for _, commitment := range allCommitments.CommittedTokens {
		coinsAfter = coinsAfter.Add(sdk.NewCoin(commitment.Denom, commitment.Amount))
	}

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

	return &types.MsgPerformActionExitPoolResponse{
		TokenOut:           exitCoins,
		WeightBalanceRatio: weightBalanceBonus.Dec(),
		Slippage:           slippage.Dec(),
	}, nil
}

func (k msgServer) PerformActionSwapByDenom(goCtx context.Context, req *types.MsgPerformActionSwapByDenom) (*types.MsgPerformActionSwapByDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrVaultNotFound, "vault %d not found", req.VaultId)
	}
	if vault.Manager != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "vault %d is not managed by %s", req.VaultId, req.Creator)
	}
	vaultAddress := types.NewVaultAddress(req.VaultId)

	// denomOut should be in the allowed coins
	found = false
	for _, allowedCoin := range vault.AllowedCoins {
		if req.DenomOut == allowedCoin {
			found = true
			break
		}
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidAction, "vault %d does not allow this action", req.VaultId)
	}

	swapByDenom := &ammtypes.MsgSwapByDenom{
		Sender:    vaultAddress.String(),
		Amount:    req.Amount,
		MinAmount: req.MinAmount,
		MaxAmount: req.MaxAmount,
		DenomIn:   req.DenomIn,
		DenomOut:  req.DenomOut,
		Recipient: vaultAddress.String(),
	}

	swapCoins, err := k.amm.SwapByDenom(ctx, swapByDenom)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
	}

	return &types.MsgPerformActionSwapByDenomResponse{
		OutAmount: swapCoins.Amount,
	}, nil
}

// func (k Keeper) AllowedAction(ctx sdk.Context, action interface{}, vaultAddress string) bool {
// 	if action == nil {
// 		return false
// 	}

// 	switch perform_action := action.(type) {
// 	case *types.Action_JoinPool:
// 		// Verify join pool fields
// 		if perform_action.JoinPool == nil {
// 			return false
// 		}
// 		if perform_action.JoinPool.PoolId == 0 {
// 			return false
// 		}
// 		if len(perform_action.JoinPool.MaxAmountsIn) == 0 {
// 			return false
// 		}
// 		if perform_action.JoinPool.ShareAmountOut.IsNil() || perform_action.JoinPool.ShareAmountOut.IsZero() {
// 			return false
// 		}
// 		return true
// 	case *types.Action_ExitPool:
// 		// Verify exit pool fields
// 		if perform_action.ExitPool == nil {
// 			return false
// 		}
// 		if perform_action.ExitPool.PoolId == 0 {
// 			return false
// 		}
// 		if len(perform_action.ExitPool.MinAmountsOut) == 0 {
// 			return false
// 		}
// 		if perform_action.ExitPool.ShareAmountIn.IsNil() || perform_action.ExitPool.ShareAmountIn.IsZero() {
// 			return false
// 		}
// 		return true
// 	case *types.Action_SwapByDenom:
// 		// Verify swap by denom fields
// 		if perform_action.SwapByDenom == nil {
// 			return false
// 		}
// 		if perform_action.SwapByDenom.Amount.IsNil() || perform_action.SwapByDenom.Amount.IsZero() {
// 			return false
// 		}
// 		if perform_action.SwapByDenom.MinAmount.IsNil() || perform_action.SwapByDenom.MinAmount.IsZero() {
// 			return false
// 		}
// 		if perform_action.SwapByDenom.MaxAmount.IsNil() || perform_action.SwapByDenom.MaxAmount.IsZero() {
// 			return false
// 		}
// 		if perform_action.SwapByDenom.DenomIn == "" || perform_action.SwapByDenom.DenomOut == "" {
// 			return false
// 		}
// 		if perform_action.SwapByDenom.Recipient != vaultAddress {
// 			return false
// 		}
// 		return true
// 	case *types.Action_ClaimRewards:
// 		if perform_action.ClaimRewards == nil {
// 			return false
// 		}
// 		return true
// 	case *types.Action_CommitClaimedRewards:
// 		// Verify commit claimed rewards fields
// 		if perform_action.CommitClaimedRewards == nil {
// 			return false
// 		}
// 		if perform_action.CommitClaimedRewards.Denom == "" {
// 			return false
// 		}
// 		return true
// 	case *types.Action_UncommitTokens:
// 		// Verify uncommit tokens fields
// 		if perform_action.UncommitTokens == nil {
// 			return false
// 		}
// 		if perform_action.UncommitTokens.Amount.IsNil() || perform_action.UncommitTokens.Amount.IsZero() {
// 			return false
// 		}
// 		if perform_action.UncommitTokens.Denom == "" {
// 			return false
// 		}
// 		return true
// 	case *types.Action_Vest:
// 		// Verify vest fields
// 		if perform_action.Vest == nil {
// 			return false
// 		}
// 		if perform_action.Vest.Amount.IsNil() || perform_action.Vest.Amount.IsZero() {
// 			return false
// 		}
// 		if perform_action.Vest.Denom == "" {
// 			return false
// 		}
// 		return true
// 	case *types.Action_CancelVest:
// 		// Verify cancel vest fields
// 		if perform_action.CancelVest == nil {
// 			return false
// 		}
// 		return true
// 	case *types.Action_ClaimVesting:
// 		// Verify claim vesting fields
// 		if perform_action.ClaimVesting == nil {
// 			return false
// 		}
// 		return true
// 	default:
// 		return false
// 	}
// }
