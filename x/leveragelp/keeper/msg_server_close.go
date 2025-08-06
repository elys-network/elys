package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
)

func (k msgServer) Close(goCtx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	position, err := k.GetPosition(ctx, msg.PoolId, creator, msg.Id)
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

	closingRatio := msg.LpAmount.ToLegacyDec().Quo(position.LeveragedLpAmount.ToLegacyDec())
	if closingRatio.GT(math.LegacyOneDec()) {
		return nil, errors.New("invalid closing ratio for leverage lp")
	}

	finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, stopLossReached, _, exitSlippageFee, swapFee, takerFee, slippageValue, swapFeeValue, takerFeeValue, weightBreakingFeeValue, err := k.CheckHealthStopLossThenRepayAndClose(ctx, &position, &pool, closingRatio, false)
	if err != nil {
		return nil, err
	}

	k.EmitCloseEvent(ctx, "user_tx", position, finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, stopLossReached, exitSlippageFee, swapFee, takerFee, slippageValue, swapFeeValue, takerFeeValue, weightBreakingFeeValue)
	return &types.MsgCloseResponse{}, nil
}
