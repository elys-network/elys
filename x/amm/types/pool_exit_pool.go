package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// return exitingCoins, weightBalanceBonus, slippage, swapFee, slippageCoins, nil
func (p *Pool) ExitPool(ctx sdk.Context, oracleKeeper OracleKeeper, accountedPoolKeeper AccountedPoolKeeper, snapshot Pool, exitingShares math.Int, tokenOutDenom string, params Params, takerFees osmomath.BigDec, applyWeightBreakingFee bool) (exitingCoins sdk.Coins, weightBalanceBonus osmomath.BigDec, slippage osmomath.BigDec, swapFee osmomath.BigDec, takerFeesFinal osmomath.BigDec, slippageCoins sdk.Coins, err error) {
	exitingCoins, weightBalanceBonus, slippage, swapFee, takerFeesFinal, slippageCoins, err = p.CalcExitPoolCoinsFromShares(ctx, oracleKeeper, accountedPoolKeeper, snapshot, exitingShares, tokenOutDenom, params, takerFees, applyWeightBreakingFee)
	if err != nil {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, err
	}

	if err := p.processExitPool(ctx, exitingCoins, exitingShares); err != nil {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, err
	}

	return exitingCoins, weightBalanceBonus, slippage, swapFee, takerFeesFinal, slippageCoins, nil
}

// exitPool exits the pool given exitingCoins and exitingShares.
// updates the pool's liquidity and totalShares.
func (p *Pool) processExitPool(_ sdk.Context, exitingCoins sdk.Coins, exitingShares math.Int) error {
	balances := p.GetTotalPoolLiquidity().Sub(exitingCoins...)
	if err := p.UpdatePoolAssetBalances(balances); err != nil {
		return err
	}

	totalShares := p.GetTotalShares().Amount
	p.TotalShares = sdk.NewCoin(p.TotalShares.Denom, totalShares.Sub(exitingShares))

	return nil
}
