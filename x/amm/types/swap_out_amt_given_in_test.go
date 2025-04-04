package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestNormalizedWeights(t *testing.T) {
	for _, tc := range []struct {
		desc        string
		poolAssets  []types.PoolAsset
		poolWeights []types.AssetWeight
	}{
		{
			desc:        "empty assets case",
			poolAssets:  []types.PoolAsset{},
			poolWeights: []types.AssetWeight{},
		},
		{
			desc: "total weight zero case",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Elys, 1000),
					Weight: sdkmath.ZeroInt(),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  ptypes.Elys,
					Weight: elystypes.ZeroDec34(),
				},
			},
		},
		{
			desc: "positive weights with one zero",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Elys, 1000),
					Weight: sdkmath.ZeroInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000),
					Weight: sdkmath.OneInt(),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  ptypes.Elys,
					Weight: elystypes.ZeroDec34(),
				},
				{
					Asset:  ptypes.Eden,
					Weight: elystypes.OneDec34(),
				},
			},
		},
		{
			desc: "all positive weights",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Elys, 1000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000),
					Weight: sdkmath.OneInt(),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  ptypes.Elys,
					Weight: elystypes.NewDec34WithPrec(5, 1),
				},
				{
					Asset:  ptypes.Eden,
					Weight: elystypes.NewDec34WithPrec(5, 1),
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			require.Equal(t, tc.poolWeights, types.NormalizedWeights(tc.poolAssets))
		})
	}
}

func (suite *TestSuite) TestOraclePoolNormalizedWeights() {
	for _, tc := range []struct {
		desc        string
		poolAssets  []types.PoolAsset
		poolWeights []types.AssetWeight
		expError    bool
	}{
		{
			desc: "oracle pool all asset prices set case",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  ptypes.BaseCurrency,
					Weight: elystypes.NewDec34WithPrec(5, 1),
				},
				{
					Asset:  "uusdt",
					Weight: elystypes.NewDec34WithPrec(5, 1),
				},
			},
			expError: false,
		},
		{
			desc: "oracle pool all asset prices set and amount zero case",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 0), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 0), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  ptypes.BaseCurrency,
					Weight: elystypes.NewDec34FromString("0.000000000000000000"),
				},
				{
					Asset:  "uusdt",
					Weight: elystypes.NewDec34FromString("0.000000000000000000"),
				},
			},
			expError: false,
		},
		{
			desc: "oracle pool one asset price not set",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uUSDT", 1000_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			poolWeights: []types.AssetWeight{},
			expError:    true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			now := time.Now()
			suite.ctx = suite.ctx.WithBlockTime(now)

			// prices set for USDT and USDC
			suite.SetupStableCoinPrices()

			// execute function
			weights, err := types.GetOraclePoolNormalizedWeights(suite.ctx, uint64(1), suite.app.OracleKeeper, tc.poolAssets)
			if tc.expError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(weights, tc.poolWeights)
			}
		})
	}
}

func (suite *TestSuite) TestNewPoolAssetsAfterSwap() {
	for _, tc := range []struct {
		desc            string
		poolAssets      []types.PoolAsset
		inCoins         sdk.Coins
		outCoins        sdk.Coins
		poolAssetsAfter []types.PoolAsset
		expErr          bool
	}{
		{
			desc: "positive in and no out case",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			inCoins:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000)},
			outCoins: sdk.Coins{},
			poolAssetsAfter: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1001_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			expErr: false,
		},
		{
			desc: "no in and positive out case",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			inCoins:  sdk.Coins{},
			outCoins: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000)},
			poolAssetsAfter: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 999_000_000), // 999 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			expErr: false,
		},
		{
			desc: "positive in and positive out case",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			inCoins:  sdk.Coins{sdk.NewInt64Coin("uusdt", 1000_000)},
			outCoins: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000)},
			poolAssetsAfter: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 999_000_000), // 999 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1001_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			expErr: false,
		},
		{
			desc: "withdrawing more than pool size",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			inCoins:         sdk.Coins{sdk.NewInt64Coin("uusdt", 1000_000)},
			outCoins:        sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1001_000_000)},
			poolAssetsAfter: []types.PoolAsset{},
			expErr:          true,
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
					UseOracle: true,
				},
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdkmath.ZeroInt(),
			}
			poolAssets, err := pool.NewPoolAssetsAfterSwap(suite.ctx, tc.inCoins, tc.outCoins, tc.poolAssets)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(poolAssets, tc.poolAssetsAfter)
			}
		})
	}
}

