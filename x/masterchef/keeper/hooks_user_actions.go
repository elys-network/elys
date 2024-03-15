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
		ptypes.Eden,
		user,
	)
}

func (k Keeper) AfterDeposit(ctx sdk.Context, poolId uint64, user string, amount sdk.Int) {
	// Update Eden
	k.AfterDepositPerReward(ctx, poolId, ptypes.Eden, user, amount)

	// Update BaseCurrency
	k.AfterDepositPerReward(ctx, poolId, ptypes.BaseCurrency, user, amount)

	externalIncentives := k.GetAllExternalIncentives(ctx)

	for _, externalIncentive := range externalIncentives {
		if externalIncentive.PoolId == poolId {
			k.AfterDepositPerReward(ctx, poolId, externalIncentive.RewardDenom, user, amount)
		}
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
		ptypes.Eden,
		user,
	)
}

func (k Keeper) AfterWithdraw(ctx sdk.Context, poolId uint64, user string, amount sdk.Int) {
	// Update Eden
	k.AfterWithdrawPerReward(ctx, poolId, ptypes.Eden, user, amount)

	// Update BaseCurrency
	k.AfterWithdrawPerReward(ctx, poolId, ptypes.BaseCurrency, user, amount)

	externalIncentives := k.GetAllExternalIncentives(ctx)

	for _, externalIncentive := range externalIncentives {
		if externalIncentive.PoolId == poolId {
			k.AfterWithdrawPerReward(ctx, poolId, externalIncentive.RewardDenom, user, amount)
		}
	}
}
