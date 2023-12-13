package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/accountedpool/keeper"
)

// Querier handles queries for the Accounted Pool module.
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
	case query.AccountedPoolAccountedPool != nil:
		return oq.queryAccountedPool(ctx, query.AccountedPoolAccountedPool)
	case query.AccountedPoolAccountedPoolAll != nil:
		return oq.queryAccountedPoolAll(ctx, query.AccountedPoolAccountedPoolAll)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
