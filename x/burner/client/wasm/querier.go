package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/burner/keeper"
)

// Querier handles queries for the Burner module.
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
	case query.BurnerParams != nil:
		return oq.queryParams(ctx, query.BurnerParams)
	case query.BurnerHistory != nil:
		return oq.queryHistory(ctx, query.BurnerHistory)
	case query.BurnerHistoryAll != nil:
		return oq.queryHistoryAll(ctx, query.BurnerHistoryAll)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
