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
	case query.CommitmentDelegations != nil:
		return oq.queryDelegations(ctx, query.CommitmentDelegations)
	case query.CommitmentUnbondingDelegations != nil:
		return oq.queryUnbondingDelegations(ctx, query.CommitmentUnbondingDelegations)
	case query.CommitmentStakedBalanceOfDenom != nil:
		return oq.queryStakedBalanceOfDenom(ctx, query.CommitmentStakedBalanceOfDenom)
	case query.CommitmentRewardsBalanceOfDenom != nil:
		return oq.queryRewardBalanceOfDenom(ctx, query.CommitmentRewardsBalanceOfDenom)
	case query.CommitmentAllValidators != nil:
		return oq.queryAllValidators(ctx, query.CommitmentAllValidators)
	case query.CommitmentDelegatorValidators != nil:
		return oq.queryDelegatorValidators(ctx, query.CommitmentDelegatorValidators)
	case query.CommitmentStakedPositions != nil:
		return oq.queryStakedPositions(ctx, query.CommitmentStakedPositions)
	case query.CommitmentUnStakedPositions != nil:
		return oq.queryUnStakedPositions(ctx, query.CommitmentUnStakedPositions)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