func (suite *TestSuite) TestWeightDistanceFromTarget() {
	for _, tc := range []struct {
		desc        string
		poolAssets  []types.PoolAsset
		expDistance elystypes.Dec34
	}{
		{
			desc: "zero balance for one asset",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 0), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			expDistance: elystypes.NewDec34WithPrec(5, 1),
		},
		{
			desc: "zero for all assets",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 0), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 0), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			expDistance: elystypes.NewDec34WithPrec(5, 1),
		},
		{
			desc: "all positive",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			expDistance: elystypes.NewDec34WithPrec(25, 2),
		},
		{
			desc: "zero distance",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
					Weight: sdkmath.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdkmath.NewInt(50),
				},
			},
			expDistance: elystypes.ZeroDec34(),
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
					UseOracle: true,
				},
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdkmath.ZeroInt(),
			}
			distance := pool.WeightDistanceFromTarget(suite.ctx, suite.app.OracleKeeper, tc.poolAssets)
			suite.Require().Equal(distance.String(), tc.expDistance.String())
		})
	}
}

func (suite *TestSuite) TestSwapOutAmtGivenIn() {
	for _, tc := range []struct {
		desc                   string
		poolAssets             []types.PoolAsset
		useOracle              bool
		externalLiquidityRatio sdkmath.LegacyDec
		thresholdWeightDiff    sdkmath.LegacyDec
		tokenIn                sdk.Coin
		outTokenDenom          string
		swapFee                sdkmath.LegacyDec
		expRecoveryBonus       elystypes.Dec34
		expTokenOut            sdk.Coin
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
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),               // 10x
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),    // 20%
			tokenIn:                sdk.NewInt64Coin("uusdt", 100_000_000), // 100 USDC
			outTokenDenom:          ptypes.BaseCurrency,
			swapFee:                sdkmath.LegacyZeroDec(),
			expRecoveryBonus:       elystypes.ZeroDec34(),
			expTokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 99009900),
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
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),               // 10x
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),    // 20%
			tokenIn:                sdk.NewInt64Coin("uusdt", 100_000_000), // 100 USDC
			outTokenDenom:          ptypes.BaseCurrency,
			swapFee:                sdkmath.LegacyZeroDec(),
			expRecoveryBonus:       elystypes.NewDec34FromString("-0.006347556007845347576004236722241546"),
			expTokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 98054944),
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
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),                           // 10x
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),                // 20%
			tokenIn:                sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000_000), // 100 USDC
			outTokenDenom:          "uusdt",
			swapFee:                sdkmath.LegacyZeroDec(),
			expRecoveryBonus:       elystypes.NewDec34FromString("0.001558845726811989564174701707355285"),
			expTokenOut:            sdk.NewInt64Coin("uusdt", 98687060),
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
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),                           // 10x
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),                // 20%
			tokenIn:                sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000_000), // 100 USDC
			outTokenDenom:          "uusdt",
			swapFee:                sdkmath.LegacyNewDecWithPrec(1, 2), // 1%
			expRecoveryBonus:       elystypes.ZeroDec34(),
			expTokenOut:            sdk.NewInt64Coin("uusdt", 247913188),
			expErr:                 false,
		},

		{
			desc: "tokenOut is zero",
			poolAssets: []types.PoolAsset{
				{
					Token:                  sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000_000),
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
				{
					Token:                  sdk.NewInt64Coin("uusdt", 1500_000_000),
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
				},
			},
			useOracle:              false,
			externalLiquidityRatio: sdkmath.LegacyNewDec(10),
			thresholdWeightDiff:    sdkmath.LegacyNewDecWithPrec(20, 2),
			tokenIn:                sdk.NewInt64Coin(ptypes.BaseCurrency, 0),
			outTokenDenom:          "uusdt",
			swapFee:                sdkmath.LegacyNewDecWithPrec(1, 2),
			expRecoveryBonus:       elystypes.ZeroDec34(),
			expTokenOut:            sdk.NewInt64Coin("uusdt", 0),
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
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdkmath.ZeroInt(),
			}
			params := suite.app.AmmKeeper.GetParams(suite.ctx)
			params.ThresholdWeightDifference = tc.thresholdWeightDiff
			params.WeightBreakingFeeMultiplier = sdkmath.LegacyNewDecWithPrec(2, 4)
			params.WeightBreakingFeeExponent = sdkmath.LegacyNewDecWithPrec(25, 1) // 2.5
			params.WeightRecoveryFeePortion = sdkmath.LegacyNewDecWithPrec(50, 2)  // 50%
			tokenOut, _, _, weightBonus, _, _, err := pool.SwapOutAmtGivenIn(suite.ctx, suite.app.OracleKeeper, &pool, sdk.Coins{tc.tokenIn}, tc.outTokenDenom, tc.swapFee, suite.app.AccountedPoolKeeper, sdkmath.LegacyOneDec(), params, sdkmath.LegacyZeroDec())
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().EqualError(err, "token out amount is zero")
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tokenOut.String(), tc.expTokenOut.String())
				suite.Require().Equal(weightBonus.String(), tc.expRecoveryBonus.String())
			}
		})
	}
}
