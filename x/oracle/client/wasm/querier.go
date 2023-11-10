package wasm

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
	case query.OracleParams != nil:
		return oq.queryParams(ctx, query.OracleParams)
	case query.OracleBandPriceResult != nil:
		return oq.queryBandPriceResult(ctx, query.OracleBandPriceResult)
	case query.OracleLastBandRequestId != nil:
		return oq.queryLastBandRequestId(ctx, query.OracleLastBandRequestId)
	case query.OracleAssetInfo != nil:
		return oq.queryAssetInfo(ctx, query.OracleAssetInfo)
	case query.OracleAssetInfoAll != nil:
		return oq.queryAssetInfoAll(ctx, query.OracleAssetInfoAll)
	case query.OraclePriceAll != nil:
		return oq.queryPriceAll(ctx, query.OraclePriceAll)
	case query.OraclePriceFeeder != nil:
		return oq.queryPriceFeeder(ctx, query.OraclePriceFeeder)
	case query.OraclePriceFeederAll != nil:
		return oq.queryPriceFeederAll(ctx, query.OraclePriceFeederAll)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
