package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) SetPriceFeeder(goCtx context.Context, msg *types.MsgSetPriceFeeder) (*types.MsgSetPriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if msg.Creator != params.ModuleAdmin {
		return nil, types.ErrNotModuleAdmin
	}
	k.Keeper.SetPriceFeeder(ctx, types.PriceFeeder{
		Feeder:   msg.Feeder,
		IsActive: msg.IsActive,
	})
	return &types.MsgSetPriceFeederResponse{}, nil
}

func (k msgServer) DeletePriceFeeder(goCtx context.Context, msg *types.MsgDeletePriceFeeder) (*types.MsgDeletePriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if msg.Creator != params.ModuleAdmin {
		return nil, types.ErrNotModuleAdmin
	}
	k.RemovePriceFeeder(ctx, msg.Feeder)
	return &types.MsgDeletePriceFeederResponse{}, nil
}
