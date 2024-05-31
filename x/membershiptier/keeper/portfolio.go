package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/membershiptier/types"
)

// SetPortfolio set a specific portfolio in the store from its index
func (k Keeper) SetPortfolio(ctx sdk.Context, portfolio types.Portfolio, assetType string) {
	assetKey := assetType + portfolio.Creator
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
) (list []types.Portfolio) {
	assetKey := assetType + user
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

// RemovePortfolio removes a portfolio from the store
func (k Keeper) RemovePortfolio(
	ctx sdk.Context,
	user string,
	assetType string,
	denom string,
) {
	assetKey := assetType + user
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(assetKey))
	store.Delete(types.PortfolioKey(
		denom,
	))
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
