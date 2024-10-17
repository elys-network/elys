package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MsgOpenResponse, error) {
	poolId := existingMtp.AmmPoolId
	if !k.OpenDefineAssetsChecker.IsPoolEnabled(ctx, poolId) {
		return nil, errorsmod.Wrap(types.ErrMTPDisabled, existingMtp.CustodyAsset)
	}

	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, existingMtp)

	existingMtp, err = k.OpenConsolidateMergeMtp(ctx, poolId, existingMtp, newMtp, baseCurrency)
	if err != nil {
		return nil, err
	}

	consolidatedOpenPrice := (existingMtp.Custody.ToLegacyDec().Mul(existingMtp.OpenPrice).Add(newMtp.Custody.ToLegacyDec().Mul(newMtp.OpenPrice))).Quo(existingMtp.Custody.ToLegacyDec().Add(newMtp.Custody.ToLegacyDec()))
	existingMtp.OpenPrice = consolidatedOpenPrice

	consolidatedTakeProfitPrice := existingMtp.Custody.ToLegacyDec().Mul(existingMtp.TakeProfitPrice).Add(newMtp.Custody.ToLegacyDec().Mul(newMtp.TakeProfitPrice)).Quo(existingMtp.Custody.ToLegacyDec().Add(newMtp.Custody.ToLegacyDec()))
	existingMtp.TakeProfitPrice = consolidatedTakeProfitPrice

	existingMtp.TakeProfitCustody = existingMtp.TakeProfitCustody.Add(newMtp.TakeProfitCustody)
	existingMtp.TakeProfitLiabilities = existingMtp.TakeProfitLiabilities.Add(newMtp.TakeProfitLiabilities)

	// Set existing MTP
	if err = k.SetMTP(ctx, existingMtp); err != nil {
		return nil, err
	}

	k.EmitOpenEvent(ctx, existingMtp)

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.hooks != nil {
		err = k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool, creator)
		if err != nil {
			return nil, err
		}
	}

	if err = k.CheckPoolHealth(ctx, poolId); err != nil {
		return nil, err
	}

	return &types.MsgOpenResponse{
		Id: existingMtp.Id,
	}, nil
}

func (k Keeper) EmitOpenEvent(ctx sdk.Context, mtp *types.MTP) {
	ctx.EventManager().EmitEvent(types.GenerateOpenEvent(mtp))
}
