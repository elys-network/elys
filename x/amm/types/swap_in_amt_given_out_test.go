package types_test

import (
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *TestSuite) TestSwapInAmtGivenOut() {
	for _, tc := range []struct {
		desc                   string
		poolAssets             []types.PoolAsset
		useOracle              bool
		externalLiquidityRatio sdk.Dec
		thresholdWeightDiff    sdk.Dec
		tokenOut               sdk.Coin
		inTokenDenom           string
		swapFee                sdk.Dec
		expRecoveryBonus       sdk.Dec
		expTokenIn             sdk.Coin
		expErr                 bool
	}{
		// scenario1 - oracle based
		// - USDT/USDC pool
		// - USDT price $1
		// - $1000 USDT / $1000 USDC
		// - External liquidity on Osmosis $10000 USDT / $10000 USDC
		// - Slippage reduction: 90%
		// - Target Weight 50% / 50%
		// - Weight Broken Threshold: 20%
		// - Swap USDT -> USDC
		// - Check weight breaking fee
		// - Check slippage
		// - Check bonus being zero
		{
			desc: "oracle pool scenario1",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle:              true,
			externalLiquidityRatio: sdk.NewDec(10),                                     // 10x
			thresholdWeightDiff:    sdk.NewDecWithPrec(20, 2),                          // 20%
			tokenOut:               sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000_000), // 100 USDC
			inTokenDenom:           "uusdt",
			swapFee:                sdk.ZeroDec(),
			expRecoveryBonus:       sdk.ZeroDec(),
			expTokenIn:             sdk.NewInt64Coin("uusdt", 105108336),
			expErr:                 false,
		},
		// scenario2 - oracle based
		// - USDT/USDC pool
		// - USDT price $1
		// - $500 USDT / $1500 USDC
		// - External liquidity on Osmosis $10000 USDT / $10000 USDC
		// - Slippage reduction: 90%
		// - Target Weight 50% / 50%
		// - Weight Broken Threshold: 20%
		// - Swap USDT -> USDC
		// - Check weight breaking fee zero
		// - Check slippage
		// - Check bonus not be zero
		{
			desc: "oracle pool scenario2",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle:              true,
			externalLiquidityRatio: sdk.NewDec(10),                                     // 10x
			thresholdWeightDiff:    sdk.NewDecWithPrec(20, 2),                          // 20%
			tokenOut:               sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000_000), // 100 USDC
			inTokenDenom:           "uusdt",
			swapFee:                sdk.ZeroDec(),
			expRecoveryBonus:       sdk.ZeroDec(),
			expTokenIn:             sdk.NewInt64Coin("uusdt", 105143571),
			expErr:                 false,
		},
		// scenario3 - oracle based
		// - USDT/USDC pool
		// - USDT price $1
		// - $500 USDT / $1500 USDC
		// - External liquidity on Osmosis $10000 USDT / $10000 USDC
		// - Slippage reduction: 90%
		// - Target Weight 50% / 50%
		// - Weight Broken Threshold: 20%
		// - Swap USDC -> USDT
		// - Check weight breaking fee
		// - Check slippage
		// - Check bonus be zero
		{
			desc: "oracle pool scenario3",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle:              true,
			externalLiquidityRatio: sdk.NewDec(10),                         // 10x
			thresholdWeightDiff:    sdk.NewDecWithPrec(20, 2),              // 20%
			tokenOut:               sdk.NewInt64Coin("uusdt", 100_000_000), // 100 USDC
			inTokenDenom:           ptypes.BaseCurrency,
			swapFee:                sdk.ZeroDec(),
			expRecoveryBonus:       sdk.MustNewDecFromStr("0.050047187318866899"),
			expTokenIn:             sdk.NewInt64Coin(ptypes.BaseCurrency, 100134830),
			expErr:                 false,
		},
		// scenario1 - non-oracle based
		// - USDT/USDC pool
		// - USDT price $1
		// - $500 USDT / $1500 USDC
		// - External liquidity on Osmosis $10000 USDT / $10000 USDC
		// - Slippage reduction: 90%
		// - Target Weight 50% / 50%
		// - Weight Broken Threshold: 20%
		// - Swap USDC -> USDT
		// - Check weight breaking fee
		// - Check slippage
		// - Check bonus be zero
		{
			desc: "non-oracle pool scenario1",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle:              false,
			externalLiquidityRatio: sdk.NewDec(10),                         // 10x
			thresholdWeightDiff:    sdk.NewDecWithPrec(20, 2),              // 20%
			tokenOut:               sdk.NewInt64Coin("uusdt", 100_000_000), // 100 USDC
			inTokenDenom:           ptypes.BaseCurrency,
			swapFee:                sdk.NewDecWithPrec(1, 2), // 1%
			expRecoveryBonus:       sdk.ZeroDec(),
			expTokenIn:             sdk.NewInt64Coin(ptypes.BaseCurrency, 36075037),
			expErr:                 false,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			now := time.Now()
			suite.ctx = suite.ctx.WithBlockTime(now)

			// bootstrap accounts
			poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// prices set for USDT and USDC
			suite.SetupStableCoinPrices()

			// execute function
			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: types.PoolParams{
					SwapFee:                     sdk.ZeroDec(),
					UseOracle:                   tc.useOracle,
					ExternalLiquidityRatio:      tc.externalLiquidityRatio,
					ThresholdWeightDifference:   tc.thresholdWeightDiff,
					WeightBreakingFeeMultiplier: sdk.OneDec(),
					WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
				},
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdk.ZeroInt(),
			}
			tokenOut, _, weightBonus, err := pool.SwapInAmtGivenOut(suite.ctx, suite.app.OracleKeeper, &pool, sdk.Coins{tc.tokenOut}, tc.inTokenDenom, tc.swapFee, suite.app.AccountedPoolKeeper)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tokenOut.String(), tc.expTokenIn.String())
				suite.Require().Equal(weightBonus.String(), tc.expRecoveryBonus.String())
			}
		})
	}
}
