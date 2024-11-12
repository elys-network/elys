package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	minComission, err := sdkmath.LegacyNewDecFromStr(msg.MinCommission)
	if err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)
	params.MinCommissionRate = minComission
	k.SetParams(ctx, params)
	return &types.MsgUpdateMinCommissionResponse{}, nil
}

func (k msgServer) UpdateMaxVotingPower(goCtx context.Context, msg *types.MsgUpdateMaxVotingPower) (*types.MsgUpdateMaxVotingPowerResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	maxVotingPower, err := sdkmath.LegacyNewDecFromStr(msg.MaxVotingPower)
	if err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)
	params.MaxVotingPower = maxVotingPower
	k.SetParams(ctx, params)
	return &types.MsgUpdateMaxVotingPowerResponse{}, nil
}

func (k msgServer) UpdateMinSelfDelegation(goCtx context.Context, msg *types.MsgUpdateMinSelfDelegation) (*types.MsgUpdateMinSelfDelegationResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	minSelfDelegation, ok := sdkmath.NewIntFromString(msg.MinSelfDelegation)
	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	params := k.GetParams(ctx)
	params.MinSelfDelegation = minSelfDelegation
	k.SetParams(ctx, params)
	return &types.MsgUpdateMinSelfDelegationResponse{}, nil
}

func (k msgServer) UpdateBrokerAddress(goCtx context.Context, msg *types.MsgUpdateBrokerAddress) (*types.MsgUpdateBrokerAddressResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.BrokerAddress = msg.BrokerAddress
	k.SetParams(ctx, params)
	return &types.MsgUpdateBrokerAddressResponse{}, nil
}

func (k msgServer) UpdateTotalBlocksPerYear(goCtx context.Context, msg *types.MsgUpdateTotalBlocksPerYear) (*types.MsgUpdateTotalBlocksPerYearResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.TotalBlocksPerYear = msg.TotalBlocksPerYear
	k.SetParams(ctx, params)
	return &types.MsgUpdateTotalBlocksPerYearResponse{}, nil
}
