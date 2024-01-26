package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, mtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MsgOpenResponse, error) {
	poolId := mtp.AmmPoolId
	pool, found := k.OpenLongChecker.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, mtp.CustodyAsset)
	}

	if !k.OpenLongChecker.IsPoolEnabled(ctx, poolId) {
		return nil, errorsmod.Wrap(types.ErrMTPDisabled, mtp.CustodyAsset)
	}

	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId, mtp.CustodyAsset)
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
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	ctx.EventManager().EmitEvent(types.GenerateOpenEvent(mtp))

	if k.hooks != nil {
		k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{}, nil
}
