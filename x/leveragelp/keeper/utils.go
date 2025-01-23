package keeper

import (
	"errors"
	"fmt"

	storetypes "cosmossdk.io/core/store"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
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

func (k Keeper) CheckSamePosition(ctx sdk.Context, msg *types.MsgOpen) (*types.Position, error) {
	positions, _, err := k.GetPositionsForAddress(ctx, sdk.MustAccAddressFromBech32(msg.Creator), &query.PageRequest{})
	if err != nil {
		return nil, err
	}
	for _, position := range positions {
		if position.AmmPoolId == msg.AmmPoolId && position.Collateral.Denom == msg.CollateralAsset {
			return position, nil
		}
	}

	return nil, nil
}

// TODO simplify this design. Double check happening, one at pool level, one at global level
func (k Keeper) CheckPoolHealth(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrapf(types.ErrPoolDoesNotExist, "leverage lp pool: %d", poolId)
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.GetPoolOpenThreshold(ctx)) {
		return errors.New("pool health too low to open new positions")
	}
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrapf(types.ErrPoolDoesNotExist, "amm pool: %d", poolId)
	}

	poolLeveragelpRatio := pool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	if poolLeveragelpRatio.GT(pool.MaxLeveragelpRatio) {
		return errorsmod.Wrap(types.ErrMaxLeverageLpExists, "pool is unhealthy")
	}
	return nil
}

func (k Keeper) CheckMaxOpenPositions(ctx sdk.Context) error {

	openPositions := k.GetOpenPositionCount(ctx)
	maxOpenPositions := k.GetMaxOpenPositions(ctx)

	if openPositions >= maxOpenPositions {
		return errorsmod.Wrap(types.ErrMaxOpenPositions, fmt.Sprintf("cannot open new positions, open positions %d - max positions %d", openPositions, maxOpenPositions))
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
	for _, position := range positions {
		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if !found {
			return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
		}
		exitCoins, _, err := k.amm.ExitPoolEst(ctx, position.GetAmmPoolId(), position.LeveragedLpAmount, baseCurrency)
		if err != nil {
			return nil, err
		}

		exitAmountAfterFee := exitCoins.AmountOf(baseCurrency)

		updated_leverage := sdkmath.LegacyZeroDec()
		denomimator := exitAmountAfterFee.ToLegacyDec().Sub(position.Liabilities.ToLegacyDec())
		if denomimator.IsPositive() {
			updated_leverage = exitAmountAfterFee.ToLegacyDec().Quo(denomimator)
		}
		if position.Liabilities.IsPositive() {
			position.PositionHealth = exitAmountAfterFee.ToLegacyDec().Quo(position.Liabilities.ToLegacyDec())
		}

		updatedLeveragePositions = append(updatedLeveragePositions, &types.QueryPosition{
			Position:         position,
			UpdatedLeverage:  updated_leverage,
			PositionUsdValue: sdkmath.LegacyNewDecFromIntWithPrec(exitAmountAfterFee, 6),
		})
	}
	return updatedLeveragePositions, nil
}

func (k Keeper) GetInterestRateUsd(ctx sdk.Context, positions []*types.QueryPosition) ([]*types.PositionAndInterest, error) {
	positions_and_interest := []*types.PositionAndInterest{}
	params := k.stableKeeper.GetParams(ctx)
	hours := sdkmath.LegacyNewDec(365 * 24)

	for _, position := range positions {
		var positionAndInterest types.PositionAndInterest
		positionAndInterest.Position = position
		price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, position.Position.Collateral.Denom)
		interestRateHour := params.InterestRate.Quo(hours)
		positionAndInterest.InterestRateHour = interestRateHour
		positionAndInterest.InterestRateHourUsd = interestRateHour.Mul(sdkmath.LegacyDec(position.Position.Liabilities.Mul(price.RoundInt())))
		positions_and_interest = append(positions_and_interest, &positionAndInterest)
	}

	return positions_and_interest, nil
}

// migrating eixsting position and setting position health to max dec when liablities is zero
func (k Keeper) MigratePositionHealth(ctx sdk.Context) {
	iterator := k.GetPositionIterator(ctx)
	defer func(iterator storetypes.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		bytesValue := iterator.Value()
		err := k.cdc.Unmarshal(bytesValue, &position)
		if err == nil {
			positionHealth, err := k.GetPositionHealth(ctx, position)
			if err == nil {
				position.PositionHealth = positionHealth
				k.SetPosition(ctx, &position)
			}
		}
	}
}
