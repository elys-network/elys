package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	accountedpoolkeeper "github.com/elys-network/elys/x/accountedpool/keeper"
	"github.com/elys-network/elys/x/amm/keeper"
	assetprofilekeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
)

// Querier handles queries for the AMM module.
type Querier struct {
	keeper              *keeper.Keeper
	bankKeeper          *bankkeeper.BaseKeeper
	commitmentKeeper    *commitmentkeeper.Keeper
	assetProfileKeeper  *assetprofilekeeper.Keeper
	perpetualKeeper     *perpetualkeeper.Keeper
	incentiveKeeper     *incentivekeeper.Keeper
	oraclekeeper        *oraclekeeper.Keeper
	leveragelpKeeper    *leveragelpkeeper.Keeper
	accountedpoolKeeper *accountedpoolkeeper.Keeper
	stablestakeKeeper   *stablestakekeeper.Keeper
}

func NewQuerier(
	keeper *keeper.Keeper,
	bankKeeper *bankkeeper.BaseKeeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	assetProfileKeeper *assetprofilekeeper.Keeper,
	perpetualKeeper *perpetualkeeper.Keeper,
	incentiveKeeper *incentivekeeper.Keeper,
	oraclekeeper *oraclekeeper.Keeper,
	leveragelpKeeper *leveragelpkeeper.Keeper,
	accountedpoolKeeper *accountedpoolkeeper.Keeper,
	stablestakeKeeper *stablestakekeeper.Keeper) *Querier {
	return &Querier{
		keeper:              keeper,
		bankKeeper:          bankKeeper,
		commitmentKeeper:    commitmentKeeper,
		assetProfileKeeper:  assetProfileKeeper,
		perpetualKeeper:     perpetualKeeper,
		incentiveKeeper:     incentiveKeeper,
		oraclekeeper:        oraclekeeper,
		leveragelpKeeper:    leveragelpKeeper,
		accountedpoolKeeper: accountedpoolKeeper,
		stablestakeKeeper:   stablestakeKeeper,
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
	case query.AmmSwapEstimation != nil:
		return oq.querySwapEstimation(ctx, query.AmmSwapEstimation)
	case query.AmmSwapEstimationByDenom != nil:
		return oq.querySwapEstimationByDenom(ctx, query.AmmSwapEstimationByDenom)
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
	case query.AmmPriceByDenom != nil:
		return oq.queryAmmPriceByDenom(ctx, query.AmmPriceByDenom)
	case query.AmmEarnMiningPoolAll != nil:
		return oq.queryEarnMiningPoolAll(ctx, query.AmmEarnMiningPoolAll)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}
