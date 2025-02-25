package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestApplyExitPoolStateChange() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"withdraw from commitment module",
			func() {
				suite.ResetSuite()

				suite.SetupStableCoinPrices()
				suite.SetupAssetProfile()

				app := suite.app
				amm, bk := app.AmmKeeper, app.BankKeeper
				ctx := suite.ctx

				err := simapp.SetStakingParam(app, ctx)
				suite.Require().NoError(err)
				// Generate 1 random account with 1000stake balanced
				addrs := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

				// Mint 100000USDC+100000USDT
				coins := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)), sdk.NewCoin("uusdt", sdkmath.NewInt(100000)))
				err = app.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
				suite.Require().NoError(err)
				err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addrs[0], coins)
				suite.Require().NoError(err)

				poolAssets := []types.PoolAsset{
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin("uusdt", sdkmath.NewInt(100000)),
					},
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)),
					},
				}

				swapFee, err := sdkmath.LegacyNewDecFromStr("0.01")
				suite.Require().NoError(err)

				poolParams := types.PoolParams{
					SwapFee:   swapFee,
					UseOracle: true,
					FeeDenom:  ptypes.BaseCurrency,
				}

				msg := types.NewMsgCreatePool(
					addrs[0].String(),
					poolParams,
					poolAssets,
				)

				// Create a USDT+USDC pool
				poolId, err := amm.CreatePool(ctx, msg)
				suite.Require().NoError(err)
				suite.Require().Equal(poolId, uint64(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

				pool, found := amm.GetPool(ctx, poolId)
				suite.Require().True(found)

				lpTokenDenom := types.GetPoolShareDenom(poolId)
				lpTokenBalance := bk.GetBalance(ctx, addrs[0], lpTokenDenom)
				suite.Require().True(lpTokenBalance.Amount.Equal(sdkmath.ZeroInt()))

				ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour))
				err = app.AmmKeeper.ApplyExitPoolStateChange(ctx, pool, addrs[0], pool.TotalShares.Amount, coins, false, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{})
				suite.Require().NoError(err)
			},
			func() {},
		},
		{
			"asset profile not found",
			func() {
				suite.ResetSuite()

				pool := types.Pool{
					PoolId: 1,
				}

				addr := suite.AddAccounts(1, nil)[0]

				coins := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)), sdk.NewCoin("uusdt", sdkmath.NewInt(100000)))

				err := suite.app.AmmKeeper.ApplyExitPoolStateChange(suite.ctx, pool, addr, pool.TotalShares.Amount, coins, false, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{})
				suite.Require().Error(err)
			},
			func() {},
		},
		{
			"withdraw from commitment module with invalid denom",
			func() {
				suite.ResetSuite()

				suite.SetupStableCoinPrices()
				suite.SetupAssetProfile()

				app := suite.app
				amm, bk := app.AmmKeeper, app.BankKeeper
				ctx := suite.ctx

				err := simapp.SetStakingParam(app, ctx)
				suite.Require().NoError(err)
				// Generate 1 random account with 1000stake balanced
				addr := suite.AddAccounts(1, nil)[0]

				// Mint 100000USDC+100000USDT
				coins := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)), sdk.NewCoin("uusdt", sdkmath.NewInt(100000)))
				err = app.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
				suite.Require().NoError(err)
				err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
				suite.Require().NoError(err)

				poolAssets := []types.PoolAsset{
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin("uusdt", sdkmath.NewInt(100000)),
					},
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)),
					},
				}

				swapFee, err := sdkmath.LegacyNewDecFromStr("0.01")
				suite.Require().NoError(err)

				poolParams := types.PoolParams{
					SwapFee:   swapFee,
					UseOracle: true,
					FeeDenom:  ptypes.BaseCurrency,
				}

				msg := types.NewMsgCreatePool(
					addr.String(),
					poolParams,
					poolAssets,
				)

				// Create a USDT+USDC pool
				poolId, err := amm.CreatePool(ctx, msg)
				suite.Require().NoError(err)
				suite.Require().Equal(poolId, uint64(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

				pool, found := amm.GetPool(ctx, poolId)
				suite.Require().True(found)

				lpTokenDenom := types.GetPoolShareDenom(poolId)
				lpTokenBalance := bk.GetBalance(ctx, addr, lpTokenDenom)
				suite.Require().True(lpTokenBalance.Amount.Equal(sdkmath.ZeroInt()))

				coins = sdk.NewCoins(sdk.NewCoin("invalid_denom", sdkmath.NewInt(100000000)))

				ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour))
				err = app.AmmKeeper.ApplyExitPoolStateChange(ctx, pool, addr, pool.TotalShares.Amount, coins, false, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{})
				suite.Require().Error(err)
			},
			func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
