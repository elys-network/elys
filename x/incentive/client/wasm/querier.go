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
	case query.IncentiveApr != nil:
		return oq.queryApr(ctx, query.IncentiveApr)
	case query.IncentiveAprs != nil:
		return oq.queryAprs(ctx, query.IncentiveAprs)
	case query.IncentiveAllProgramRewards != nil:
		return oq.queryAllProgramRewards(ctx, query.IncentiveAllProgramRewards)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
