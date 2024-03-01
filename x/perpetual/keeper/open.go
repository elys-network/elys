package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen, isBroker bool) (*types.MsgOpenResponse, error) {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Determine the type of position (long or short) and validate assets accordingly.
	switch msg.Position {
	case types.Position_LONG:
		if err := types.CheckLongAssets(msg.Collateral.Denom, msg.TradingAsset, baseCurrency); err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		if err := types.CheckShortAssets(msg.Collateral.Denom, msg.TradingAsset, baseCurrency); err != nil {
			return nil, err
		}
	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	if err := k.OpenChecker.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}

	// check if existing mtp to consolidate
	existingMtp := k.OpenChecker.CheckSameAssetPosition(ctx, msg)

	if err := k.OpenChecker.CheckMaxOpenPositions(ctx); err != nil {
		return nil, err
	}

	// Get pool id, amm pool, and perpetual pool
	poolId, ammPool, pool, err := k.OpenChecker.PreparePools(ctx, msg.Collateral.Denom, msg.TradingAsset)
	if err != nil {
		return nil, err
	}

	if err := k.OpenChecker.CheckPoolHealth(ctx, poolId); err != nil {
		return nil, err
	}

	var mtp *types.MTP
	switch msg.Position {
	case types.Position_LONG:
		mtp, err = k.OpenChecker.OpenLong(ctx, poolId, msg, baseCurrency, isBroker)
		if err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		mtp, err = k.OpenChecker.OpenShort(ctx, poolId, msg, baseCurrency, isBroker)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	if existingMtp != nil {
		return k.OpenConsolidate(ctx, existingMtp, mtp, msg, baseCurrency)
	}

	// calc and update open price
	k.OpenChecker.UpdateOpenPrice(ctx, mtp, ammPool, baseCurrency)

	k.OpenChecker.EmitOpenEvent(ctx, mtp)

	if k.hooks != nil {
		k.hooks.AfterPerpetualPositionOpen(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{
		Id: mtp.Id,
	}, nil
}
