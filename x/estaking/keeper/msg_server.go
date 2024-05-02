package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/estaking/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	k.SetParams(ctx, req.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) WithdrawReward(goCtx context.Context, msg *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delAddr := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	amount, err := k.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawRewardResponse{Amount: amount}, nil
}

func (k msgServer) WithdrawElysStakingRewards(goCtx context.Context, msg *types.MsgWithdrawElysStakingRewards) (*types.MsgWithdrawElysStakingRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delAddr := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)
	delegations := k.Keeper.Keeper.GetDelegatorDelegations(ctx, delAddr, 1024)
	rewards := sdk.Coins{}
	for _, del := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			return nil, err
		}
		amount, err := k.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
		if err != nil {
			return nil, err
		}
		rewards = rewards.Add(amount...)
	}

	return &types.MsgWithdrawElysStakingRewardsResponse{Amount: rewards}, nil
}

func (k Keeper) WithdrawAllRewards(goCtx context.Context, msg *types.MsgWithdrawAllRewards) (*types.MsgWithdrawAllRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delAddr := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)
	var amount sdk.Coins
	var err error = nil
	var rewards = sdk.Coins{}
	k.IterateDelegations(ctx, delAddr, func(index int64, del stakingtypes.DelegationI) (stop bool) {
		valAddr := del.GetValidatorAddr()
		amount, err = k.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
		if err != nil {
			return true
		}
		rewards = rewards.Add(amount...)
		return false
	})

	if err != nil {
		return nil, err
	}
	return &types.MsgWithdrawAllRewardsResponse{Amount: rewards}, nil
}
