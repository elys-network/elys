package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
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

	if err := k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}

	// check if existing mtp to consolidate
	existingMtp := k.CheckSameAssetPosition(ctx, msg)

	if existingMtp == nil && msg.Leverage.Equal(math.LegacyOneDec()) {
		return nil, fmt.Errorf("cannot open new position with leverage 1")
	}

	if err := k.CheckMaxOpenPositions(ctx); err != nil {
		return nil, err
	}

	poolId := msg.PoolId
	// Get pool id, amm pool, and perpetual pool
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "amm pool not found for pool %d", poolId)
	}

	if !ammPool.PoolParams.UseOracle {
		return nil, types.ErrPoolHasToBeOracle
	}

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, types.ErrPoolDoesNotExist
	}

	if !k.PoolChecker.IsPoolEnabled(ctx, poolId) || k.PoolChecker.IsPoolClosed(ctx, poolId) {
		return nil, errorsmod.Wrap(types.ErrMTPDisabled, "pool is disabled or closed")
	}

	if err = k.CheckLowPoolHealth(ctx, poolId); err != nil {
		return nil, err
	}

	mtp, err := k.OpenDefineAssets(ctx, poolId, msg, baseCurrency, isBroker)
	if err != nil {
		return nil, err
	}

	// calc and update open price
	err = k.UpdateOpenPrice(ctx, mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	if existingMtp != nil {
		return k.OpenConsolidate(ctx, existingMtp, mtp, msg, baseCurrency)
	}

	if err = k.CheckLowPoolHealth(ctx, poolId); err != nil {
		return nil, err
	}

	k.EmitOpenEvent(ctx, mtp)

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.hooks != nil {
		// pool values has been updated
		pool, found = k.GetPool(ctx, poolId)
		if !found {
			return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
		}

		err = k.hooks.AfterPerpetualPositionOpen(ctx, ammPool, pool, creator)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgOpenResponse{
		Id: mtp.Id,
	}, nil
}
