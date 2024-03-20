package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/launchpad/keeper"
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
	case query.LaunchpadParams != nil:
		return oq.queryParams(ctx, query.LaunchpadParams)
	case query.LaunchpadBonus != nil:
		return oq.queryBonus(ctx, query.LaunchpadBonus)
	case query.LaunchpadBuyElysEst != nil:
		return oq.queryBuyElysEst(ctx, query.LaunchpadBuyElysEst)
	case query.LaunchpadReturnElysEst != nil:
		return oq.queryReturnElysEst(ctx, query.LaunchpadReturnElysEst)
	case query.LaunchpadOrders != nil:
		return oq.queryOrders(ctx, query.LaunchpadOrders)
	case query.LaunchpadAllOrders != nil:
		return oq.queryAllOrders(ctx, query.LaunchpadAllOrders)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
