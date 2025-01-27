package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) Close(goCtx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	position, err := k.GetPosition(ctx, creator, msg.Id)
	if err != nil {
		return nil, err
	}

	// Retrieve Pool
	pool, found := k.GetPool(ctx, position.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	positionHealth, err := k.GetPositionHealth(ctx, position)
	if err != nil {
		return nil, err
	}
	safetyFactor := k.GetSafetyFactor(ctx)

	// If lpAmount is lower than zero or position is unhealthy, close full amount
	lpAmount := msg.LpAmount
	if lpAmount.IsNil() || lpAmount.LTE(math.ZeroInt()) || positionHealth.LTE(safetyFactor) {
		lpAmount = position.LeveragedLpAmount
	}

	closingRatio := lpAmount.ToLegacyDec().Quo(position.LeveragedLpAmount.ToLegacyDec())

	finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, finalUserRewards, err := k.CheckHealthStopLossThenRepayAndClose(ctx, &position, &pool, closingRatio, false)
	if err != nil {
		return nil, err
	}

	if k.hooks != nil {
		ammPool, err := k.GetAmmPool(ctx, position.AmmPoolId)
		if err != nil {
			return nil, err
		}
		err = k.hooks.AfterLeverageLpPositionClose(ctx, sdk.MustAccAddressFromBech32(msg.Creator), ammPool)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("closing_ratio", finalClosingRatio.String()),
		sdk.NewAttribute("lp_amount_closed", totalLpAmountToClose.String()),
		sdk.NewAttribute("coins_to_amm", coinsForAmm.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("user_rewards", finalUserRewards.String()),
	))
	return &types.MsgCloseResponse{}, nil
}
