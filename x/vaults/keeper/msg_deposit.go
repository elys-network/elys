package keeper

import (
	"context"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	tiertypes "github.com/elys-network/elys/x/tier/types"

	"github.com/elys-network/elys/x/vaults/types"
)

func (k msgServer) Deposit(goCtx context.Context, req *types.MsgDeposit) (*types.MsgDepositResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return &types.MsgDepositResponse{}, types.ErrVaultNotFound
	}

	depositer := sdk.MustAccAddressFromBech32(req.Depositor)
	redemptionRate := k.CalculateRedemptionRateForVault(ctx, vault.Id)
	vaultName := types.GetVaultIdModuleName(vault.Id)

	depositCoin := sdk.NewCoin(vault.DepositDenom, req.Amount.Amount)
	err := k.bk.SendCoinsFromAccountToModule(ctx, depositer, vaultName, sdk.Coins{depositCoin})
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetShareDenomForVault(vault.Id)
	// Initial case
	if redemptionRate.IsZero() {
		redemptionRate = sdkmath.LegacyOneDec()
	}
	shareAmount := depositCoin.Amount.ToLegacyDec().Quo(redemptionRate).RoundInt()
	shareCoins := sdk.NewCoins(sdk.NewCoin(shareDenom, shareAmount))

	err = k.bk.MintCoins(ctx, vaultName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, vaultName, depositer, shareCoins)
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositResponse{}, nil
}

func (k Keeper) VaultUsdValue(ctx sdk.Context, vaultId uint64) (sdkmath.LegacyDec, error) {
	vault, found := k.GetVault(ctx, vaultId)
	if !found {
		return sdkmath.LegacyZeroDec(), types.ErrVaultNotFound
	}
	vaultAddress := types.NewVaultAddress(vaultId)
	totalValue := sdkmath.LegacyZeroDec()
	commitments := k.commitement.GetCommitments(ctx, vaultAddress)
	// TODO: Handle zero values for denom, we should issue shares if price is not available
	for _, commitment := range commitments.CommittedTokens {
		// Pool balance
		if strings.HasPrefix(commitment.Denom, "amm/pool") {
			poolId, err := ammtypes.GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				continue
			}
			pool, found := k.amm.GetPool(ctx, poolId)
			if !found {
				continue
			}
			info := k.amm.PoolExtraInfo(ctx, pool, tiertypes.OneDay)
			amount := commitment.Amount.ToLegacyDec()
			totalValue = totalValue.Add(amount.Mul(info.LpTokenPrice).QuoInt(ammtypes.OneShare))
		}
	}
	for _, coin := range vault.AllowedCoins {
		if !strings.HasPrefix(coin, "amm/pool") {
			balance := k.bk.GetBalance(ctx, vaultAddress, coin)
			totalValue = totalValue.Add(k.amm.CalculateUSDValue(ctx, coin, balance.Amount))
		}
	}
	for _, coin := range vault.RewardCoins {
		if !strings.HasPrefix(coin, "amm/pool") {
			balance := k.bk.GetBalance(ctx, vaultAddress, coin)
			totalValue = totalValue.Add(k.amm.CalculateUSDValue(ctx, coin, balance.Amount))
		}
	}
	return totalValue, nil
}

func (k Keeper) CalculateRedemptionRateForVault(ctx sdk.Context, vaultId uint64) sdkmath.LegacyDec {
	totalShares := k.bk.GetSupply(ctx, types.GetShareDenomForVault(vaultId))

	if totalShares.Amount.IsZero() {
		return sdkmath.LegacyZeroDec()
	}

	// TODO: Handle zero values for denom, we should not issue shares if price is not available
	usdValue, err := k.VaultUsdValue(ctx, vaultId)
	if err != nil {
		return sdkmath.LegacyZeroDec()
	}

	return usdValue.Quo(totalShares.Amount.ToLegacyDec())
}
