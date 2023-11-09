package querier

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/oracle/keeper"
)

// Querier handles queries for the Oracle module.
type Querier struct {
	keeper *keeper.Keeper
}

func NewQuerier(keeper *keeper.Keeper) *Querier {
	return &Querier{keeper: keeper}
}

func (oq *Querier) HandleQuery(ctx sdk.Context, query wasmbindingstypes.ElysQuery) ([]byte, error) {
	switch {
	case query.PriceAll != nil:
		return oq.queryPriceAll(ctx, query.PriceAll)
	case query.AssetInfo != nil:
		return oq.queryAssetInfo(ctx, query.AssetInfo)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
