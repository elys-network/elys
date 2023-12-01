package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/margin/keeper"
)

// Querier handles queries for the Margin module.
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
	case query.MarginParams != nil:
		return oq.queryParams(ctx, query.MarginParams)
	case query.MarginQueryPositions != nil:
		return oq.queryPositions(ctx, query.MarginQueryPositions)
	case query.MarginQueryPositionsByPool != nil:
		return oq.queryPositionsByPool(ctx, query.MarginQueryPositionsByPool)
	case query.MarginGetStatus != nil:
		return oq.queryGetStatus(ctx, query.MarginGetStatus)
	case query.MarginGetPositionsForAddress != nil:
		return oq.queryPositionsForAddress(ctx, query.MarginGetPositionsForAddress)
	case query.MarginGetWhitelist != nil:
		return oq.queryGetWhitelist(ctx, query.MarginGetWhitelist)
	case query.MarginIsWhitelisted != nil:
		return oq.queryIsWhitelisted(ctx, query.MarginIsWhitelisted)
	case query.MarginPool != nil:
		return oq.queryPool(ctx, query.MarginPool)
	case query.MarginPools != nil:
		return oq.queryPools(ctx, query.MarginPools)
	case query.MarginMTP != nil:
		return oq.queryMtp(ctx, query.MarginMTP)
	case query.MarginMinCollateral != nil:
		return oq.queryMinCollateral(ctx, query.MarginMinCollateral)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
