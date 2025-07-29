package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/x/parameter/types"
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

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.MinCommissionRate = msg.MinCommission
	k.SetParams(ctx, params)
	return &types.MsgUpdateMinCommissionResponse{}, nil
}

func (k msgServer) UpdateMaxVotingPower(goCtx context.Context, msg *types.MsgUpdateMaxVotingPower) (*types.MsgUpdateMaxVotingPowerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.MaxVotingPower = msg.MaxVotingPower
	k.SetParams(ctx, params)
	return &types.MsgUpdateMaxVotingPowerResponse{}, nil
}

func (k msgServer) UpdateMinSelfDelegation(goCtx context.Context, msg *types.MsgUpdateMinSelfDelegation) (*types.MsgUpdateMinSelfDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.MinSelfDelegation = msg.MinSelfDelegation
	k.SetParams(ctx, params)
	return &types.MsgUpdateMinSelfDelegationResponse{}, nil
}

func (k msgServer) UpdateTotalBlocksPerYear(goCtx context.Context, msg *types.MsgUpdateTotalBlocksPerYear) (*types.MsgUpdateTotalBlocksPerYearResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.TotalBlocksPerYear = msg.TotalBlocksPerYear
	k.SetParams(ctx, params)
	return &types.MsgUpdateTotalBlocksPerYearResponse{}, nil
}

func (k msgServer) UpdateTakerFees(goCtx context.Context, msg *types.MsgUpdateTakerFees) (*types.MsgUpdateTakerFeesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.TakerFees = msg.TakerFees
	params.EnableTakerFeeSwap = msg.EnableTakerFeeSwap
	params.TakerFeeCollectionInterval = msg.TakerFeeCollectionInterval
	k.SetParams(ctx, params)
	return &types.MsgUpdateTakerFeesResponse{}, nil
}
