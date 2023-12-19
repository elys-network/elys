package types_test

import (
	fmt "fmt"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestPoolTVL() {
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
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDT
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
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000), // 1000 USDC
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

func TestPool_GetPoolAssetAndIndex(t *testing.T) {
	poolAssets := []types.PoolAsset{
		{
			Token:  sdk.NewCoin("token1", sdk.NewInt(100)),
			Weight: sdk.NewInt(10),
		},
		{
			Token:  sdk.NewCoin("token2", sdk.NewInt(200)),
			Weight: sdk.NewInt(20),
		},
	}

	pool := types.Pool{
		PoolAssets: poolAssets,
	}

	// Test case 1: Existing PoolAsset
	index, poolAsset, err := pool.GetPoolAssetAndIndex("token1")
	require.NoError(t, err)
	require.Equal(t, 0, index)
	require.Equal(t, poolAssets[0], poolAsset)

	// Test case 1b: Existing PoolAsset
	index, poolAsset, err = pool.GetPoolAssetAndIndex("token2")
	require.NoError(t, err)
	require.Equal(t, 1, index)
	require.Equal(t, poolAssets[1], poolAsset)

	// Test case 2: Non-existing PoolAsset
	nonExistingDenom := "nonExistingToken"
	_, _, err = pool.GetPoolAssetAndIndex(nonExistingDenom)
	expectedErr := errorsmod.Wrapf(types.ErrDenomNotFoundInPool, fmt.Sprintf(types.FormatNoPoolAssetFoundErrFormat, nonExistingDenom))
	require.EqualError(t, err, expectedErr.Error())

	// Test case 3: Empty denom
	_, _, err = pool.GetPoolAssetAndIndex("")
	require.EqualError(t, err, "you tried to find the PoolAsset with empty denom")
}

func TestPool_GetAllPoolAssets(t *testing.T) {
	poolAssets := []types.PoolAsset{
		{
			Token:  sdk.NewCoin("token1", sdk.NewInt(100)),
			Weight: sdk.NewInt(10),
		},
		{
			Token:  sdk.NewCoin("token2", sdk.NewInt(200)),
			Weight: sdk.NewInt(20),
		},
	}

	pool := types.Pool{
		PoolAssets: poolAssets,
	}

	result := pool.GetAllPoolAssets()

	require.Equal(t, len(poolAssets), len(result))
	for i := 0; i < len(poolAssets); i++ {
		require.Equal(t, poolAssets[i], result[i])
	}
}
