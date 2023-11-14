package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
)

// Querier handles queries for the Stable Stake module.
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
	case query.StableStakeParams != nil:
		return oq.queryParams(ctx, query.StableStakeParams)
	case query.StableStakeBorrowRatio != nil:
		return oq.queryBorrowRatio(ctx, query.StableStakeBorrowRatio)
	case query.BalanceOfBorrow != nil:
		return oq.queryBorrowedAmount(ctx, query.BalanceOfBorrow)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
