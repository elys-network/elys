package keeper

import (
	"context"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/vaults/types"
)

func (k Keeper) ClaimRewards(ctx sdk.Context, sender sdk.AccAddress, poolIds []uint64, recipient sdk.AccAddress) error {
	coins := sdk.NewCoins()
	rewardPoolIds := []string{}
	for _, poolId := range poolIds {
		k.AfterWithdraw(ctx, poolId, sender, math.ZeroInt())

		for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
			userRewardInfo, found := k.GetUserRewardInfo(ctx, sender, poolId, rewardDenom)
			if found && userRewardInfo.RewardPending.IsPositive() {
				coin := sdk.NewCoin(rewardDenom, userRewardInfo.RewardPending.TruncateInt())
				coins = coins.Add(coin)
				rewardPoolIds = append(rewardPoolIds, strconv.FormatUint(poolId, 10))

				userRewardInfo.RewardPending = math.LegacyZeroDec()
				if userRewardInfo.RewardDebt.IsZero() {
					k.RemoveUserRewardInfo(ctx, userRewardInfo.GetUserAccount(), userRewardInfo.PoolId, userRewardInfo.RewardDenom)
				} else {
					k.SetUserRewardInfo(ctx, userRewardInfo)
				}
			}
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtClaimRewards,
			sdk.NewAttribute(types.AttributeSender, sender.String()),
			sdk.NewAttribute(types.AttributeRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeVaultIds, strings.Join(rewardPoolIds, ",")),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
		),
	})

	// Transfer rewards (Eden/EdenB is transferred through commitment module)
	err := k.commitment.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
	if err != nil {
		return err
	}

	return nil
}

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	if len(msg.VaultIds) == 0 {
		allPools := k.GetAllPoolInfos(ctx)
		for _, pool := range allPools {
			msg.VaultIds = append(msg.VaultIds, pool.PoolId)
		}
	}

	err := k.Keeper.ClaimRewards(ctx, sender, msg.VaultIds, sender)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
