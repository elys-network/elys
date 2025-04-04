package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
)

// Querier handles queries for the Auth module.
type Querier struct {
	keeper *keeper.AccountKeeper
}

func NewQuerier(keeper *keeper.AccountKeeper) *Querier {
	return &Querier{
		keeper: keeper,
	}
}

func (oq *Querier) HandleQuery(ctx sdk.Context, query wasmbindingstypes.ElysQuery) ([]byte, error) {
	switch {
	case query.AuthAddresses != nil:
		return oq.queryAddresses(ctx, query.AuthAddresses)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
