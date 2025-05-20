package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	"github.com/elys-network/elys/v4/x/masterchef/types"
	stablestaketypes "github.com/elys-network/elys/v4/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetPoolTotalCommit(ctx sdk.Context, poolId uint64) math.Int {
	var shareDenom string
	if poolId >= stablestaketypes.UsdcPoolId {
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
	if poolId >= stablestaketypes.UsdcPoolId {
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
	poolRewardInfo.PoolAccRewardPerShare = poolRewardInfo.GetBigDecPoolAccRewardPerShare().Add(
		(osmomath.BigDecFromSDKInt(amount).Mul(ammtypes.OneShareBigDec)).Quo(osmomath.BigDecFromSDKInt(totalCommit)),
	).Dec()
	poolRewardInfo.LastUpdatedBlock = uint64(ctx.BlockHeight())
	k.SetPoolRewardInfo(ctx, poolRewardInfo)
}

func (k Keeper) UpdateUserRewardPending(ctx sdk.Context, poolId uint64, rewardDenom string, user sdk.AccAddress, isDeposit bool, amount math.Int) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)
	poolAccRewardPerShare := osmomath.ZeroBigDec()
	if found {
		poolAccRewardPerShare = poolRewardInfo.GetBigDecPoolAccRewardPerShare()
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

	userRewardInfo.RewardPending = userRewardInfo.GetBigDecRewardPending().Add(
		poolAccRewardPerShare.
			Mul(osmomath.BigDecFromSDKInt(userBalance)).
			Sub(userRewardInfo.GetBigDecRewardDebt()).
			Quo(ammtypes.OneShareBigDec),
	).Dec()

	k.SetUserRewardInfo(ctx, userRewardInfo)
}

func (k Keeper) UpdateUserRewardDebt(ctx sdk.Context, poolId uint64, rewardDenom string, user sdk.AccAddress) {
	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)
	poolAccRewardPerShare := osmomath.ZeroBigDec()
	if found {
		poolAccRewardPerShare = poolRewardInfo.GetBigDecPoolAccRewardPerShare()
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
		osmomath.BigDecFromSDKInt(k.GetPoolBalance(ctx, poolId, user)),
	).Dec()

	k.SetUserRewardInfo(ctx, userRewardInfo)
}
