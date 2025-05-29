package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

func (k Keeper) AfterDepositPerReward(ctx sdk.Context, poolId uint64, rewardDenom string, user sdk.AccAddress, amount math.Int) {
	k.UpdateUserRewardPending(
		ctx,
		poolId,
		rewardDenom,
		user,
		true,
		amount,
	)
	k.UpdateUserRewardDebt(
		ctx,
		poolId,
		rewardDenom,
		user,
	)
}

func (k Keeper) AfterDeposit(ctx sdk.Context, poolId uint64, user sdk.AccAddress, amount math.Int) {
	for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
		k.AfterDepositPerReward(ctx, poolId, rewardDenom, user, amount)
	}
}

func (k Keeper) AfterWithdrawPerReward(ctx sdk.Context, poolId uint64, rewardDenom string, user sdk.AccAddress, amount math.Int) {
	k.UpdateUserRewardPending(
		ctx,
		poolId,
		rewardDenom,
		user,
		false,
		amount,
	)
	k.UpdateUserRewardDebt(
		ctx,
		poolId,
		rewardDenom,
		user,
	)
}

func (k Keeper) AfterWithdraw(ctx sdk.Context, poolId uint64, user sdk.AccAddress, amount math.Int) {
	for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
		k.AfterWithdrawPerReward(ctx, poolId, rewardDenom, user, amount)
	}
}

func (k Keeper) GetRewardDenoms(ctx sdk.Context, poolId uint64) []string {
	rewardDenoms := make(map[string]bool)
	rewardDenoms[ptypes.Eden] = true
	rewardDenoms[k.GetBaseCurrencyDenom(ctx)] = true
	keys := []string{k.GetBaseCurrencyDenom(ctx)}

	poolInfo, found := k.GetPoolInfo(ctx, poolId)
	if !found {
		return []string{}
	}

	if poolInfo.EnableEdenRewards {
		keys = append(keys, ptypes.Eden)
	}

	for _, denom := range poolInfo.ExternalRewardDenoms {
		if rewardDenoms[denom] {
			continue
		}

		keys = append(keys, denom)
		rewardDenoms[denom] = true
	}

	return keys
}

func (k Keeper) GetBaseCurrencyDenom(ctx sdk.Context) string {
	baseCurrency, _ := k.assetProfileKeeper.GetUsdcDenom(ctx)
	return baseCurrency
}
