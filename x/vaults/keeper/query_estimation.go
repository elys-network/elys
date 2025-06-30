package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"

	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (k Keeper) DepositEstimation(goCtx context.Context, req *types.QueryDepositEstimationRequest) (*types.QueryDepositEstimationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}

	redemptionRate := k.CalculateRedemptionRateForVault(ctx, vault.Id)
	var usdValue osmomath.BigDec
	var shareAmount sdkmath.Int
	if redemptionRate.IsZero() {
		usdValue = k.amm.CalculateUSDValue(ctx, vault.DepositDenom, req.Amount)
		if usdValue.IsZero() {
			return nil, types.ErrDepositValueZero
		}
		shareAmount = usdValue.Mul(osmomath.BigDecFromSDKInt(sdkmath.NewInt(1000000))).Dec().RoundInt()
	} else {
		shareAmount = (osmomath.BigDecFromSDKInt(req.Amount).Quo(redemptionRate)).Dec().RoundInt()
	}

	return &types.QueryDepositEstimationResponse{
		SharesAmount:   shareAmount,
		SharesUsdValue: usdValue.Dec(),
	}, nil
}

func (k Keeper) WithdrawEstimation(goCtx context.Context, req *types.QueryWithdrawEstimationRequest) (*types.QueryWithdrawEstimationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}

	vaultAddress := types.NewVaultAddress(vault.Id)
	shareDenom := types.GetShareDenomForVault(vault.Id)
	totalShares := k.bk.GetSupply(ctx, shareDenom).Amount
	shareRatio := req.SharesAmount.ToLegacyDec().Quo(totalShares.ToLegacyDec())

	toSendCoins := sdk.NewCoins()
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
			shareCoins, _, _, _, _, err := k.amm.ExitPool(ctx, vaultAddress, poolId, coin.Amount, sdk.Coins{}, coin.Denom, false, true)
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
			}

			toSendCoins = toSendCoins.Sub(coin)
			toSendCoins = toSendCoins.Add(shareCoins...)
		}
	}

	// Convert to coins value to usd value
	// TODO: Maybe use direct swap for accurate usd value
	usdValue := osmomath.ZeroBigDec()
	for _, coin := range toSendCoins {
		usdValue = usdValue.Add(k.amm.CalculateUSDValue(ctx, coin.Denom, coin.Amount))
	}

	return &types.QueryWithdrawEstimationResponse{
		Amount:   usdValue.Dec().TruncateInt(),
		UsdValue: usdValue.Dec(),
	}, nil
}
