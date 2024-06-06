package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/membershiptier/types"
)

func (k msgServer) SetPortfolio(goCtx context.Context, msg *types.MsgSetPortfolio) (*types.MsgSetPortfolioResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.RetreiveAllPortfolio(ctx, msg.User)

	return &types.MsgSetPortfolioResponse{}, nil
}
