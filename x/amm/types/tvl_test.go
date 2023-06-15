package types_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *TestSuite) TestTVL() {
	for _, tc := range []struct {
		desc       string
		poolAssets []*types.PoolAsset
		useOracle  bool
		expTVL     sdk.Dec
		expError   bool
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
			useOracle: true,
			expTVL:    sdk.NewDec(2000),
			expError:  false,
		},
		{
			desc: "oracle pool one asset price not set",
			poolAssets: []*types.PoolAsset{
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
			poolAssets: []*types.PoolAsset{
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
			poolAssets: []*types.PoolAsset{
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
			provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// prices set for USDT and USDC
			suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
				Denom:   "uusdc",
				Display: "USDC",
				Decimal: 6,
			})
			suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
				Denom:   "uusdt",
				Display: "USDT",
				Decimal: 6,
			})
			suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
				Asset:     "USDC",
				Price:     sdk.NewDec(1),
				Source:    "elys",
				Provider:  provider.String(),
				Timestamp: uint64(now.Unix()),
			})
			suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
				Asset:     "USDT",
				Price:     sdk.NewDec(1),
				Source:    "elys",
				Provider:  provider.String(),
				Timestamp: uint64(now.Unix()),
			})

			// execute function
			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: &types.PoolParams{
					SwapFee:   sdk.ZeroDec(),
					UseOracle: tc.useOracle,
				},
				TotalShares: sdk.Coin{},
				PoolAssets:  tc.poolAssets,
				TotalWeight: sdk.ZeroInt(),
			}
			tvl, err := pool.TVL(suite.ctx, suite.app.OracleKeeper)
			if tc.expError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tvl.String(), tc.expTVL.String())
			}
		})
	}
}
