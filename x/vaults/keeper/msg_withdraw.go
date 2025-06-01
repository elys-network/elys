package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/vaults/types"
)

func (k msgServer) Withdraw(goCtx context.Context, req *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	creator := sdk.MustAccAddressFromBech32(req.Withdrawer)

	k.DeductPerformanceFee(ctx)

	shareDenom := types.GetShareDenomForVault(req.VaultId)
	shareCoin := sdk.NewCoin(shareDenom, req.Shares)
	shareCoins := sdk.NewCoins(shareCoin)

	err := k.commitment.UncommitTokens(ctx, creator, shareDenom, req.Shares, false)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	totalShares := k.bk.GetSupply(ctx, shareDenom).Amount
	shareRatio := req.Shares.ToLegacyDec().Quo(totalShares.ToLegacyDec())

	vaultAddress := types.NewVaultAddress(req.VaultId)
	toSendCoins := sdk.NewCoins()
	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}
	commitments := k.commitment.GetCommitments(ctx, creator)

	for _, commitment := range commitments.CommittedTokens {
		amount := commitment.Amount.ToLegacyDec().Mul(shareRatio).RoundInt()
		toSendCoins = toSendCoins.Add(sdk.NewCoin(commitment.Denom, amount))
	}

	for _, coin := range vault.AllowedCoins {
		balance := k.bk.GetBalance(ctx, vaultAddress, coin)
		amount := balance.Amount.ToLegacyDec().Mul(shareRatio).RoundInt()
		toSendCoins = toSendCoins.Add(sdk.NewCoin(coin, amount))
	}

	err = k.bk.SendCoins(ctx, vaultAddress, creator, toSendCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.BurnCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawResponse{}, nil
}

// TODO: Add withdraw fee for amount received
