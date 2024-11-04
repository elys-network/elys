package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MsgOpenResponse, error) {
	poolId := existingMtp.AmmPoolId
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, existingMtp)

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	err = k.SettleFunding(ctx, existingMtp, &pool, ammPool)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "error handling funding fee")
	}

	existingMtp, err = k.OpenConsolidateMergeMtp(ctx, existingMtp, newMtp)
	if err != nil {
		return nil, err
	}

	if !newMtp.Liabilities.IsZero() {
		consolidatedOpenPrice := (existingMtp.Custody.ToLegacyDec().Mul(existingMtp.OpenPrice).Add(newMtp.Custody.ToLegacyDec().Mul(newMtp.OpenPrice))).Quo(existingMtp.Custody.ToLegacyDec().Add(newMtp.Custody.ToLegacyDec()))
		existingMtp.OpenPrice = consolidatedOpenPrice

		consolidatedTakeProfitPrice := existingMtp.Custody.ToLegacyDec().Mul(existingMtp.TakeProfitPrice).Add(newMtp.Custody.ToLegacyDec().Mul(newMtp.TakeProfitPrice)).Quo(existingMtp.Custody.ToLegacyDec().Add(newMtp.Custody.ToLegacyDec()))
		existingMtp.TakeProfitPrice = consolidatedTakeProfitPrice
	}

	existingMtp.TakeProfitCustody = existingMtp.TakeProfitCustody.Add(newMtp.TakeProfitCustody)
	existingMtp.TakeProfitLiabilities = existingMtp.TakeProfitLiabilities.Add(newMtp.TakeProfitLiabilities)

	// no need to update pool's TakeProfitCustody, TakeProfitLiabilities, Custody and Liabilities as it was already in OpenDefineAssets

	existingMtp.MtpHealth, err = k.GetMTPHealth(ctx, *existingMtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.GetSafetyFactor(ctx)
	if existingMtp.MtpHealth.LTE(safetyFactor) {
		return nil, errorsmod.Wrapf(types.ErrMTPUnhealthy, "(MtpHealth: %s)", existingMtp.MtpHealth.String())
	}

	// Set existing MTP
	if err = k.SetMTP(ctx, existingMtp); err != nil {
		return nil, err
	}

	k.EmitOpenEvent(ctx, existingMtp)

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.hooks != nil {
		params := k.GetParams(ctx)
		// The pool value above was sent in pointer so its updated
		err = k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool, creator, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return nil, err
		}
	}

	if err = k.CheckLowPoolHealthAndMinimumCustody(ctx, poolId); err != nil {
		return nil, err
	}

	return &types.MsgOpenResponse{
		Id: existingMtp.Id,
	}, nil
}

func (k Keeper) EmitOpenEvent(ctx sdk.Context, mtp *types.MTP) {
	ctx.EventManager().EmitEvent(types.GenerateOpenEvent(mtp))
}
