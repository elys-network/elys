package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/vaults/types"
)

// func (k msgServer) Withdraw(goCtx context.Context, req *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(goCtx)
// 	creator := sdk.MustAccAddressFromBech32(req.Withdrawer)

// 	shareDenom := types.GetShareDenomForVault(req.VaultId)
// 	shareCoin := sdk.NewCoin(shareDenom, req.Shares)
// 	shareCoins := sdk.NewCoins(shareCoin)
// 	err := k.bk.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, shareCoins)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = k.bk.BurnCoins(ctx, types.ModuleName, shareCoins)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalShares := k.bk.GetSupply(ctx, shareDenom).Amount
// 	shareRatio := req.Shares.ToLegacyDec().Quo(totalShares.ToLegacyDec())

// 	vaultAddress := types.NewVaultAddress(req.VaultId)
// 	toSendCoins := sdk.NewCoins()
// 	vault, found := k.GetVault(ctx, req.VaultId)
// 	if !found {
// 		return nil, types.ErrVaultNotFound
// 	}
// 	for _, coin := range vault.AllowedCoins {
// 		balance := k.bk.GetBalance(ctx, vaultAddress, coin)
// 		amount := balance.Amount.ToLegacyDec().Mul(shareRatio).RoundInt()
// 		toSendCoins = toSendCoins.Add(sdk.NewCoin(coin, amount))
// 	}

// 	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, toSendCoins)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Commit tokens or not ?

// 	return &types.MsgWithdrawResponse{}, nil
// }

// Withdraw handles the withdrawal of funds from a vault
func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Deduct performance fee if applicable
	k.DeductPerformanceFee(ctx)

	// Get the vault
	vault, found := k.GetVault(ctx, msg.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}

	// Get user's share balance
	shareDenom := types.GetShareDenomForVault(vault.Id)
	userShareBalance := k.bk.GetBalance(ctx, sdk.MustAccAddressFromBech32(msg.Withdrawer), shareDenom)
	if userShareBalance.Amount.LT(msg.Shares) {
		return nil, types.ErrInsufficientShares
	}

	// Calculate user's share ratio
	totalShares := k.bk.GetSupply(ctx, shareDenom).Amount
	shareRatio := sdkmath.LegacyNewDecFromInt(userShareBalance.Amount).Quo(sdkmath.LegacyNewDecFromInt(totalShares))

	// Burn user's shares
	err := k.bk.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(shareDenom, msg.Shares)))
	if err != nil {
		return nil, err
	}

	balancesBefore := k.bk.GetAllBalances(ctx, types.NewVaultAddress(vault.Id))

	// Execute withdrawal strategy if it exists
	if vault.WithdrawStrategy != nil {
		for _, action := range vault.WithdrawStrategy {
			// Scale action amounts based on user's share ratio
			scaledAction := scaleActionAmounts(action, shareRatio)

			// Create perform action message
			performActionMsg := &types.MsgPerformAction{
				Creator: vault.Manager,
				Action:  scaledAction,
				VaultId: msg.VaultId,
			}

			// Execute the action
			_, err := k.PerformAction(goCtx, performActionMsg)
			if err != nil {
				return nil, fmt.Errorf("failed to execute withdrawal strategy action: %w", err)
			}
		}
	}

	balancesAfter := k.bk.GetAllBalances(ctx, types.NewVaultAddress(vault.Id))
	amountsToSend := balancesAfter.Sub(balancesBefore...)

	// Send coins to user
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(msg.Withdrawer), amountsToSend)
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawResponse{
		VaultId: msg.VaultId,
		Amount:  amountsToSend,
	}, nil
}

// scaleActionAmounts scales the amounts in an action based on the share ratio
func scaleActionAmounts(action *types.Action, shareRatio sdkmath.LegacyDec) *types.Action {
	scaledAction := *action // Create a copy of the action

	switch perform_action := action.Action.(type) {
	case *types.Action_ExitPool:
		scaledAction.Action = &types.Action_ExitPool{
			ExitPool: &ammtypes.MsgExitPool{
				Sender:        perform_action.ExitPool.Sender,
				PoolId:        perform_action.ExitPool.PoolId,
				MinAmountsOut: perform_action.ExitPool.MinAmountsOut,
				ShareAmountIn: perform_action.ExitPool.ShareAmountIn.ToLegacyDec().Mul(shareRatio).TruncateInt(),
			},
		}
	}

	// TODO: What if multiple exit from a lp shares

	return &scaledAction
}
