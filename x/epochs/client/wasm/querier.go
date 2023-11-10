package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/epochs/keeper"
)

// Querier handles queries for the Epochs module.
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
	case query.EpochsEpochInfos != nil:
		return oq.queryEpochInfos(ctx, query.EpochsEpochInfos)
	case query.EpochsCurrentEpoch != nil:
		return oq.queryCurrentEpoch(ctx, query.EpochsCurrentEpoch)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
