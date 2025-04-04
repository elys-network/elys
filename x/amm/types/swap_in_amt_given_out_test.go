package types_test

import (
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *TestSuite) TestSwapInAmtGivenOut() {
	for _, tc := range []struct {
		desc                   string
		poolAssets             []types.PoolAsset
		useOracle              bool
		externalLiquidityRatio sdkmath.LegacyDec
		thresholdWeightDiff    sdkmath.LegacyDec
		tokenOut               sdk.Coin
		inTokenDenom           string
		swapFee                osmomath.BigDec
		expRecoveryBonus       osmomath.BigDec
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
					Token:                  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
				{
					Token:                  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
			},
			useOracle:              true,
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),                           // 10x
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),                // 20%
			tokenOut:               sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000_000), // 100 USDC
			inTokenDenom:           "uusdt",
			swapFee:                osmomath.ZeroBigDec(),
			expRecoveryBonus:       osmomath.MustNewBigDecFromStr("0"),
			expTokenIn:             sdk.NewInt64Coin("uusdt", 101010110),
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
					Token:                  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000), // 1000 USDT
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
				{
					Token:                  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
			},
			useOracle:              true,
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),                           // 10x
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),                // 20%
			tokenOut:               sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000_000), // 100 USDC
			inTokenDenom:           "uusdt",
			swapFee:                osmomath.ZeroBigDec(),
			expRecoveryBonus:       osmomath.MustNewBigDecFromStr("-0.006413552900341021378891128886988564"),
			expTokenIn:             sdk.NewInt64Coin("uusdt", 102008668),
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
					Token:                  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000), // 1000 USDT
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
				{
					Token:                  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
			},
			useOracle:              true,
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),               // 10x
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),    // 20%
			tokenOut:               sdk.NewInt64Coin("uusdt", 100_000_000), // 100 USDC
			inTokenDenom:           ptypes.BaseCurrency,
			swapFee:                osmomath.ZeroBigDec(),
			expRecoveryBonus:       osmomath.MustNewBigDecFromStr("0.001558845726811989564174701707355285"),
			expTokenIn:             sdk.NewInt64Coin(ptypes.BaseCurrency, 101348300),
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
					Token:                  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000), // 1000 USDT
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
				{
					Token:                  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
			},
			useOracle:              false,
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),               // 10x
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),    // 20%
			tokenOut:               sdk.NewInt64Coin("uusdt", 100_000_000), // 100 USDC
			inTokenDenom:           ptypes.BaseCurrency,
			swapFee:                osmomath.NewBigDecWithPrec(1, 2), // 1%
			expRecoveryBonus:       osmomath.ZeroBigDec(),
			expTokenIn:             sdk.NewInt64Coin(ptypes.BaseCurrency, 36075037),
			expErr:                 false,
		},

		{
			desc: "tokenOut is zero",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000),
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000),
					Weight: sdkmath.NewInt(50),
				},
			},
			useOracle:              false,
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),
			tokenOut:               sdk.NewInt64Coin("uusdt", 0),
			inTokenDenom:           ptypes.BaseCurrency,
			swapFee:                osmomath.NewBigDecWithPrec(1, 2),
			expRecoveryBonus:       osmomath.ZeroBigDec(),
			expTokenIn:             sdk.NewInt64Coin(ptypes.BaseCurrency, 0),
			expErr:                 true,
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
					SwapFee:   sdkmath.LegacyZeroDec(),
					UseOracle: tc.useOracle,
					FeeDenom:  ptypes.BaseCurrency,
				},
				TotalShares: sdk.NewCoin(types.GetPoolShareDenom(1), sdkmath.ZeroInt()),
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdkmath.ZeroInt(),
			}
			params := suite.app.AmmKeeper.GetParams(suite.ctx)
			params.ThresholdWeightDifference = tc.thresholdWeightDiff
			params.WeightBreakingFeeMultiplier = sdkmath.LegacyNewDecWithPrec(2, 4)
			params.WeightBreakingFeeExponent = sdkmath.LegacyNewDecWithPrec(25, 1) // 2.5
			params.WeightRecoveryFeePortion = sdkmath.LegacyNewDecWithPrec(50, 2)  // 50%
			tokenOut, _, _, weightBonus, _, _, err := pool.SwapInAmtGivenOut(suite.ctx, suite.app.OracleKeeper, &pool, sdk.Coins{tc.tokenOut}, tc.inTokenDenom, tc.swapFee, suite.app.AccountedPoolKeeper, osmomath.OneBigDec(), params, osmomath.ZeroBigDec())
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().EqualError(err, "amount too low")
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tokenOut.String(), tc.expTokenIn.String())
				suite.Require().Equal(weightBonus.String(), tc.expRecoveryBonus.String())
			}
		})
	}
}
