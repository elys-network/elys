package app

import (
	"fmt"

	"cosmossdk.io/core/appmodule"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/ojo-network/ojo/x/oracle/abci"
	"github.com/ojo-network/ojo/x/oracle/types"
)

// PreBlocker is run before finalize block to update the aggregrate exchange rate votes on the oracle module
// that were verified by the vote etension handler so that the exchange rate votes are available during the
// entire block execution (from BeginBlock). It will execute the preblockers of the other modules set in
// SetOrderPreBlockers as well.
func (app *ElysApp) PreBlocker(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	if req == nil {
		err := fmt.Errorf("preblocker received a nil request")
		app.Logger().Error(err.Error())
		return nil, err
	}

	// execute preblockers of modules in OrderPreBlockers first.
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	paramsChanged := false
	for _, moduleName := range app.mm.OrderPreBlockers {
		if module, ok := app.mm.Modules[moduleName].(appmodule.HasPreBlocker); ok {
			rsp, err := module.PreBlock(ctx)
			if err != nil {
				return nil, err
			}
			if rsp.IsConsensusParamsChanged() {
				paramsChanged = true
			}
		}
	}

	res := &sdk.ResponsePreBlock{
		ConsensusParamsChanged: paramsChanged,
	}

	if len(req.Txs) == 0 {
		return res, nil
	}
	voteExtensionsEnabled := abci.VoteExtensionsEnabled(ctx)
	if voteExtensionsEnabled {
		var injectedVoteExtTx types.InjectedVoteExtensionTx
		if err := injectedVoteExtTx.Unmarshal(req.Txs[0]); err != nil {
			app.Logger().Error("failed to decode injected vote extension tx", "err", err)
			return nil, err
		}
		for _, exchangeRateVote := range injectedVoteExtTx.ExchangeRateVotes {
			app.OracleKeeper.SetAggregateExchangeRateVote(ctx, exchangeRateVote.Voter, exchangeRateVote)
		}

		for _, externalLiquidity := range injectedVoteExtTx.ExternalLiquidity {
			pool, found := app.AmmKeeper.GetPool(ctx, externalLiquidity.PoolId)
			if !found {
				app.Logger().Error("failed to get pool", "poolId", externalLiquidity.PoolId, "err", ammtypes.ErrInvalidPoolId)
				continue
			}

			assetAmountDepthInfo := make([]ammtypes.AssetAmountDepth, 0)
			for _, amountDepthInfo := range externalLiquidity.AmountDepthInfo {
				assetAmountDepthInfo = append(assetAmountDepthInfo, ammtypes.AssetAmountDepth{
					Asset:  amountDepthInfo.Asset,
					Amount: amountDepthInfo.Amount,
					Depth:  amountDepthInfo.Depth,
				})
			}

			// Get external liquidity ratio for each of the asset separately
			poolAssets, err := app.AmmKeeper.GetExternalLiquidityRatio(ctx, pool, assetAmountDepthInfo)
			if err != nil {
				app.Logger().Error("failed to get pool assets", "poolId", pool.PoolId, "err", ammtypes.ErrInvalidPoolId)
				continue
			}

			pool.PoolAssets = poolAssets
			app.AmmKeeper.SetPool(ctx, pool)
		}

	}

	app.Logger().Info(
		"preblocker executed",
		"vote_extensions_enabled", voteExtensionsEnabled,
	)

	return res, nil
}
