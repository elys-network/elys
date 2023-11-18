package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, mtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MsgOpenResponse, error) {
	// Get token asset other than base currency
	tradingAsset := types.GetTradingAsset(msg.CollateralAsset, msg.BorrowAsset, baseCurrency)

	poolId := mtp.AmmPoolId
	pool, found := k.OpenLongChecker.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, tradingAsset)
	}

	if !k.OpenLongChecker.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, tradingAsset)
	}

	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId, tradingAsset)
	if err != nil {
		return nil, err
	}

	switch msg.Position {
	case types.Position_LONG:
		mtp, err = k.OpenConsolidateLong(ctx, poolId, mtp, msg, baseCurrency)
		if err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		mtp, err = k.OpenConsolidateShort(ctx, poolId, mtp, msg, baseCurrency)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	ctx.EventManager().EmitEvent(k.GenerateOpenEvent(mtp))

	if k.hooks != nil {
		k.hooks.AfterMarginPositionModified(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{}, nil
}
