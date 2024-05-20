package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/estaking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Rewards(goCtx context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	fmt.Println("DEBUG:: EstakingRewards-1")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	fmt.Println("DEBUG:: EstakingRewards-2")
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	fmt.Println("DEBUG:: EstakingRewards-3")
	total := sdk.DecCoins{}
	var delRewards []types.DelegationDelegatorReward

	delAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	fmt.Println("DEBUG:: EstakingRewards-4")
	k.IterateDelegations(
		ctx, delAddr,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			fmt.Println("DEBUG:: EstakingRewards-4-1")
			valAddr := del.GetValidatorAddr()
			val := k.Validator(ctx, valAddr)
			fmt.Println("DEBUG:: EstakingRewards-4-2", val)
			endingPeriod := k.distrKeeper.IncrementValidatorPeriod(ctx, val)
			fmt.Println("DEBUG:: EstakingRewards-4-3", del)
			delReward := k.distrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)

			fmt.Println("DEBUG:: EstakingRewards-4-4", delReward)
			finalRewards, _ := delReward.TruncateDecimal()
			if finalRewards == nil {
				finalRewards = []sdk.Coin{}
			}
			fmt.Println("DEBUG:: EstakingRewards-4-5", finalRewards)
			delRewards = append(delRewards, types.DelegationDelegatorReward{
				ValidatorAddress: valAddr.String(),
				Reward:           finalRewards,
			})
			total = total.Add(delReward...)
			fmt.Println("DEBUG:: EstakingRewards-4-6", total)
			return false
		},
	)
	finalTotalRewards, _ := total.TruncateDecimal()
	fmt.Println("DEBUG:: EstakingRewards-5")

	return &types.QueryRewardsResponse{Rewards: delRewards, Total: finalTotalRewards}, nil
}
