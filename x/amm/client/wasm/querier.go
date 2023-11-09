package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/amm/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
)

// Querier handles queries for the AMM module.
type Querier struct {
	keeper           *keeper.Keeper
	bankKeeper       *bankkeeper.BaseKeeper
	commitmentKeeper *commitmentkeeper.Keeper
}

func NewQuerier(keeper *keeper.Keeper, bankKeeper *bankkeeper.BaseKeeper, commitmentKeeper *commitmentkeeper.Keeper) *Querier {
	return &Querier{
		keeper:           keeper,
		bankKeeper:       bankKeeper,
		commitmentKeeper: commitmentKeeper,
	}
}

func (oq *Querier) HandleQuery(ctx sdk.Context, query wasmbindingstypes.ElysQuery) ([]byte, error) {
	switch {
	case query.PriceAll != nil:
		return oq.querySwapEstimation(ctx, query.SwapEstimation)
	case query.BalanceOfDenom != nil:
		return oq.queryBalanceOfDenom(ctx, query.BalanceOfDenom)
	case query.Params != nil:
		return oq.queryParams(ctx, query.Params)
	case query.Pool != nil:
		return oq.queryPool(ctx, query.Pool)
	case query.PoolAll != nil:
		return oq.queryPoolAll(ctx, query.PoolAll)
	case query.DenomLiquidity != nil:
		return oq.queryDenomLiquidity(ctx, query.DenomLiquidity)
	case query.DenomLiquidityAll != nil:
		return oq.queryDenomLiquidityAll(ctx, query.DenomLiquidityAll)
	case query.SlippageTrack != nil:
		return oq.querySlippageTrack(ctx, query.SlippageTrack)
	case query.SlippageTrackAll != nil:
		return oq.querySlippageTrackAll(ctx, query.SlippageTrackAll)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
