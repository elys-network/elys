package keeper

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/membershiptier/types"
)

func (k Keeper) ProcessPortfolioChange(ctx sdk.Context, assetType string, user string, denom string, amount sdk.Int) {
	sender := sdk.MustAccAddressFromBech32(user)
	switch assetType {
	case types.LiquidKeyPrefix:
		{
			balances := k.bankKeeper.GetAllBalances(ctx, sender)
			for _, balance := range balances {
				tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, balance.Denom)
				prevMin, found := k.GetPortfolioMinimumToday(ctx, user, assetType, k.GetDateFromBlock(ctx.BlockTime()), balance.Denom)
				totalValue := balance.Amount.ToLegacyDec().Mul(tokenPrice)
				if prevMin.LT(totalValue) && found {
					totalValue = prevMin
				}
				k.SetPortfolio(ctx, types.Portfolio{
					Creator:      user,
					Assetkey:     types.LiquidKeyPrefix,
					Denom:        denom,
					Amount:       amount.Uint64(),
					MinimumToday: totalValue,
				}, types.LiquidKeyPrefix)
			}
		}
	case types.PerpetualKeyPrefix:
		{
			prevMin, found := k.GetPortfolioMinimumToday(ctx, user, assetType, k.GetDateFromBlock(ctx.BlockTime()), denom)
			totalValue := amount.ToLegacyDec()
			if prevMin.LT(totalValue) && found {
				totalValue = prevMin
			}
			k.SetPortfolio(ctx, types.Portfolio{
				Creator:      user,
				Assetkey:     types.PerpetualKeyPrefix,
				Denom:        denom,
				Amount:       amount.Uint64(),
				MinimumToday: totalValue,
			}, types.PerpetualKeyPrefix)
		}
	case types.PoolKeyPrefix:
		{
			// TODO: Check commitment logic to enable pool value tracking
			poolId, err := GetPoolIdFromShareDenom(denom)
			if err != nil {
				return
			}
			pool, found := k.amm.GetPool(ctx, poolId)
			if !found {
				return
			}
			info := k.amm.PoolExtraInfo(ctx, pool)

			prevMin, found := k.GetPortfolioMinimumToday(ctx, user, assetType, k.GetDateFromBlock(ctx.BlockTime()), denom)
			totalValue := amount.ToLegacyDec().Mul(info.LpTokenPrice)
			if prevMin.LT(totalValue) && found {
				totalValue = prevMin
			}
			k.SetPortfolio(ctx, types.Portfolio{
				Creator:      user,
				Assetkey:     types.PoolKeyPrefix,
				Denom:        denom,
				Amount:       amount.Uint64(),
				MinimumToday: totalValue,
			}, types.PoolKeyPrefix)
		}
	case types.StakedKeyPrefix:
		{
			// Get USDC earn program
			// Get elys earn program
			// Get eden earn program
			// Get eden boost program
			// Check all events and add hook for these
		}
	default:
	}
}

// SetPortfolio set a specific portfolio in the store from its index
func (k Keeper) SetPortfolio(ctx sdk.Context, portfolio types.Portfolio, assetType string) {
	assetKey := k.GetDateFromBlock(ctx.BlockTime()) + assetType + portfolio.Creator
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(assetKey))
	b := k.cdc.MustMarshal(&portfolio)
	store.Set(types.PortfolioKey(
		portfolio.Denom,
	), b)
}

func (k Keeper) GetMembershipTier(ctx sdk.Context, user string) (total_portfoilio sdk.Dec, tier string, discount uint64) {
	year, month, day := ctx.BlockTime().Date()
	dateToday := time.Date(year, month, day, 0, 0, 0, 0, ctx.BlockTime().Location())
	startDate := dateToday.AddDate(0, 0, -7)
	minTotal := sdk.NewDec(math.MaxInt64)
	for d := startDate; !d.After(dateToday); d = d.AddDate(0, 0, 1) {
		// Traverse all possible portfolio data
		portLiq := k.GetPortfolioTotal(ctx, user, types.LiquidKeyPrefix, d.Format("2006-01-02"))
		portPerp := k.GetPortfolioTotal(ctx, user, types.PerpetualKeyPrefix, d.Format("2006-01-02"))
		portPool := k.GetPortfolioTotal(ctx, user, types.PoolKeyPrefix, d.Format("2006-01-02"))
		portStaked := k.GetPortfolioTotal(ctx, user, types.StakedKeyPrefix, d.Format("2006-01-02"))
		totalPort := portLiq.Add(portPool).Add(portPerp).Add(portStaked)
		// TODO: add rewards
		if totalPort.LT(minTotal) {
			minTotal = totalPort
		}
	}

	if minTotal.GTE(sdk.NewDec(500000)) {
		return minTotal, "platinum", 30
	}

	if minTotal.GTE(sdk.NewDec(250000)) {
		return minTotal, "gold", 20
	}

	if minTotal.GTE(sdk.NewDec(50000)) {
		return minTotal, "silver", 10
	}

	return minTotal, "bronze", 0
}

// GetPortfolio returns a portfolio from its index
func (k Keeper) GetPortfolioTotal(
	ctx sdk.Context,
	user string,
	assetType string,
	timestamp string,
) (total sdk.Dec) {
	assetKey := timestamp + assetType + user
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(assetKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	total = sdk.NewDec(0)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Portfolio
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		total = total.Add(val.MinimumToday)
	}

	return
}

// GetPortfolio returns a portfolio from its index
func (k Keeper) GetPortfolio(
	ctx sdk.Context,
	user string,
	assetType string,
	timestamp string,
) (list []types.Portfolio) {
	assetKey := timestamp + assetType + user
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(assetKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Portfolio
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPortfolio returns a portfolio from its index
func (k Keeper) GetPortfolioMinimumToday(
	ctx sdk.Context,
	user string,
	assetType string,
	timestamp string,
	denom string,
) (sdk.Dec, bool) {
	assetKey := timestamp + assetType + user
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(assetKey))

	found := store.Has(types.PortfolioKey(
		denom,
	))

	if !found {
		return sdk.NewDec(0), false
	}

	portfolio := store.Get(types.PortfolioKey(
		denom,
	))
	var val types.Portfolio
	k.cdc.MustUnmarshal(portfolio, &val)
	return val.MinimumToday, true
}

// RemovePortfolio removes a portfolio from the store
func (k Keeper) RemovePortfolio(
	ctx sdk.Context,
	user string,
	assetType string,
	timestamp string,
) {
	assetKey := timestamp + assetType + user
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(assetKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(types.PortfolioKey(
			string(iterator.Key()),
		))
	}
}

// GetAllPortfolio returns all portfolio
func (k Keeper) GetAllPortfolio(ctx sdk.Context, timestamp string) (list []types.Portfolio) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(timestamp))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Portfolio
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetDateFromBlock(blockTime time.Time) string {
	// Extract the year, month, and day
	year, month, day := blockTime.Date()
	// Create a new time.Time object with the extracted date and time set to midnight
	blockDate := time.Date(year, month, day, 0, 0, 0, 0, blockTime.Location())
	// Format the date as a string in the "%Y-%m-%d" format
	return blockDate.Format("2006-01-02")
}

func GetPoolIdFromShareDenom(shareDenom string) (uint64, error) {
	poolId, err := strconv.Atoi(strings.TrimPrefix(shareDenom, "amm/pool/"))
	if err != nil {
		return 0, err
	}
	return uint64(poolId), nil
}
