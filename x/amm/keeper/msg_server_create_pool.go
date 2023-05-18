package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CreatePool attempts to create a pool returning the newly created pool ID or an error upon failure.
// The pool creation fee is used to fund the community pool.
// It will create a dedicated module account for the pool and sends the initial liquidity to the created module account.
func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// poolId, err = k.Keeper.CreatePool(ctx, msg)
	// if err != nil {
	// 	return 0, err
	// }

	// ctx.EventManager().EmitEvents(sdk.Events{
	// 	sdk.NewEvent(
	// 		types.TypeEvtPoolCreated,
	// 		sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
	// 	),
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 		sdk.NewAttribute(sdk.AttributeKeySender, msg.PoolCreator().String()),
	// 	),
	// })

	return &types.MsgCreatePoolResponse{
		PoolID: /*poolId*/ 0,
	}, nil
}
