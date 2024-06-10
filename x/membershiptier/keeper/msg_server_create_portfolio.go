package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/membershiptier/types"
)

func (k msgServer) CreatePortfolio(goCtx context.Context, msg *types.MsgCreatePortfolio) (*types.MsgCreatePortfolioResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.RetreiveAllPortfolio(ctx, msg.User)

	return &types.MsgCreatePortfolioResponse{}, nil
}
