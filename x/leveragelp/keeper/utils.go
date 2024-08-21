package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) CheckUserAuthorization(ctx sdk.Context, msg *types.MsgOpen) error {
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.IsWhitelistingEnabled(ctx) && !k.CheckIfWhitelisted(ctx, creator) {
		return errorsmod.Wrap(types.ErrUnauthorised, "unauthorised")
	}
	return nil
}

func (k Keeper) CheckSamePosition(ctx sdk.Context, msg *types.MsgOpen) *types.Position {
	positions, _, err := k.GetPositionsForAddress(ctx, sdk.MustAccAddressFromBech32(msg.Creator), &query.PageRequest{})
	if err != nil {
		return nil
	}
	for _, position := range positions {
		if position.Position.AmmPoolId == msg.AmmPoolId && position.Position.Collateral.Denom == msg.CollateralAsset {
			return position.Position
		}
	}

	return nil
}

func (k Keeper) CheckPoolHealth(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	if !k.IsPoolEnabled(ctx, poolId) || k.IsPoolClosed(ctx, poolId) {
		return errorsmod.Wrap(types.ErrPositionDisabled, "pool is disabled or closed")
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.GetPoolOpenThreshold(ctx)) {
		return errorsmod.Wrap(types.ErrInvalidPosition, "pool health too low to open new positions")
	}
	return nil
}

func (k Keeper) CheckMaxOpenPositions(ctx sdk.Context) error {
	if k.GetOpenPositionCount(ctx) >= k.GetMaxOpenPositions(ctx) {
		return errorsmod.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}
	return nil
}

func (k Keeper) GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error) {
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return ammPool, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}
	return ammPool, nil
}

func (k Keeper) GetLeverageLpUpdatedLeverage(ctx sdk.Context, positions []*types.Position) ([]*types.QueryPosition, error) {
	updatedLeveragePositions := []*types.QueryPosition{}
	for i, position := range positions {
		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if !found {
			return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
		}
		exitCoinsAfterFee, _, err := k.amm.ExitPoolEst(ctx, position.GetAmmPoolId(), position.LeveragedLpAmount, baseCurrency)
		if err != nil {
			return nil, err
		}

		exitAmountAfterFee := exitCoinsAfterFee.AmountOf(baseCurrency)

		price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, position.Collateral.Denom)
		updated_leverage := exitAmountAfterFee.Quo(exitAmountAfterFee.Sub(position.Liabilities.Mul(price.TruncateInt())))

		updatedLeveragePositions[i] = &types.QueryPosition{
			Position:        position,
			UpdatedLeverage: updated_leverage,
		}
	}
	return updatedLeveragePositions, nil
}