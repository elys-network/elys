package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MsgOpenResponse, error) {
	poolId := existingMtp.AmmPoolId
	pool, found := k.OpenDefineAssetsChecker.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, newMtp.CustodyAsset)
	}

	if !k.OpenDefineAssetsChecker.IsPoolEnabled(ctx, poolId) {
		return nil, errorsmod.Wrap(types.ErrMTPDisabled, existingMtp.CustodyAsset)
	}

	ammPool, err := k.OpenDefineAssetsChecker.GetAmmPool(ctx, poolId, existingMtp.CustodyAsset)
	if err != nil {
		return nil, err
	}

	existingMtp, err = k.OpenConsolidateMergeMtp(ctx, poolId, existingMtp, newMtp, msg, baseCurrency)
	if err != nil {
		return nil, err
	}

	// calc and update open price
	err = k.OpenDefineAssetsChecker.UpdateOpenPrice(ctx, existingMtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	k.OpenDefineAssetsChecker.EmitOpenEvent(ctx, existingMtp)

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
