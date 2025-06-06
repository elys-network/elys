package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v6/x/tokenomics/types"
)

func (k msgServer) CreateTimeBasedInflation(goCtx context.Context, msg *types.MsgCreateTimeBasedInflation) (*types.MsgCreateTimeBasedInflationResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, found := k.GetTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
	if found {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	timeBasedInflation := types.TimeBasedInflation{
		Authority:        msg.Authority,
		StartBlockHeight: msg.StartBlockHeight,
		EndBlockHeight:   msg.EndBlockHeight,
		Description:      msg.Description,
		Inflation:        msg.Inflation,
	}

	k.SetTimeBasedInflation(ctx, timeBasedInflation)
	return &types.MsgCreateTimeBasedInflationResponse{}, nil
}

func (k msgServer) UpdateTimeBasedInflation(goCtx context.Context, msg *types.MsgUpdateTimeBasedInflation) (*types.MsgUpdateTimeBasedInflationResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, found := k.GetTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	timeBasedInflation := types.TimeBasedInflation{
		Authority:        msg.Authority,
		StartBlockHeight: msg.StartBlockHeight,
		EndBlockHeight:   msg.EndBlockHeight,
		Description:      msg.Description,
		Inflation:        msg.Inflation,
	}

	k.SetTimeBasedInflation(ctx, timeBasedInflation)

	return &types.MsgUpdateTimeBasedInflationResponse{}, nil
}

func (k msgServer) DeleteTimeBasedInflation(goCtx context.Context, msg *types.MsgDeleteTimeBasedInflation) (*types.MsgDeleteTimeBasedInflationResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, found := k.GetTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
	return &types.MsgDeleteTimeBasedInflationResponse{}, nil
}
