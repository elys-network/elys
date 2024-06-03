package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/membershiptier/types"
)

func (k Keeper) ProcessPortfolioChange(ctx sdk.Context, assetType string, user string, denom string) {
	// TODO: set today's date minimum in USD and current value
	// Get USD from oracle
	sender := sdk.MustAccAddressFromBech32(user)
	switch assetType {
	case types.LiquidKeyPrefix:
		{
			balances := k.bankKeeper.GetAllBalances(ctx, sender)
			for _, balance := range balances {
				k.SetPortfolio(ctx, types.Portfolio{
					Creator:      user,
					Assetkey:     types.LiquidKeyPrefix,
					Token:        balance,
					MinimumToday: min(minimum_today, balance.Amount*usdDenomPrice),
				}, types.LiquidKeyPrefix)
			}
		}

	case types.PerpetualKeyPrefix:
	case types.PoolKeyPrefix:
	case types.StakedKeyPrefix:
	default:
	}
}

// SetPortfolio set a specific portfolio in the store from its index
func (k Keeper) SetPortfolio(ctx sdk.Context, portfolio types.Portfolio, assetType string) {
	assetKey := k.GetDateFromBlock(ctx.BlockTime()) + assetType + portfolio.Creator
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(assetKey))
	b := k.cdc.MustMarshal(&portfolio)
	store.Set(types.PortfolioKey(
		portfolio.Token.Denom,
	), b)
}

// GetPortfolio returns a portfolio from its index
func (k Keeper) GetPortfolio(
	ctx sdk.Context,
	user string,
	assetType string,
	timestamp string,
) (list []types.Portfolio) {
	ctx.BlockTime().Date()
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
) types.Portfolio {
	ctx.BlockTime().Date()
	assetKey := timestamp + assetType + user
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(assetKey))

	portfolio := store.Get(types.PortfolioKey(
		denom,
	))
	var val types.Portfolio
	k.cdc.MustUnmarshal(portfolio, &val)
	return val
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
func (k Keeper) GetAllPortfolio(ctx sdk.Context) (list []types.Portfolio) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortfolioKeyPrefix))
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
