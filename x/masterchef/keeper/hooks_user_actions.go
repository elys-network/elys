package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) AfterDepositPerReward(ctx sdk.Context, poolId uint64, rewardDenom string, user string, amount sdk.Int) {
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

func (k Keeper) AfterDeposit(ctx sdk.Context, poolId uint64, user string, amount sdk.Int) {
	for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
		k.AfterDepositPerReward(ctx, poolId, rewardDenom, user, amount)
	}
}

func (k Keeper) AfterWithdrawPerReward(ctx sdk.Context, poolId uint64, rewardDenom string, user string, amount sdk.Int) {
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

func (k Keeper) AfterWithdraw(ctx sdk.Context, poolId uint64, user string, amount sdk.Int) {
	for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
		k.AfterWithdrawPerReward(ctx, poolId, rewardDenom, user, amount)
	}
}

func (k Keeper) GetRewardDenoms(ctx sdk.Context, poolId uint64) []string {
	rewardDenoms := make(map[string]bool)

	rewardDenoms[ptypes.Eden] = true
	rewardDenoms[ptypes.BaseCurrency] = true

	poolInfo, found := k.GetPool(ctx, poolId)
	if !found {
		return []string{}
	}

	externalIncentives := poolInfo.ExternalRewardDenoms
	for _, externalIncentive := range externalIncentives {
		rewardDenoms[externalIncentive] = true
	}

	keys := make([]string, 0, len(rewardDenoms))
	for k2 := range rewardDenoms {
		keys = append(keys, k2)
	}

	return keys
}
