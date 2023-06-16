package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestNormalizedWeights(t *testing.T) {
	for _, tc := range []struct {
		desc        string
		poolAssets  []*types.PoolAsset
		poolWeights []types.AssetWeight
	}{
		{
			desc:        "empty assets case",
			poolAssets:  []*types.PoolAsset{},
			poolWeights: []types.AssetWeight{},
		},
		{
			desc: "total weight zero case",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uelys", 1000),
					Weight: sdk.ZeroInt(),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  "uelys",
					Weight: sdk.ZeroDec(),
				},
			},
		},
		{
			desc: "positive weights with one zero",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uelys", 1000),
					Weight: sdk.ZeroInt(),
				},
				{
					Token:  sdk.NewInt64Coin("ueden", 1000),
					Weight: sdk.OneInt(),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  "uelys",
					Weight: sdk.ZeroDec(),
				},
				{
					Asset:  "ueden",
					Weight: sdk.OneDec(),
				},
			},
		},
		{
			desc: "all positive weights",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uelys", 1000),
					Weight: sdk.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin("ueden", 1000),
					Weight: sdk.OneInt(),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  "uelys",
					Weight: sdk.NewDecWithPrec(5, 1),
				},
				{
					Asset:  "ueden",
					Weight: sdk.NewDecWithPrec(5, 1),
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
		poolAssets  []*types.PoolAsset
		poolWeights []types.AssetWeight
		expError    bool
	}{
		{
			desc: "oracle pool all asset prices set case",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  "uusdc",
					Weight: sdk.NewDecWithPrec(5, 1),
				},
				{
					Asset:  "uusdt",
					Weight: sdk.NewDecWithPrec(5, 1),
				},
			},
			expError: false,
		},
		{
			desc: "oracle pool all asset prices set and amount zero case",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 0), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 0), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			poolWeights: []types.AssetWeight{
				{
					Asset:  "uusdc",
					Weight: sdk.ZeroDec(),
				},
				{
					Asset:  "uusdt",
					Weight: sdk.ZeroDec(),
				},
			},
			expError: false,
		},
		{
			desc: "oracle pool one asset price not set",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uUSDT", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
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
			weights, err := types.OraclePoolNormalizedWeights(suite.ctx, suite.app.OracleKeeper, tc.poolAssets)
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
		poolAssets      []*types.PoolAsset
		inCoins         sdk.Coins
		outCoins        sdk.Coins
		poolAssetsAfter []*types.PoolAsset
		expErr          bool
	}{
		{
			desc: "positive in and no out case",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			inCoins:  sdk.Coins{sdk.NewInt64Coin("uusdc", 1000_000)},
			outCoins: sdk.Coins{},
			poolAssetsAfter: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1001_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			expErr: false,
		},
		{
			desc: "no in and positive out case",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			inCoins:  sdk.Coins{},
			outCoins: sdk.Coins{sdk.NewInt64Coin("uusdc", 1000_000)},
			poolAssetsAfter: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 999_000_000), // 999 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			expErr: false,
		},
		{
			desc: "positive in and positive out case",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			inCoins:  sdk.Coins{sdk.NewInt64Coin("uusdt", 1000_000)},
			outCoins: sdk.Coins{sdk.NewInt64Coin("uusdc", 1000_000)},
			poolAssetsAfter: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 999_000_000), // 999 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1001_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			expErr: false,
		},
		{
			desc: "withdrawing more than pool size",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			inCoins:         sdk.Coins{sdk.NewInt64Coin("uusdt", 1000_000)},
			outCoins:        sdk.Coins{sdk.NewInt64Coin("uusdc", 1001_000_000)},
			poolAssetsAfter: []*types.PoolAsset{},
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
				PoolParams: &types.PoolParams{
					SwapFee:   sdk.ZeroDec(),
					UseOracle: true,
				},
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdk.ZeroInt(),
			}
			poolAssets, err := pool.NewPoolAssetsAfterSwap(tc.inCoins, tc.outCoins)
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
		poolAssets  []*types.PoolAsset
		expDistance sdk.Dec
	}{
		{
			desc: "zero balance for one asset",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 0), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			expDistance: sdk.NewDecWithPrec(5, 1),
		},
		{
			desc: "zero for all assets",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 0), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 0), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			expDistance: sdk.NewDecWithPrec(5, 1),
		},
		{
			desc: "all positive",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 500_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			expDistance: sdk.NewDecWithPrec(25, 2),
		},
		{
			desc: "zero distance",
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			expDistance: sdk.ZeroDec(),
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
				PoolParams: &types.PoolParams{
					SwapFee:   sdk.ZeroDec(),
					UseOracle: true,
				},
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdk.ZeroInt(),
			}
			distance := pool.WeightDistanceFromTarget(suite.ctx, suite.app.OracleKeeper, tc.poolAssets)
			suite.Require().Equal(distance, tc.expDistance)
		})
	}
}

