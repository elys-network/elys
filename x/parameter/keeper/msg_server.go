package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/parameter/types"
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

func (k msgServer) UpdateMinCommission(goCtx context.Context, msg *types.MsgUpdateMinCommission) (*types.MsgUpdateMinCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	minComission, err := sdk.NewDecFromStr(msg.MinCommission)
	if err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)
	params.MinCommissionRate = minComission
	k.SetParams(ctx, params)
	return &types.MsgUpdateMinCommissionResponse{}, nil
}

func (k msgServer) UpdateMaxVotingPower(goCtx context.Context, msg *types.MsgUpdateMaxVotingPower) (*types.MsgUpdateMaxVotingPowerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	maxVotingPower, err := sdk.NewDecFromStr(msg.MaxVotingPower)
	if err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)
	params.MaxVotingPower = maxVotingPower
	k.SetParams(ctx, params)
	return &types.MsgUpdateMaxVotingPowerResponse{}, nil
}

func (k msgServer) UpdateMinSelfDelegation(goCtx context.Context, msg *types.MsgUpdateMinSelfDelegation) (*types.MsgUpdateMinSelfDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	minSelfDelegation, ok := sdk.NewIntFromString(msg.MinSelfDelegation)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	params := k.GetParams(ctx)
	params.MinSelfDelegation = minSelfDelegation
	k.SetParams(ctx, params)
	return &types.MsgUpdateMinSelfDelegationResponse{}, nil
}
