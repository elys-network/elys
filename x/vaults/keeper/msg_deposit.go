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
	if err := k.SetParams(ctx, req.Params); err != nil {
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
