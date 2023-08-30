package types_test

import (
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (suite *TestSuite) TestTVL() {
	for _, tc := range []struct {
		desc       string
		poolAssets []types.PoolAsset
		useOracle  bool
		expTVL     sdk.Dec
		expError   bool
	}{
		{
			desc: "oracle pool all asset prices set case",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDT
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle: true,
			expTVL:    sdk.NewDec(2000),
			expError:  false,
		},
		{
			desc: "oracle pool one asset price not set",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("ujuno", 1000_000_000), // 1000 JUNO
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdt", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle: true,
			expTVL:    sdk.NewDec(0),
			expError:  true,
		},
		{
			desc: "non-oracle pool not asset price set",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("ujuno", 1000_000_000), // 1000 JUNO
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("ukava", 1000_000_000), // 1000 KAVA
					Weight: sdk.NewInt(50),
				},
			},
			useOracle: false,
			expTVL:    sdk.NewDec(0),
			expError:  false,
		},
		{
			desc: "non-oracle pool one price set",
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin("ujuno", 1000_000_000), // 1000 JUNO
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewInt64Coin("uusdc", 1000_000_000), // 1000 USDC
					Weight: sdk.NewInt(50),
				},
			},
			useOracle: false,
			expTVL:    sdk.NewDec(2000),
			expError:  false,
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
					SwapFee:   sdk.ZeroDec(),
					UseOracle: tc.useOracle,
				},
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdk.ZeroInt(),
			}
			tvl, err := pool.TVL(suite.ctx, suite.app.OracleKeeper, suite.app.AccountedPoolKeeper)
			if tc.expError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tvl.String(), tc.expTVL.String())
			}
		})
	}
}
