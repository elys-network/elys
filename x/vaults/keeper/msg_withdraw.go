package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

func (k msgServer) Withdraw(goCtx context.Context, req *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//creator := sdk.MustAccAddressFromBech32(req.Withdrawer)

	k.DeductPerformanceFee(ctx)

	return &types.MsgWithdrawResponse{}, nil
}
