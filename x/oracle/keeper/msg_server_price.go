package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) FeedPrice(goCtx context.Context, msg *types.MsgFeedPrice) (*types.MsgFeedPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	provider := sdk.MustAccAddressFromBech32(msg.Provider)
	feeder, found := k.Keeper.GetPriceFeeder(ctx, provider)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, types.ErrPriceFeederNotActive
	}

	price := types.Price{
		Asset:       msg.FeedPrice.Asset,
		Price:       msg.FeedPrice.Price,
		Source:      msg.FeedPrice.Source,
		Provider:    msg.Provider,
		Timestamp:   uint64(ctx.BlockTime().Unix()),
		BlockHeight: uint64(ctx.BlockHeight()),
	}

	k.SetPrice(ctx, price)
	return &types.MsgFeedPriceResponse{}, nil
}
