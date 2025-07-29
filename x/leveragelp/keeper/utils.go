package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/utils"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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

func (k Keeper) CheckMaxLeverageRatio(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrapf(types.ErrPoolDoesNotExist, "leverage lp pool: %d", poolId)
	}

	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrapf(types.ErrPoolDoesNotExist, "amm pool: %d", poolId)
	}

	poolLeverageRatio := pool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	if poolLeverageRatio.GTE(pool.MaxLeveragelpRatio) {
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

		ammPool, err := k.GetAmmPool(ctx, position.AmmPoolId)
		if err != nil {
			return nil, err
		}

		ammTVL, err := ammPool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}

		debtDenomPrice := k.oracleKeeper.GetDenomPrice(ctx, position.Collateral.Denom)
		debtValue := position.GetBigDecLiabilities().Mul(debtDenomPrice)

		positionValue := position.GetBigDecLeveragedLpAmount().Mul(ammTVL).Quo(osmomath.BigDecFromSDKInt(ammPool.TotalShares.Amount))

		updated_leverage := osmomath.ZeroBigDec()
		denominator := positionValue.Sub(debtValue)
		if denominator.IsPositive() {
			updated_leverage = positionValue.Quo(denominator)
		}
		if debtValue.IsPositive() {
			position.PositionHealth = positionValue.Quo(debtValue).Dec()
		}

		updatedLeveragePositions = append(updatedLeveragePositions, &types.QueryPosition{
			Position:         position,
			UpdatedLeverage:  updated_leverage.Dec(),
			PositionUsdValue: positionValue.Dec(),
		})
	}
	return updatedLeveragePositions, nil
}

func (k Keeper) GetInterestRateUsd(ctx sdk.Context, positions []*types.QueryPosition) ([]*types.PositionAndInterest, error) {
	positions_and_interest := []*types.PositionAndInterest{}

	for _, position := range positions {
		pool, found := k.stableKeeper.GetPoolByDenom(ctx, position.Position.Collateral.Denom)
		if !found {
			return nil, errorsmod.Wrap(types.ErrPoolNotCreatedForBorrow, fmt.Sprintf("Asset: %s", position.Position.Collateral.Denom))
		}

		var positionAndInterest types.PositionAndInterest
		positionAndInterest.Position = position
		price := k.oracleKeeper.GetDenomPrice(ctx, position.Position.Collateral.Denom)
		interestRateHour := pool.GetBigDecInterestRate().Quo(utils.HoursInYear)
		positionAndInterest.InterestRateHour = interestRateHour.Dec()
		positionAndInterest.InterestRateHourUsd = interestRateHour.Mul(position.Position.GetBigDecLiabilities()).Mul(price).Dec()
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
				position.PositionHealth = positionHealth.Dec()
				k.SetPosition(ctx, &position)
			}
		}
	}
}
