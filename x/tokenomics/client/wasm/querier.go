package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/tokenomics/keeper"
)

// Querier handles queries for the Tokenomics module.
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
	case query.TokenomicsParams != nil:
		return oq.queryParams(ctx, query.TokenomicsParams)
	case query.TokenomicsAirdrop != nil:
		return oq.queryAirdrop(ctx, query.TokenomicsAirdrop)
	case query.TokenomicsAirdropAll != nil:
		return oq.queryAirdropAll(ctx, query.TokenomicsAirdropAll)
	case query.TokenomicsGenesisInflation != nil:
		return oq.queryGenesisInflation(ctx, query.TokenomicsGenesisInflation)
	case query.TokenomicsTimeBasedInflation != nil:
		return oq.queryTimeBasedInflation(ctx, query.TokenomicsTimeBasedInflation)
	case query.TokenomicsTimeBasedInflationAll != nil:
		return oq.queryTimeBasedInflationAll(ctx, query.TokenomicsTimeBasedInflationAll)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
