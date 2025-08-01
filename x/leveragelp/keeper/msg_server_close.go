package keeper

import (
	"context"
	"errors"
	"strconv"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k msgServer) Close(goCtx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	position, err := k.GetPosition(ctx, creator, msg.Id)
	if err != nil {
		return nil, err
	}

	if position.LeveragedLpAmount.IsZero() {
		return nil, types.ErrAmountTooLow
	}

	// Retrieve Pool
	pool, found := k.GetPool(ctx, position.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	closingRatio := osmomath.BigDecFromSDKInt(msg.LpAmount).Quo(position.GetBigDecLeveragedLpAmount())
	if closingRatio.GT(osmomath.OneBigDec()) {
		return nil, errors.New("invalid closing ratio for leverage lp")
	}

	finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, stopLossReached, _, exitSlippageFeeOnClosingPosition, swapFee, takerFee, slippageValue, swapFeeValue, takerFeeValue, weightBreakingFeeValue, err := k.CheckHealthStopLossThenRepayAndClose(ctx, &position, &pool, closingRatio, false)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
		sdk.NewAttribute("address", msg.Creator),
		sdk.NewAttribute("closing_ratio", finalClosingRatio.String()),
		sdk.NewAttribute("lp_amount_closed", totalLpAmountToClose.String()),
		sdk.NewAttribute("coins_to_amm", coinsForAmm.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("user_return_tokens", userReturnTokens.String()),
		sdk.NewAttribute("exit_fee", exitFeeOnClosingPosition.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
		sdk.NewAttribute("stop_loss_reached", strconv.FormatBool(stopLossReached)),
		sdk.NewAttribute("exit_slippage_fee", exitSlippageFeeOnClosingPosition.String()),
		sdk.NewAttribute("exit_slippage_fee_value_in_usd", slippageValue.String()),
		sdk.NewAttribute("exit_swap_fee", swapFee.String()),
		sdk.NewAttribute("exit_swap_fee_value_in_usd", swapFeeValue.String()),
		sdk.NewAttribute("exit_taker_fee", takerFee.String()),
		sdk.NewAttribute("exit_taker_fee_value_in_usd", takerFeeValue.String()),
		sdk.NewAttribute("exit_weight_breaking_fee_value_in_usd", weightBreakingFeeValue.String()),
	))
	return &types.MsgCloseResponse{}, nil
}