func (suite *TestSuite) TestSwapOutAmtGivenIn() {
	for _, tc := range []struct {
		desc                string
		poolAssets          []*types.PoolAsset
		useOracle           bool
		slippageReduction   sdk.Dec
		thresholdWeightDiff sdk.Dec
		tokenIn             sdk.Coin
		outTokenDenom       string
		expRecoveryBonus    sdk.Dec
		expTokenOut         sdk.Coin
		expErr              bool
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
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle:           true,
			slippageReduction:   sdk.NewDecWithPrec(90, 2),              // 90%
			thresholdWeightDiff: sdk.NewDecWithPrec(20, 2),              // 20%
			tokenIn:             sdk.NewInt64Coin("uusdt", 100_000_000), // 100 USDC
			outTokenDenom:       "uusdc",
			expRecoveryBonus:    sdk.ZeroDec(),
			expTokenOut:         sdk.NewInt64Coin("uusdc", 94161125),
			expErr:              false,
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
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 500_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle:           true,
			slippageReduction:   sdk.NewDecWithPrec(90, 2),              // 90%
			thresholdWeightDiff: sdk.NewDecWithPrec(20, 2),              // 20%
			tokenIn:             sdk.NewInt64Coin("uusdt", 100_000_000), // 100 USDC
			outTokenDenom:       "uusdc",
			expRecoveryBonus:    sdk.ZeroDec(),
			expTokenOut:         sdk.NewInt64Coin("uusdc", 88723966),
			expErr:              false,
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
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 500_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle:           true,
			slippageReduction:   sdk.NewDecWithPrec(90, 2),              // 90%
			thresholdWeightDiff: sdk.NewDecWithPrec(20, 2),              // 20%
			tokenIn:             sdk.NewInt64Coin("uusdc", 100_000_000), // 100 USDC
			outTokenDenom:       "uusdt",
			expRecoveryBonus:    sdk.MustNewDecFromStr("0.052267002518891688"),
			expTokenOut:         sdk.NewInt64Coin("uusdt", 115_000_000),
			expErr:              false,
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
			poolAssets: []*types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 500_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1500_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle:           false,
			slippageReduction:   sdk.NewDecWithPrec(90, 2),              // 90%
			thresholdWeightDiff: sdk.NewDecWithPrec(20, 2),              // 20%
			tokenIn:             sdk.NewInt64Coin("uusdc", 100_000_000), // 100 USDC
			outTokenDenom:       "uusdt",
			expRecoveryBonus:    sdk.ZeroDec(),
			expTokenOut:         sdk.NewInt64Coin("uusdt", 250_000_000),
			expErr:              false,
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
				PoolParams: &types.PoolParams{
					SwapFee:                     sdk.ZeroDec(),
					UseOracle:                   tc.useOracle,
					SlippageReduction:           tc.slippageReduction,
					ThresholdWeightDifference:   tc.thresholdWeightDiff,
					WeightBreakingFeeMultiplier: sdk.OneDec(),
				},
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdk.ZeroInt(),
			}
			tokenOut, weightBonus, err := pool.SwapOutAmtGivenIn(suite.ctx, suite.app.OracleKeeper, sdk.Coins{tc.tokenIn}, tc.outTokenDenom)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tokenOut.String(), tc.expTokenOut.String())
				suite.Require().Equal(weightBonus.String(), tc.expRecoveryBonus.String())
			}
		})
	}
}
