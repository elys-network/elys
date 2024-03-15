package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetPoolTotalSupply(poolId uint64) sdk.Int {
	// TODO
	return sdk.NewInt(0)
}

func (k Keeper) GetPoolBalance(poolId uint64, user string) sdk.Int {
	// TODO
	return sdk.NewInt(0)
}

func (k Keeper) UpdateAccPerShare(ctx sdk.Context, poolId uint64, rewardDenom string, amount sdk.Int) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)

	if !found {
		return
	}

	poolRewardInfo.PoolAccRewardPerShare = poolRewardInfo.PoolAccRewardPerShare.Add(
		math.LegacyDec(amount.Quo(k.GetPoolTotalSupply(poolId))),
	)
	poolRewardInfo.LastUpdatedBlock = uint64(ctx.BlockHeight())

	k.SetPoolRewardInfo(ctx, poolRewardInfo)
}

func (k Keeper) UpdateUserRewardPending(ctx sdk.Context, poolId uint64, rewardDenom string, user string, isDeposit bool, amount sdk.Int) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)

	if !found {
		return
	}

	userRewardInfo, found := k.GetUserRewardInfo(ctx, user, poolId, rewardDenom)

	if !found {
		return
	}

	// need to consider balance change on deposit/withdraw
	userBalance := k.GetPoolBalance(poolId, user)
	if isDeposit {
		userBalance = userBalance.Sub(amount)
	} else {
		userBalance = userBalance.Add(amount)
	}

	userRewardInfo.RewardPending = userRewardInfo.RewardPending.Add(
		poolRewardInfo.PoolAccRewardPerShare.Mul(
			math.LegacyDec(userBalance),
		).Sub(userRewardInfo.RewardDebt),
	)

	k.SetUserRewardInfo(ctx, userRewardInfo)
}

func (k Keeper) UpdateUserRewardDebt(ctx sdk.Context, poolId uint64, rewardDenom string, user string) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)

	if !found {
		return
	}

	userRewardInfo, found := k.GetUserRewardInfo(ctx, user, poolId, rewardDenom)

	if !found {
		return
	}

	userRewardInfo.RewardDebt = poolRewardInfo.PoolAccRewardPerShare.Mul(
		math.LegacyDec(k.GetPoolBalance(poolId, user)),
	)

	k.SetUserRewardInfo(ctx, userRewardInfo)
}
