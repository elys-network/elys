package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	assetkeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	"github.com/elys-network/elys/x/commitment/keeper"
	epochkeeper "github.com/elys-network/elys/x/epochs/keeper"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
)

// Querier handles queries for the Commitment module.
type Querier struct {
	keeper        *keeper.Keeper
	stakingKeeper *stakingkeeper.Keeper
	epochKeeper   *epochkeeper.Keeper
	ammKeeper     *ammkeeper.Keeper
	assetKeeper   *assetkeeper.Keeper
	stableKeeper  *stablekeeper.Keeper
	oracleKeeper  *oraclekeeper.Keeper
}

func NewQuerier(
	keeper *keeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	epochKeeper *epochkeeper.Keeper,
	ammKeeper *ammkeeper.Keeper,
	assetKeeper *assetkeeper.Keeper,
	stableKeeper *stablekeeper.Keeper,
	oracleKeeper *oraclekeeper.Keeper,
) *Querier {
	return &Querier{
		keeper:        keeper,
		stakingKeeper: stakingKeeper,
		epochKeeper:   epochKeeper,
		ammKeeper:     ammKeeper,
		assetKeeper:   assetKeeper,
		stableKeeper:  stableKeeper,
		oracleKeeper:  oracleKeeper,
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
	case query.CommitmentAllValidators != nil:
		return oq.queryAllValidators(ctx, query.CommitmentAllValidators)
	case query.CommitmentDelegatorValidators != nil:
		return oq.queryDelegatorValidators(ctx, query.CommitmentDelegatorValidators)
	case query.CommitmentStakedPositions != nil:
		return oq.queryStakedPositions(ctx, query.CommitmentStakedPositions)
	case query.CommitmentUnStakedPositions != nil:
		return oq.queryUnStakedPositions(ctx, query.CommitmentUnStakedPositions)
	case query.CommitmentVestingInfo != nil:
		return oq.queryVestingInfo(ctx, query.CommitmentVestingInfo)
	case query.CommitmentNumberOfCommitments != nil:
		return oq.queryNumberOfCommitments(ctx, query.CommitmentNumberOfCommitments)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
