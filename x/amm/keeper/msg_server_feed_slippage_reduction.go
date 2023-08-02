package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) FeedSlippageReduction(goCtx context.Context, msg *types.MsgFeedSlippageReduction) (*types.MsgFeedSlippageReductionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeder, found := k.oracleKeeper.GetPriceFeeder(ctx, msg.Sender)
	if !found {
		return nil, oracletypes.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, oracletypes.ErrPriceFeederNotActive
	}

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrInvalidPoolId
	}

	pool.PoolParams.SlippageReduction = msg.SlippageReduction
	k.SetPool(ctx, pool)

	return &types.MsgFeedSlippageReductionResponse{}, nil
}
