package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/leveragelp/keeper"
)

// Querier handles queries for the Leverage LP module.
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
	case query.LeveragelpParams != nil:
		return oq.queryParams(ctx, query.LeveragelpParams)
	case query.LeveragelpQueryPositions != nil:
		return oq.queryPositions(ctx, query.LeveragelpQueryPositions)
	case query.LeveragelpQueryPositionsByPool != nil:
		return oq.queryPositionsByPool(ctx, query.LeveragelpQueryPositionsByPool)
	case query.LeveragelpGetStatus != nil:
		return oq.queryGetStatus(ctx, query.LeveragelpGetStatus)
	case query.LeveragelpQueryPositionsForAddress != nil:
		return oq.queryPositionsForAddress(ctx, query.LeveragelpQueryPositionsForAddress)
	case query.LeveragelpGetWhitelist != nil:
		return oq.queryGetWhitelist(ctx, query.LeveragelpGetWhitelist)
	case query.LeveragelpIsWhitelisted != nil:
		return oq.queryIsWhitelisted(ctx, query.LeveragelpIsWhitelisted)
	case query.LeveragelpPool != nil:
		return oq.queryPool(ctx, query.LeveragelpPool)
	case query.LeveragelpPools != nil:
		return oq.queryPools(ctx, query.LeveragelpPools)
	case query.LeveragelpPosition != nil:
		return oq.queryPosition(ctx, query.LeveragelpPosition)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
