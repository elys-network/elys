package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/commitment/keeper"
)

// Querier handles queries for the Commitment module.
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
	case query.CommitmentParams != nil:
		return oq.queryParams(ctx, query.CommitmentParams)
	case query.CommitmentShowCommitments != nil:
		return oq.queryShowCommitments(ctx, query.CommitmentShowCommitments)
	case query.Delegations != nil:
		return oq.queryDelegations(ctx, query.Delegations)
	case query.UnbondingDelegations != nil:
		return oq.queryUnbondingDelegations(ctx, query.UnbondingDelegations)
	case query.StakedBalanceOfDenom != nil:
		return oq.queryStakedBalanceOfDenom(ctx, query.StakedBalanceOfDenom)
	case query.RewardsBalanceOfDenom != nil:
		return oq.queryRewardBalanceOfDenom(ctx, query.RewardsBalanceOfDenom)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
