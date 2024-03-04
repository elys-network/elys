package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) FeedMultiplePrices(goCtx context.Context, msg *types.MsgFeedMultiplePrices) (*types.MsgFeedMultiplePricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeder, found := k.Keeper.GetPriceFeeder(ctx, msg.Creator)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, types.ErrPriceFeederNotActive
	}

	for _, price := range msg.Prices {
		price.Provider = msg.Creator
		price.Timestamp = uint64(ctx.BlockTime().Unix())
		price.BlockHeight = uint64(ctx.BlockHeight())
		k.SetPrice(ctx, price)
	}

	return &types.MsgFeedMultiplePricesResponse{}, nil
}
