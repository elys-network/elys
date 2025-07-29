package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/v7/wasmbindings/types"
	"github.com/elys-network/elys/v7/x/amm/keeper"
)

// Querier handles queries for the AMM module.
type Querier struct {
	keeper *keeper.Keeper
}

func NewQuerier(
	keeper *keeper.Keeper,
) *Querier {
	return &Querier{
		keeper: keeper,
	}
}

func (oq *Querier) HandleQuery(ctx sdk.Context, query wasmbindingstypes.ElysQuery) ([]byte, error) {
	switch {
	case query.AmmParams != nil:
		return oq.queryParams(ctx, query.AmmParams)
	case query.AmmPool != nil:
		return oq.queryPool(ctx, query.AmmPool)
	case query.AmmPoolAll != nil:
		return oq.queryPoolAll(ctx, query.AmmPoolAll)
	case query.AmmDenomLiquidity != nil:
		return oq.queryDenomLiquidity(ctx, query.AmmDenomLiquidity)
	case query.AmmDenomLiquidityAll != nil:
		return oq.queryDenomLiquidityAll(ctx, query.AmmDenomLiquidityAll)
	case query.AmmSwapEstimationExactAmountOut != nil:
		return oq.QuerySwapEstimationExactAmountOut(ctx, query.AmmSwapEstimationExactAmountOut)
	case query.AmmSwapEstimation != nil:
		return oq.querySwapEstimation(ctx, query.AmmSwapEstimation)
	case query.AmmSwapEstimationByDenom != nil:
		return oq.querySwapEstimationByDenom(ctx, query.AmmSwapEstimationByDenom)
	case query.AmmJoinPoolEstimation != nil:
		return oq.queryJoinPoolEstimation(ctx, query.AmmJoinPoolEstimation)
	case query.AmmExitPoolEstimation != nil:
		return oq.queryExitPoolEstimation(ctx, query.AmmExitPoolEstimation)
	case query.AmmSlippageTrack != nil:
		return oq.querySlippageTrack(ctx, query.AmmSlippageTrack)
	case query.AmmSlippageTrackAll != nil:
		return oq.querySlippageTrackAll(ctx, query.AmmSlippageTrackAll)
	case query.AmmBalance != nil:
		return oq.queryBalance(ctx, query.AmmBalance)
	case query.AmmInRouteByDenom != nil:
		return oq.queryInRouteByDenom(ctx, query.AmmInRouteByDenom)
	case query.AmmOutRouteByDenom != nil:
		return oq.queryOutRouteByDenom(ctx, query.AmmOutRouteByDenom)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
