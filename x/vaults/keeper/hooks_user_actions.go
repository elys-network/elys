package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
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
	for _, rewardDenom := range k.GetRewardDenoms(ctx) {
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
	for _, rewardDenom := range k.GetRewardDenoms(ctx) {
		k.AfterWithdrawPerReward(ctx, poolId, rewardDenom, user, amount)
	}
}

// TODO: Extend this to support external reward denoms
func (k Keeper) GetRewardDenoms(ctx sdk.Context) []string {
	return []string{ptypes.Eden, k.GetBaseCurrencyDenom(ctx)}
}

func (k Keeper) GetBaseCurrencyDenom(ctx sdk.Context) string {
	baseCurrency, _ := k.assetProfileKeeper.GetUsdcDenom(ctx)
	return baseCurrency
}
