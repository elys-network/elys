package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestApplyJoinPoolStateChange() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"join pool",
			func() {
				suite.ResetSuite()

				suite.SetupCoinPrices()
				suite.SetupAssetProfile()
			},
			func() {
				app := suite.app
				amm, bk := app.AmmKeeper, app.BankKeeper
				ctx := suite.ctx

				err := simapp.SetStakingParam(app, ctx)
				suite.Require().NoError(err)

				addr := suite.AddAccounts(1, nil)[0]

				poolAssets := []types.PoolAsset{
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000)),
					},
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)),
					},
				}

				poolParams := types.PoolParams{
					SwapFee:   sdkmath.LegacyMustNewDecFromStr("0.01"),
					UseOracle: true,
					FeeDenom:  ptypes.BaseCurrency,
				}

				msg := types.NewMsgCreatePool(
					addr.String(),
					poolParams,
					poolAssets,
				)

				// Create a ATOM+USDC pool
				poolId, err := amm.CreatePool(ctx, msg)
				suite.Require().NoError(err)
				suite.Require().Equal(poolId, uint64(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

				pool, found := amm.GetPool(ctx, poolId)
				suite.Require().True(found)

				lpTokenDenom := types.GetPoolShareDenom(poolId)
				lpTokenBalance := bk.GetBalance(ctx, addr, lpTokenDenom)
				suite.Require().True(lpTokenBalance.Amount.Equal(sdkmath.ZeroInt()))

				joinCoins := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)), sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000)))

				ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour))
				err = app.AmmKeeper.ApplyJoinPoolStateChange(ctx, pool, addr, pool.TotalShares.Amount, joinCoins, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{})
				suite.Require().NoError(err)
			},
		},
		{
			"join pool fails because join coins exceeds user balance",
			func() {
				suite.ResetSuite()

				suite.SetupCoinPrices()
				suite.SetupAssetProfile()
			},
			func() {
				app := suite.app
				amm, bk := app.AmmKeeper, app.BankKeeper
				ctx := suite.ctx

				err := simapp.SetStakingParam(app, ctx)
				suite.Require().NoError(err)

				addr := suite.AddAccounts(1, nil)[0]

				poolAssets := []types.PoolAsset{
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000)),
					},
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)),
					},
				}

				poolParams := types.PoolParams{
					SwapFee:   sdkmath.LegacyMustNewDecFromStr("0.01"),
					UseOracle: true,
					FeeDenom:  ptypes.BaseCurrency,
				}

				msg := types.NewMsgCreatePool(
					addr.String(),
					poolParams,
					poolAssets,
				)

				// Create a ATOM+USDC pool
				poolId, err := amm.CreatePool(ctx, msg)
				suite.Require().NoError(err)
				suite.Require().Equal(poolId, uint64(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

				pool, found := amm.GetPool(ctx, poolId)
				suite.Require().True(found)

				lpTokenDenom := types.GetPoolShareDenom(poolId)
				lpTokenBalance := bk.GetBalance(ctx, addr, lpTokenDenom)
				suite.Require().True(lpTokenBalance.Amount.Equal(sdkmath.ZeroInt()))

				joinCoins := sdk.NewCoins(
					sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount().Mul(sdkmath.NewInt(10))),
					sdk.NewCoin(ptypes.ATOM, suite.GetAccountIssueAmount().Mul(sdkmath.NewInt(10))),
				)

				ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour))
				err = app.AmmKeeper.ApplyJoinPoolStateChange(ctx, pool, addr, pool.TotalShares.Amount, joinCoins, sdkmath.LegacyZeroDec(), sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdk.Coins{})
				suite.Require().Error(err)
			},
		},
		{
			"join pool fails because numshares is negative",
			func() {
				suite.ResetSuite()

				suite.SetupCoinPrices()
				suite.SetupAssetProfile()
			},
			func() {
				app := suite.app
				amm, bk := app.AmmKeeper, app.BankKeeper
				ctx := suite.ctx

				err := simapp.SetStakingParam(app, ctx)
				suite.Require().NoError(err)

				addr := suite.AddAccounts(1, nil)[0]

				poolAssets := []types.PoolAsset{
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000)),
					},
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)),
					},
				}

				poolParams := types.PoolParams{
					SwapFee:   sdkmath.LegacyMustNewDecFromStr("0.01"),
					UseOracle: true,
					FeeDenom:  ptypes.BaseCurrency,
				}

				msg := types.NewMsgCreatePool(
					addr.String(),
					poolParams,
					poolAssets,
				)

				// Create a ATOM+USDC pool
				poolId, err := amm.CreatePool(ctx, msg)
				suite.Require().NoError(err)
				suite.Require().Equal(poolId, uint64(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

				pool, found := amm.GetPool(ctx, poolId)
				suite.Require().True(found)

				lpTokenDenom := types.GetPoolShareDenom(poolId)
				lpTokenBalance := bk.GetBalance(ctx, addr, lpTokenDenom)
				suite.Require().True(lpTokenBalance.Amount.Equal(sdkmath.ZeroInt()))

				joinCoins := sdk.NewCoins(
					sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)),
					sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000)),
				)

				ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour))

				// must panic
				suite.Require().Panics(func() {
					err = app.AmmKeeper.ApplyJoinPoolStateChange(ctx, pool, addr, sdkmath.NewInt(-1000), joinCoins, sdkmath.LegacyZeroDec(), sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdk.Coins{})
					suite.Require().Error(err)
				})
			},
		},
		{
			"join pool with positive weight balance bonus",
			func() {
				suite.ResetSuite()

				suite.SetupCoinPrices()
				suite.SetupAssetProfile()
			},
			func() {
				app := suite.app
				amm, bk := app.AmmKeeper, app.BankKeeper
				ctx := suite.ctx

				err := simapp.SetStakingParam(app, ctx)
				suite.Require().NoError(err)

				addr := suite.AddAccounts(1, nil)[0]

				poolAssets := []types.PoolAsset{
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000)),
					},
					{
						Weight: sdkmath.NewInt(50),
						Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)),
					},
				}

				poolParams := types.PoolParams{
					SwapFee:   sdkmath.LegacyMustNewDecFromStr("0.01"),
					UseOracle: true,
					FeeDenom:  ptypes.BaseCurrency,
				}

				msg := types.NewMsgCreatePool(
					addr.String(),
					poolParams,
					poolAssets,
				)

				// Create a ATOM+USDC pool
				poolId, err := amm.CreatePool(ctx, msg)
				suite.Require().NoError(err)
				suite.Require().Equal(poolId, uint64(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

				pool, found := amm.GetPool(ctx, poolId)
				suite.Require().True(found)

				lpTokenDenom := types.GetPoolShareDenom(poolId)
				lpTokenBalance := bk.GetBalance(ctx, addr, lpTokenDenom)
				suite.Require().True(lpTokenBalance.Amount.Equal(sdkmath.ZeroInt()))

				rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
				// send 1000ATOM and 1000USDC to treasury
				coins := sdk.NewCoins(
					sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(1000)),
					sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(1000)),
				)
				// mint coins
				err = bk.MintCoins(ctx, types.ModuleName, coins)
				suite.Require().NoError(err)
				// send coins to treasury
				err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, rebalanceTreasuryAddr, coins)
				suite.Require().NoError(err)

				joinCoins := sdk.NewCoins(
					sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)),
					sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000)),
				)

				ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour))
				err = app.AmmKeeper.ApplyJoinPoolStateChange(ctx, pool, addr, pool.TotalShares.Amount, joinCoins, sdkmath.LegacyNewDecWithPrec(10, 2), sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdk.Coins{})
				suite.Require().NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
