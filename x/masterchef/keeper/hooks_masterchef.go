package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/masterchef/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) GetPoolTotalCommit(ctx sdk.Context, poolId uint64) sdk.Int {
	shareDenom := ammtypes.GetPoolShareDenom(poolId)
	if poolId == stablestaketypes.PoolId {
		shareDenom = stablestaketypes.GetShareDenom()
	}

	params := k.cmk.GetParams(ctx)
	return params.TotalCommitted.AmountOf(shareDenom)
}

func (k Keeper) GetPoolBalance(ctx sdk.Context, poolId uint64, user string) sdk.Int {
	commitments := k.cmk.GetCommitments(ctx, user)
	shareDenom := stablestaketypes.GetShareDenom()
	if poolId != stablestaketypes.PoolId {
		shareDenom = ammtypes.GetPoolShareDenom(poolId)
	}

	return commitments.GetCommittedAmountForDenom(shareDenom)
}

func (k Keeper) UpdateAccPerShare(ctx sdk.Context, poolId uint64, rewardDenom string, amount sdk.Int) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)

	if !found {
		poolRewardInfo = types.PoolRewardInfo{
			PoolId:                poolId,
			RewardDenom:           rewardDenom,
			PoolAccRewardPerShare: sdk.NewDec(0),
			LastUpdatedBlock:      0,
		}
	}

	totalCommit := k.GetPoolTotalCommit(ctx, poolId)
	if totalCommit.IsZero() {
		return
	}
	poolRewardInfo.PoolAccRewardPerShare = poolRewardInfo.PoolAccRewardPerShare.Add(
		math.LegacyNewDecFromInt(amount.Mul(ammtypes.OneShare)).
			Quo(math.LegacyNewDecFromInt(totalCommit)),
	)
	poolRewardInfo.LastUpdatedBlock = uint64(ctx.BlockHeight())
	k.SetPoolRewardInfo(ctx, poolRewardInfo)
}

func (k Keeper) UpdateUserRewardPending(ctx sdk.Context, poolId uint64, rewardDenom string, user string, isDeposit bool, amount sdk.Int) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)
	if !found {
		poolRewardInfo = types.PoolRewardInfo{
			PoolId:                poolId,
			RewardDenom:           rewardDenom,
			PoolAccRewardPerShare: sdk.NewDec(0),
			LastUpdatedBlock:      0,
		}
	}

	userRewardInfo, found := k.GetUserRewardInfo(ctx, user, poolId, rewardDenom)
	if !found {
		userRewardInfo = types.UserRewardInfo{
			User:          user,
			PoolId:        poolId,
			RewardDenom:   rewardDenom,
			RewardDebt:    sdk.NewDec(0),
			RewardPending: sdk.NewDec(0),
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
		poolRewardInfo.PoolAccRewardPerShare.
			Mul(math.LegacyNewDecFromInt(userBalance)).
			Sub(userRewardInfo.RewardDebt).
			Quo(math.LegacyNewDecFromInt(ammtypes.OneShare)),
	)

	k.SetUserRewardInfo(ctx, userRewardInfo)
}

func (k Keeper) UpdateUserRewardDebt(ctx sdk.Context, poolId uint64, rewardDenom string, user string) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)

	if !found {
		poolRewardInfo = types.PoolRewardInfo{
			PoolId:                poolId,
			RewardDenom:           rewardDenom,
			PoolAccRewardPerShare: sdk.NewDec(0),
			LastUpdatedBlock:      0,
		}
	}

	userRewardInfo, found := k.GetUserRewardInfo(ctx, user, poolId, rewardDenom)

	if !found {
		userRewardInfo = types.UserRewardInfo{
			User:          user,
			PoolId:        poolId,
			RewardDenom:   rewardDenom,
			RewardDebt:    sdk.NewDec(0),
			RewardPending: sdk.NewDec(0),
		}
	}

	userRewardInfo.RewardDebt = poolRewardInfo.PoolAccRewardPerShare.Mul(
		math.LegacyNewDecFromInt(k.GetPoolBalance(ctx, poolId, user)),
	)

	k.SetUserRewardInfo(ctx, userRewardInfo)
}
