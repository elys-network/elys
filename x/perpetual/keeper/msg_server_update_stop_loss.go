package keeper

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) UpdateStopLoss(goCtx context.Context, msg *types.MsgUpdateStopLoss) (*types.MsgUpdateStopLossResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Load existing mtp
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, creator, msg.Id)
	if err != nil {
		return nil, err
	}

	poolId := mtp.AmmPoolId
	_, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	if mtp.Position == types.Position_LONG {
		if !msg.Price.IsZero() && msg.Price.GTE(tradingAssetPrice) {
			return nil, fmt.Errorf("stop loss price cannot be greater than equal to tradingAssetPrice for long (Stop loss: %s, asset price: %s)", msg.Price.String(), tradingAssetPrice.String())
		}
	}
	if mtp.Position == types.Position_SHORT {
		if !msg.Price.IsZero() && msg.Price.LTE(tradingAssetPrice) {
			return nil, fmt.Errorf("stop loss price cannot be less than equal to tradingAssetPrice for short (Stop loss: %s, asset price: %s)", msg.Price.String(), tradingAssetPrice.String())
		}
	}

	mtp.StopLossPrice = msg.Price
	err = k.SetMTP(ctx, &mtp)
	if err != nil {
		return nil, err
	}

	event := sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("stop_loss", mtp.StopLossPrice.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgUpdateStopLossResponse{}, nil
}
