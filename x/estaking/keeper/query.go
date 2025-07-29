package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/v7/x/estaking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Rewards(goCtx context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	total := sdk.DecCoins{}
	var delRewards []types.DelegationDelegatorReward

	delAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	err = k.IterateDelegations(
		ctx, delAddr,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			valAddr, err := sdk.ValAddressFromBech32(del.GetValidatorAddr())
			if err != nil {
				panic(err)
			}
			val, err := k.Validator(ctx, valAddr)
			if val == nil || err != nil {
				return false
			}
			endingPeriod, err := k.distrKeeper.IncrementValidatorPeriod(ctx, val)
			if err != nil {
				return false
			}
			// check if the delegation is valid
			startingInfo, err := k.distrKeeper.GetDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(delAddr))
			if err != nil {
				return false
			}
			if startingInfo.Stake.IsNil() {
				return false
			}
			delReward, err := k.distrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
			if err != nil {
				return false
			}

			finalRewards, _ := delReward.TruncateDecimal()
			if finalRewards == nil {
				finalRewards = []sdk.Coin{}
			}
			delRewards = append(delRewards, types.DelegationDelegatorReward{
				ValidatorAddress: valAddr.String(),
				Reward:           finalRewards,
			})
			total = total.Add(delReward...)
			return false
		},
	)
	if err != nil {
		return nil, err
	}
	finalTotalRewards, _ := total.TruncateDecimal()

	return &types.QueryRewardsResponse{Rewards: delRewards, Total: finalTotalRewards}, nil
}

func (k Keeper) Invariant(goCtx context.Context, req *types.QueryInvariantRequest) (*types.QueryInvariantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	valTokensSum := math.ZeroInt()
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		valTokensSum = valTokensSum.Add(validator.GetTokens())
		return false
	})
	if err != nil {
		return nil, err
	}

	totalBondedTokens, err := k.TotalBondedTokens(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryInvariantResponse{
		TotalBonded:              totalBondedTokens,
		BondedValidatorTokensSum: valTokensSum,
	}, nil
}
