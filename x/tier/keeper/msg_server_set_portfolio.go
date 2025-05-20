package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/tier/types"
)

func (k msgServer) SetPortfolio(goCtx context.Context, msg *types.MsgSetPortfolio) (*types.MsgSetPortfolioResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	user := sdk.MustAccAddressFromBech32(msg.User)
	k.RetrieveAllPortfolio(ctx, user)

	return &types.MsgSetPortfolioResponse{}, nil
}
