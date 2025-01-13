package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/masterchef/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) GetPoolTotalCommit(ctx sdk.Context, poolId uint64) math.Int {
	var shareDenom string
	if poolId >= stablestaketypes.PoolId {
		shareDenom = stablestaketypes.GetShareDenomForPool(poolId)
	} else {
		shareDenom = ammtypes.GetPoolShareDenom(poolId)
	}

	params := k.commitmentKeeper.GetParams(ctx)
	return params.TotalCommitted.AmountOf(shareDenom)
}

func (k Keeper) GetPoolBalance(ctx sdk.Context, poolId uint64, user sdk.AccAddress) math.Int {
	commitments := k.commitmentKeeper.GetCommitments(ctx, user)
	var shareDenom string
	if poolId >= stablestaketypes.PoolId {
		shareDenom = stablestaketypes.GetShareDenomForPool(poolId)
	} else {
		shareDenom = ammtypes.GetPoolShareDenom(poolId)
	}

	return commitments.GetCommittedAmountForDenom(shareDenom)
}

func (k Keeper) UpdateAccPerShare(ctx sdk.Context, poolId uint64, rewardDenom string, amount math.Int) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)
	if !found {
		poolRewardInfo = types.PoolRewardInfo{
			PoolId:                poolId,
			RewardDenom:           rewardDenom,
			PoolAccRewardPerShare: math.LegacyNewDec(0),
			LastUpdatedBlock:      uint64(ctx.BlockHeight()),
		}
	}

	totalCommit := k.GetPoolTotalCommit(ctx, poolId)
	if totalCommit.IsZero() {
		return
	}
	poolRewardInfo.PoolAccRewardPerShare = poolRewardInfo.PoolAccRewardPerShare.Add(
		math.LegacyNewDecFromInt(amount.Mul(ammtypes.OneShare)).QuoInt(totalCommit),
	)
	poolRewardInfo.LastUpdatedBlock = uint64(ctx.BlockHeight())
	k.SetPoolRewardInfo(ctx, poolRewardInfo)
}

func (k Keeper) UpdateUserRewardPending(ctx sdk.Context, poolId uint64, rewardDenom string, user sdk.AccAddress, isDeposit bool, amount math.Int) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)
	poolAccRewardPerShare := math.LegacyNewDec(0)
	if found {
		poolAccRewardPerShare = poolRewardInfo.PoolAccRewardPerShare
	}

	userRewardInfo, found := k.GetUserRewardInfo(ctx, user, poolId, rewardDenom)
	if !found {
		userRewardInfo = types.UserRewardInfo{
			User:          user.String(),
			PoolId:        poolId,
			RewardDenom:   rewardDenom,
			RewardDebt:    math.LegacyNewDec(0),
			RewardPending: math.LegacyNewDec(0),
		}
	}

	// need to consider balance change on deposit/withdraw
	userBalance := k.GetPoolBalance(ctx, poolId, user)
	if isDeposit {
		userBalance = userBalance.Sub(amount)
	} else {
		userBalance = userBalance.Add(amount)
	}

	userRewardInfo.RewardPending = userRewardInfo.RewardPending.Add(
		poolAccRewardPerShare.
			MulInt(userBalance).
			Sub(userRewardInfo.RewardDebt).
			QuoInt(ammtypes.OneShare),
	)

	k.SetUserRewardInfo(ctx, userRewardInfo)
}

func (k Keeper) UpdateUserRewardDebt(ctx sdk.Context, poolId uint64, rewardDenom string, user sdk.AccAddress) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)
	poolAccRewardPerShare := math.LegacyNewDec(0)
	if found {
		poolAccRewardPerShare = poolRewardInfo.PoolAccRewardPerShare
	}

	userRewardInfo, found := k.GetUserRewardInfo(ctx, user, poolId, rewardDenom)
	if !found {
		userRewardInfo = types.UserRewardInfo{
			User:          user.String(),
			PoolId:        poolId,
			RewardDenom:   rewardDenom,
			RewardDebt:    math.LegacyNewDec(0),
			RewardPending: math.LegacyNewDec(0),
		}
	}

	userRewardInfo.RewardDebt = poolAccRewardPerShare.Mul(
		math.LegacyNewDecFromInt(k.GetPoolBalance(ctx, poolId, user)),
	)

	k.SetUserRewardInfo(ctx, userRewardInfo)
}
