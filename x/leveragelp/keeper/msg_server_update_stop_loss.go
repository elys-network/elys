package keeper

import (
	"context"
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

func (k msgServer) UpdateStopLoss(goCtx context.Context, msg *types.MsgUpdateStopLoss) (*types.MsgUpdateStopLossResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if !params.StopLossEnabled {
		return nil, errors.New("stop loss price not enabled")
	}

	position, found := k.GetPositionWithId(ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Position)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPositionDoesNotExist, fmt.Sprintf("positionId: %d", msg.Position))

	}

	poolId := position.AmmPoolId
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	position.StopLossPrice = msg.Price
	k.SetPosition(ctx, position)

	// Add trigger function
	_, closeAttempted, _, err := k.CheckAndLiquidateUnhealthyPosition(ctx, position, pool)
	if closeAttempted && err != nil {
		return nil, err
	}

	event := sdk.NewEvent(types.EventUpdateStopLoss,
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
		sdk.NewAttribute("updated_stop_loss", position.StopLossPrice.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgUpdateStopLossResponse{}, nil
}
