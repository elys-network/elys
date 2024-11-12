package keeper

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) UpdateStopLoss(goCtx context.Context, msg *types.MsgUpdateStopLoss) (*types.MsgUpdateStopLossResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	position, found := k.GetPositionWithId(ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Position)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPositionDoesNotExist, fmt.Sprintf("positionId: %d", msg.Position))

	}

	poolId := position.AmmPoolId
	_, found = k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	position.StopLossPrice = msg.Price
	k.SetPosition(ctx, position)

	event := sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
		sdk.NewAttribute("stop_loss", position.StopLossPrice.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgUpdateStopLossResponse{}, nil
}
