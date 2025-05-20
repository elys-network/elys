package keeper

import (
	"context"
	"cosmossdk.io/math"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/elys-network/elys/v4/x/tokenomics/types"
)

func (k msgServer) CreateAirdrop(goCtx context.Context, msg *types.MsgCreateAirdrop) (*types.MsgCreateAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, found := k.GetAirdrop(ctx, msg.Intent)
	if found {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	airdrop := types.Airdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
		Amount:    msg.Amount,
		Expiry:    msg.Expiry,
	}

	k.SetAirdrop(ctx, airdrop)
	return &types.MsgCreateAirdropResponse{}, nil
}

func (k msgServer) UpdateAirdrop(goCtx context.Context, msg *types.MsgUpdateAirdrop) (*types.MsgUpdateAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, found := k.GetAirdrop(ctx, msg.Intent)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	airdrop := types.Airdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
		Amount:    msg.Amount,
		Expiry:    msg.Expiry,
	}

	k.SetAirdrop(ctx, airdrop)
	return &types.MsgUpdateAirdropResponse{}, nil
}

func (k msgServer) DeleteAirdrop(goCtx context.Context, msg *types.MsgDeleteAirdrop) (*types.MsgDeleteAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, found := k.GetAirdrop(ctx, msg.Intent)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveAirdrop(ctx, msg.Intent)
	return &types.MsgDeleteAirdropResponse{}, nil
}

func (k msgServer) ClaimAirdrop(goCtx context.Context, msg *types.MsgClaimAirdrop) (*types.MsgClaimAirdropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	airdrop, found := k.GetAirdrop(ctx, msg.Sender)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Sender != airdrop.Authority {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if ctx.BlockTime().Unix() > int64(airdrop.Expiry) {
		return nil, types.ErrAirdropExpired
	}

	// Add commitments
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	commitments := k.commitmentKeeper.GetCommitments(ctx, sender)
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, math.NewInt(int64(airdrop.Amount))))
	k.commitmentKeeper.SetCommitments(ctx, commitments)

	k.RemoveAirdrop(ctx, msg.Sender)
	return &types.MsgClaimAirdropResponse{}, nil
}
