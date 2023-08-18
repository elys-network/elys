package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.CheckMaxOpenPositions(ctx); err != nil {
		return nil, err
	}

	// Get token asset other than USDC
	noneNativeAsset := k.GetNoneNativeAsset(msg.CollateralAsset, msg.BorrowAsset)

	// Get the first valid pool
	poolId, err := k.GetFirstValidPool(ctx, noneNativeAsset)
	if err != nil {
		return nil, err
	}

	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId, noneNativeAsset)
	if err != nil {
		return nil, err
	}

	pool, found := k.PoolChecker.GetPool(ctx, poolId)
	// If margin pool doesn't exist yet, we should initiate it according to its corresponding ammPool
	if !found {
		pool = types.NewPool(poolId)
		pool.InitiatePool(ctx, ammPool)

		k.OpenLongChecker.SetPool(ctx, pool)
	}

	if err := k.CheckPoolHealth(ctx, poolId); err != nil {
		return nil, err
	}

	var mtp *types.MTP
	switch msg.Position {
	case types.Position_LONG:
		mtp, err = k.OpenLong(ctx, poolId, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	ctx.EventManager().EmitEvent(k.GenerateOpenEvent(mtp))
	return &types.MsgOpenResponse{}, nil
}
