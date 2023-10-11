package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, mtp *types.MTP, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	// Get token asset other than base currency
	tradingAsset := k.OpenLongChecker.GetTradingAsset(msg.CollateralAsset, msg.BorrowAsset)

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
		mtp, err = k.OpenConsolidateLong(ctx, poolId, mtp, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	ctx.EventManager().EmitEvent(k.GenerateOpenEvent(mtp))

	if k.hooks != nil {
		k.hooks.AfterLeveragelpPositionModified(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{}, nil
}
