package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
)

// Querier handles queries for the Perpetual module.
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
	case query.PerpetualParams != nil:
		return oq.queryParams(ctx, query.PerpetualParams)
	case query.PerpetualQueryPositions != nil:
		return oq.queryPositions(ctx, query.PerpetualQueryPositions)
	case query.PerpetualQueryPositionsByPool != nil:
		return oq.queryPositionsByPool(ctx, query.PerpetualQueryPositionsByPool)
	case query.PerpetualGetStatus != nil:
		return oq.queryGetStatus(ctx, query.PerpetualGetStatus)
	case query.PerpetualGetPositionsForAddress != nil:
		return oq.queryPositionsForAddress(ctx, query.PerpetualGetPositionsForAddress)
	case query.PerpetualGetWhitelist != nil:
		return oq.queryGetWhitelist(ctx, query.PerpetualGetWhitelist)
	case query.PerpetualIsWhitelisted != nil:
		return oq.queryIsWhitelisted(ctx, query.PerpetualIsWhitelisted)
	case query.PerpetualPool != nil:
		return oq.queryPool(ctx, query.PerpetualPool)
	case query.PerpetualPools != nil:
		return oq.queryPools(ctx, query.PerpetualPools)
	case query.PerpetualMTP != nil:
		return oq.queryMtp(ctx, query.PerpetualMTP)
	case query.PerpetualOpenEstimation != nil:
		return oq.queryOpenEstimation(ctx, query.PerpetualOpenEstimation)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
