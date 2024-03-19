package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/launchpad/types"
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

func (k msgServer) BuyElys(goCtx context.Context, msg *types.MsgBuyElys) (*types.MsgBuyElysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx
	// TODO: implement

	return &types.MsgBuyElysResponse{}, nil
}

func (k msgServer) ReturnElys(goCtx context.Context, msg *types.MsgReturnElys) (*types.MsgReturnElysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx
	// TODO: implement

	return &types.MsgReturnElysResponse{}, nil
}

func (k msgServer) WithdrawRaised(goCtx context.Context, msg *types.MsgWithdrawRaised) (*types.MsgWithdrawRaisedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if params.WithdrawAddress != msg.Sender {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "expected %s, got %s", params.WithdrawAddress, msg.Sender)
	}

	// TODO: implement

	return &types.MsgWithdrawRaisedResponse{}, nil
}
