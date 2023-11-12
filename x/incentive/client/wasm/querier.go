package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/incentive/keeper"
)

// Querier handles queries for the Incentive module.
type Querier struct {
	keeper        *keeper.Keeper
	stakingKeeper *stakingkeeper.Keeper
}

func NewQuerier(keeper *keeper.Keeper, stakingKeeper *stakingkeeper.Keeper) *Querier {
	return &Querier{
		keeper:        keeper,
		stakingKeeper: stakingKeeper,
	}
}

func (oq *Querier) HandleQuery(ctx sdk.Context, query wasmbindingstypes.ElysQuery) ([]byte, error) {
	switch {
	case query.IncentiveParams != nil:
		return oq.queryParams(ctx, query.IncentiveParams)
	case query.IncentiveCommunityPool != nil:
		return oq.queryCommunityPool(ctx, query.IncentiveCommunityPool)
	case query.AllValidators != nil:
		return oq.queryAllValidators(ctx, query.AllValidators)
	case query.DelegatorValidators != nil:
		return oq.queryDelegatorValidators(ctx, query.DelegatorValidators)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
