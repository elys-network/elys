package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/incentive/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CommunityPool(goCtx context.Context, req *types.QueryCommunityPoolRequest) (*types.QueryCommunityPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pool := k.GetFeePoolCommunityCoins(ctx)

	return &types.QueryCommunityPoolResponse{Pool: pool}, nil
}

func (k Keeper) DelegationRewards(goCtx context.Context, req *types.QueryDelegationRewardsRequest) (*types.QueryDelegationRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	valAdr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	val := k.stk.Validator(ctx, valAdr)
	if val == nil {
		return nil, sdkerrors.Wrap(types.ErrNoValidatorExists, req.ValidatorAddress)
	}

	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	del := k.stk.Delegation(ctx, delAdr, valAdr)
	if del == nil {
		return nil, types.ErrNoDelegationExists
	}

	endingPeriod := k.IncrementValidatorPeriod(ctx, val)
	rewards := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	return &types.QueryDelegationRewardsResponse{Rewards: rewards}, nil
}

func (k Keeper) DelegationTotalRewards(goCtx context.Context, req *types.QueryDelegationTotalRewardsRequest) (*types.QueryDelegationTotalRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	total := sdk.DecCoins{}
	var delRewards []types.DelegationDelegatorReward

	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	k.stk.IterateDelegations(
		ctx, delAdr,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			valAddr := del.GetValidatorAddr()
			val := k.stk.Validator(ctx, valAddr)
			endingPeriod := k.IncrementValidatorPeriod(ctx, val)
			delReward := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

			delRewards = append(delRewards, types.NewDelegationDelegatorReward(valAddr, delReward))
			total = total.Add(delReward...)
			return false
		},
	)

	return &types.QueryDelegationTotalRewardsResponse{Rewards: delRewards, Total: total}, nil
}

func (k Keeper) DelegatorValidators(goCtx context.Context, req *types.QueryDelegatorValidatorsRequest) (*types.QueryDelegatorValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	var validators []string

	k.stk.IterateDelegations(
		ctx, delAdr,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			validators = append(validators, del.GetValidatorAddr().String())
			return false
		},
	)

	return &types.QueryDelegatorValidatorsResponse{Validators: validators}, nil
}

func (k Keeper) DelegatorWithdrawAddress(goCtx context.Context, req *types.QueryDelegatorWithdrawAddressRequest) (*types.QueryDelegatorWithdrawAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}
	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, delAdr)

	return &types.QueryDelegatorWithdrawAddressResponse{WithdrawAddress: withdrawAddr.String()}, nil
}

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) ValidatorCommission(goCtx context.Context, req *types.QueryValidatorCommissionRequest) (*types.QueryValidatorCommissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	valAdr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	commission := k.GetValidatorAccumulatedCommission(ctx, valAdr)

	return &types.QueryValidatorCommissionResponse{Commission: commission}, nil
}

func (k Keeper) ValidatorOutstandingRewards(goCtx context.Context, req *types.QueryValidatorOutstandingRewardsRequest) (*types.QueryValidatorOutstandingRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	valAdr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	rewards := k.GetValidatorOutstandingRewards(ctx, valAdr)

	return &types.QueryValidatorOutstandingRewardsResponse{Rewards: rewards}, nil
}

func (k Keeper) ValidatorSlashes(goCtx context.Context, req *types.QueryValidatorSlashesRequest) (*types.QueryValidatorSlashesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	if req.EndingHeight < req.StartingHeight {
		return nil, status.Errorf(codes.InvalidArgument, "starting height greater than ending height (%d > %d)", req.StartingHeight, req.EndingHeight)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address")
	}
	slashesStore := prefix.NewStore(store, types.GetValidatorSlashEventPrefix(valAddr))

	events, pageRes, err := query.GenericFilteredPaginate(k.cdc, slashesStore, req.Pagination, func(key []byte, result *types.ValidatorSlashEvent) (*types.ValidatorSlashEvent, error) {
		if result.ValidatorPeriod < req.StartingHeight || result.ValidatorPeriod > req.EndingHeight {
			return nil, nil
		}

		return result, nil
	}, func() *types.ValidatorSlashEvent {
		return &types.ValidatorSlashEvent{}
	})
	if err != nil {
		return nil, err
	}

	slashes := []types.ValidatorSlashEvent{}
	for _, event := range events {
		slashes = append(slashes, *event)
	}

	return &types.QueryValidatorSlashesResponse{Slashes: slashes, Pagination: pageRes}, nil
}
