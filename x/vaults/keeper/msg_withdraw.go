package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/vaults/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k msgServer) Withdraw(goCtx context.Context, req *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	creator := sdk.MustAccAddressFromBech32(req.Withdrawer)
	vaultAddress := types.NewVaultAddress(req.VaultId)

	_, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}

	k.DeductPerformanceFee(ctx)

	// claim pending rewards
	k.ClaimRewards(ctx, &types.MsgClaimRewards{
		Sender:   creator.String(),
		VaultIds: k.GetAllPoolIds(ctx, vaultAddress),
	})

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

	err = k.bk.BurnCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	if totalShares.IsZero() {
		return nil, types.ErrNoShares
	}

	shareRatio := req.Shares.ToLegacyDec().Quo(totalShares.ToLegacyDec())

	toSendCoins := sdk.NewCoins()
	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}
	commitments := k.commitment.GetCommitments(ctx, vaultAddress)

	for _, commitment := range commitments.CommittedTokens {
		amount := commitment.Amount.ToLegacyDec().Mul(shareRatio).RoundInt()
		toSendCoins = toSendCoins.Add(sdk.NewCoin(commitment.Denom, amount))
	}

	for _, coin := range vault.AllowedCoins {
		balance := k.bk.GetBalance(ctx, vaultAddress, coin)
		amount := balance.Amount.ToLegacyDec().Mul(shareRatio).RoundInt()
		toSendCoins = toSendCoins.Add(sdk.NewCoin(coin, amount))
	}

	for _, coin := range toSendCoins {
		// FOR AMM LP
		if strings.HasPrefix(coin.Denom, "amm/pool/") {
			poolId, err := GetPoolIdFromShareDenom(coin.Denom)
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
			}
			vaultAddress := types.NewVaultAddress(poolId)

			// exit pool
			shareCoins, _, _, _, _, err = k.amm.ExitPool(ctx, vaultAddress, poolId, coin.Amount, sdk.Coins{}, coin.Denom, false, true)
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
			}

			toSendCoins = toSendCoins.Sub(coin)
			toSendCoins = toSendCoins.Add(shareCoins...)
		}
	}

	err = k.bk.SendCoins(ctx, vaultAddress, creator, toSendCoins)
	if err != nil {
		return nil, err
	}

	// Set withdrawal usd value
	usdValue, err := k.VaultUsdValue(ctx, req.VaultId)
	if err != nil {
		return nil, err
	}
	usdValue = osmomath.BigDecFromDec(usdValue.Dec().Mul(shareRatio))
	vault.WithdrawalUsdValue = vault.WithdrawalUsdValue.Add(usdValue.Dec())
	k.SetVault(ctx, vault)

	k.AfterWithdraw(ctx, req.VaultId, creator, req.Shares)

	return &types.MsgWithdrawResponse{}, nil
}
