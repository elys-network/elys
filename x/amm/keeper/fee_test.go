package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestPortionCoins(t *testing.T) {
	coins := sdk.Coins{sdk.NewInt64Coin("ueden", 1000), sdk.NewInt64Coin("uelys", 10000)}
	portion := keeper.PortionCoins(coins, sdk.ZeroDec())
	require.Equal(t, portion, sdk.Coins{})

	portion = keeper.PortionCoins(coins, sdk.NewDecWithPrec(1, 1))
	require.Equal(t, portion, sdk.Coins{sdk.NewInt64Coin("ueden", 100), sdk.NewInt64Coin("uelys", 1000)})

	portion = keeper.PortionCoins(coins, sdk.NewDec(1))
	require.Equal(t, portion, coins)
}

func (suite *KeeperTestSuite) TestOnCollectFee() {
	for _, tc := range []struct {
		desc              string
		fee               sdk.Coins
		poolInitBalance   sdk.Coins
		expRevenueBalance sdk.Coins
		expPass           bool
	}{
		{
			desc:              "multiple fees collected",
			fee:               sdk.Coins{sdk.NewInt64Coin("uelys", 1000), sdk.NewInt64Coin("uusdc", 1000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			expRevenueBalance: sdk.Coins{sdk.NewInt64Coin("uusdc", 1799)},
			expPass:           true,
		},
		{
			desc:              "zero fees collected",
			fee:               sdk.Coins{},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			expRevenueBalance: sdk.Coins{},
			expPass:           true,
		},
		{
			desc:              "usdc fee collected",
			fee:               sdk.Coins{sdk.NewInt64Coin("uusdc", 1000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			expRevenueBalance: sdk.Coins{sdk.NewInt64Coin("uusdc", 900)},
			expPass:           true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
			poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			revenueAddr := types.NewPoolRevenueAddress(1)

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.fee)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, treasuryAddr, tc.fee)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.poolInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, tc.poolInitBalance)
			suite.Require().NoError(err)

			// execute function
			for _, coin := range tc.poolInitBalance {
				suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
					Denom:     coin.Denom,
					Liquidity: coin.Amount,
				})
			}
			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: types.PoolParams{
					SwapFee:                     sdk.ZeroDec(),
					ExitFee:                     sdk.ZeroDec(),
					UseOracle:                   false,
					WeightBreakingFeeMultiplier: sdk.ZeroDec(),
					ExternalLiquidityRatio:      sdk.NewDec(1),
					LpFeePortion:                sdk.ZeroDec(),
					StakingFeePortion:           sdk.ZeroDec(),
					WeightRecoveryFeePortion:    sdk.ZeroDec(),
					ThresholdWeightDifference:   sdk.ZeroDec(),
					FeeDenom:                    "uusdc",
				},
				TotalShares: sdk.Coin{},
				PoolAssets: []types.PoolAsset{
					{
						Token:  tc.poolInitBalance[0],
						Weight: sdk.NewInt(10),
					},
					{
						Token:  tc.poolInitBalance[1],
						Weight: sdk.NewInt(10),
					},
				},
				TotalWeight: sdk.ZeroInt(),
			}
			err = suite.app.AmmKeeper.OnCollectFee(suite.ctx, pool, tc.fee)
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// check pool balance increase/decrease
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, revenueAddr)
				suite.Require().Equal(balances.String(), tc.expRevenueBalance.String())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSwapFeesToRevenueToken() {
	for _, tc := range []struct {
		desc              string
		fee               sdk.Coins
		poolInitBalance   sdk.Coins
		expRevenueBalance sdk.Coins
		expPass           bool
	}{
		{
			desc:              "multiple fees collected",
			fee:               sdk.Coins{sdk.NewInt64Coin("uelys", 1000), sdk.NewInt64Coin("uusdc", 1000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			expRevenueBalance: sdk.Coins{sdk.NewInt64Coin("uusdc", 1999)},
			expPass:           true,
		},
		{
			desc:              "zero fees collected",
			fee:               sdk.Coins{},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			expRevenueBalance: sdk.Coins{},
			expPass:           true,
		},
		{
			desc:              "usdc fee collected",
			fee:               sdk.Coins{sdk.NewInt64Coin("uusdc", 1000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			expRevenueBalance: sdk.Coins{sdk.NewInt64Coin("uusdc", 1000)},
			expPass:           true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
			poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			revenueAddr := types.NewPoolRevenueAddress(1)

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.fee)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, revenueAddr, tc.fee)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.poolInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, tc.poolInitBalance)
			suite.Require().NoError(err)

			// execute function
			for _, coin := range tc.poolInitBalance {
				suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
					Denom:     coin.Denom,
					Liquidity: coin.Amount,
				})
			}
			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: types.PoolParams{
					SwapFee:                     sdk.ZeroDec(),
					ExitFee:                     sdk.ZeroDec(),
					UseOracle:                   false,
					WeightBreakingFeeMultiplier: sdk.ZeroDec(),
					ExternalLiquidityRatio:      sdk.NewDec(1),
					LpFeePortion:                sdk.ZeroDec(),
					StakingFeePortion:           sdk.ZeroDec(),
					WeightRecoveryFeePortion:    sdk.ZeroDec(),
					ThresholdWeightDifference:   sdk.ZeroDec(),
					FeeDenom:                    "uusdc",
				},
				TotalShares: sdk.Coin{},
				PoolAssets: []types.PoolAsset{
					{
						Token:  tc.poolInitBalance[0],
						Weight: sdk.NewInt(10),
					},
					{
						Token:  tc.poolInitBalance[1],
						Weight: sdk.NewInt(10),
					},
				},
				TotalWeight: sdk.ZeroInt(),
			}
			err = suite.app.AmmKeeper.SwapFeesToRevenueToken(suite.ctx, pool, tc.fee)
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// check pool balance increase/decrease
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, revenueAddr)
				suite.Require().Equal(balances.String(), tc.expRevenueBalance.String())
			}
		})
	}

	// // No fee management required when doing swap from fees to revenue token
	// func (k Keeper) SwapFeesToRevenueToken(ctx sdk.Context, pool types.Pool, fee sdk.Coins) error {
	// 	poolRevenueAddress := types.NewPoolRevenueAddress(pool.PoolId)
	// 	for _, tokenIn := range fee {
	// 		// skip for fee denom
	// 		if tokenIn.Denom == pool.PoolParams.FeeDenom {
	// 			continue
	// 		}
	// 		// Executes the swap in the pool and stores the output. Updates pool assets but
	// 		// does not actually transfer any tokens to or from the pool.
	// 		tokenOutCoin, _, err := pool.SwapOutAmtGivenIn(ctx, k.oracleKeeper, sdk.Coins{tokenIn}, pool.PoolParams.FeeDenom, sdk.ZeroDec())
	// 		if err != nil {
	// 			return err
	// 		}

	// 		tokenOutAmount := tokenOutCoin.Amount

	// 		if !tokenOutAmount.IsPositive() {
	// 			return sdkerrors.Wrapf(types.ErrInvalidMathApprox, "token amount must be positive")
	// 		}

	// 		// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
	// 		// Also emits a swap event and updates related liquidity metrics.
	// 		err, _ = k.UpdatePoolForSwap(ctx, pool, poolRevenueAddress, tokenIn, tokenOutCoin, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// }
}
