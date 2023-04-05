package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) FeedPrice(goCtx context.Context, msg *types.MsgFeedPrice) (*types.MsgFeedPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetPrice(
		ctx,
		msg.Asset,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var price = types.Price{
		Provider:  msg.Provider,
		Asset:     msg.Asset,
		Price:     msg.Price,
		Source:    msg.Source,
		Timestamp: uint64(ctx.BlockTime().Unix()),
	}

	k.SetPrice(ctx, price)
	return &types.MsgFeedPriceResponse{}, nil
}
