package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) FeedPrice(goCtx context.Context, msg *types.MsgFeedPrice) (*types.MsgFeedPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeder, found := k.Keeper.GetPriceFeeder(ctx, msg.Provider)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, types.ErrPriceFeederNotActive
	}

	price := types.Price{
		Provider:    msg.Provider,
		Asset:       msg.Asset,
		Price:       msg.Price,
		Source:      msg.Source,
		Timestamp:   uint64(ctx.BlockTime().Unix()),
		BlockHeight: uint64(ctx.BlockHeight()),
	}

	k.SetPrice(ctx, price)
	return &types.MsgFeedPriceResponse{}, nil
}
