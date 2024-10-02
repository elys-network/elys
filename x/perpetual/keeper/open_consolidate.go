package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MsgOpenResponse, error) {
	poolId := existingMtp.AmmPoolId
	pool, found := k.OpenLongChecker.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, newMtp.CustodyAsset)
	}

	if !k.OpenLongChecker.IsPoolEnabled(ctx, poolId) {
		return nil, errorsmod.Wrap(types.ErrMTPDisabled, existingMtp.CustodyAsset)
	}

	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId, existingMtp.CustodyAsset)
	if err != nil {
		return nil, err
	}

	switch msg.Position {
	case types.Position_LONG:
		existingMtp, err = k.OpenConsolidateLong(ctx, poolId, existingMtp, newMtp, msg, baseCurrency)
		if err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		existingMtp, err = k.OpenConsolidateShort(ctx, poolId, existingMtp, newMtp, msg, baseCurrency)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	// calc and update open price
	err = k.UpdateOpenPrice(ctx, existingMtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	k.EmitOpenEvent(ctx, existingMtp)

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.hooks != nil {
		k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool, creator)
	}

	return &types.MsgOpenResponse{
		Id: existingMtp.Id,
	}, nil
}

func (k Keeper) EmitOpenEvent(ctx sdk.Context, mtp *types.MTP) {
	ctx.EventManager().EmitEvent(types.GenerateOpenEvent(mtp))
}
