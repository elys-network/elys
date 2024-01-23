package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/assetprofile/keeper"
)

// Querier handles queries for the Asset Profile module.
type Querier struct {
	keeper *keeper.Keeper
}

func NewQuerier(keeper *keeper.Keeper) *Querier {
	return &Querier{
		keeper: keeper,
	}
}

func (oq *Querier) HandleQuery(ctx sdk.Context, query wasmbindingstypes.ElysQuery) ([]byte, error) {
	switch {
	case query.AssetProfileParams != nil:
		return oq.queryParams(ctx, query.AssetProfileParams)
	case query.AssetProfileEntry != nil:
		return oq.queryEntry(ctx, query.AssetProfileEntry)
	case query.AssetProfileEntryByDenom != nil:
		return oq.queryEntryByDenom(ctx, query.AssetProfileEntryByDenom)
	case query.AssetProfileEntryAll != nil:
		return oq.queryEntryAll(ctx, query.AssetProfileEntryAll)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
