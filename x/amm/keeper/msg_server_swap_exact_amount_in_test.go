package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/v6/x/amm/keeper"
	"github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestMsgServerSwapExactAmountIn() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		swapFee           sdkmath.LegacyDec
		tokenIn           sdk.Coin
		tokenOutMin       sdkmath.Int
		tokenOut          sdk.Coin
		swapRoute         []types.SwapAmountInRoute
		expSenderBalance  sdk.Coins
		expPass           bool
	}{
		{
			desc:              "successful execution with positive swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyNewDecWithPrec(1, 2), // 1%
			tokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:       sdkmath.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin(ptypes.BaseCurrency, 9802),
			swapRoute: []types.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: ptypes.BaseCurrency,
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009802)},
			expPass:          true,
		},
		{
			desc:              "successful execution with zero swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:       sdkmath.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			swapRoute: []types.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: ptypes.BaseCurrency,
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009900)},
			expPass:          true,
		},
		{
			desc:              "not available pool as route",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:       sdkmath.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			swapRoute: []types.SwapAmountInRoute{
				{
					PoolId:        3,
					TokenOutDenom: ptypes.BaseCurrency,
				},
			},
			expSenderBalance: sdk.Coins{},
			expPass:          false,
		},
		{
			desc:              "multiple routes",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenIn:           sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			tokenOutMin:       sdkmath.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin("uusdt", 9802),
			swapRoute: []types.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: ptypes.Elys,
				},
				{
					PoolId:        2,
					TokenOutDenom: "uusdt",
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990000), sdk.NewInt64Coin("uusdt", 9802)},
			expPass:          true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolAddr := types.NewPoolAddress(uint64(1))
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolAddr2 := types.NewPoolAddress(uint64(2))
			treasuryAddr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolCoins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}
			pool2Coins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin("uusdt", 1000000)}

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolCoins)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolCoins)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, pool2Coins)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr2, pool2Coins)
			suite.Require().NoError(err)

			// execute function
			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
				Denom:     ptypes.Elys,
				Liquidity: sdkmath.NewInt(2000000),
			})
			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
				Denom:     ptypes.BaseCurrency,
				Liquidity: sdkmath.NewInt(1000000),
			})
			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
				Denom:     "uusdt",
				Liquidity: sdkmath.NewInt(1000000),
			})

			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: types.PoolParams{
					SwapFee:  tc.swapFee,
					FeeDenom: ptypes.BaseCurrency,
				},
				TotalShares: sdk.Coin{},
				PoolAssets: []types.PoolAsset{
					{
						Token:  poolCoins[0],
						Weight: sdkmath.NewInt(10),
					},
					{
						Token:  poolCoins[1],
						Weight: sdkmath.NewInt(10),
					},
				},
				TotalWeight: sdkmath.ZeroInt(),
			}
			pool2 := types.Pool{
				PoolId:            2,
				Address:           poolAddr2.String(),
				RebalanceTreasury: treasuryAddr2.String(),
				PoolParams: types.PoolParams{
					SwapFee:  tc.swapFee,
					FeeDenom: ptypes.BaseCurrency,
				},
				TotalShares: sdk.Coin{},
				PoolAssets: []types.PoolAsset{
					{
						Token:  pool2Coins[0],
						Weight: sdkmath.NewInt(10),
					},
					{
						Token:  pool2Coins[1],
						Weight: sdkmath.NewInt(10),
					},
				},
				TotalWeight: sdkmath.ZeroInt(),
			}
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)
			suite.app.AmmKeeper.SetPool(suite.ctx, pool2)
			suite.Require().True(suite.VerifyPoolAssetWithBalance(1))
			suite.Require().True(suite.VerifyPoolAssetWithBalance(2))

			msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
			resp, err := msgServer.SwapExactAmountIn(
				suite.ctx,
				&types.MsgSwapExactAmountIn{
					Sender:            sender.String(),
					Routes:            tc.swapRoute,
					TokenIn:           tc.tokenIn,
					TokenOutMinAmount: tc.tokenOutMin,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(resp.TokenOutAmount.String(), tc.tokenOut.Amount.String())
				suite.app.AmmKeeper.EndBlocker(suite.ctx)
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(2))

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())
			}
		})
	}
}

func (suite *AmmKeeperTestSuite) TestMsgServerSlippageDifferenceWhenSplit() {
	//suite.SetupTest()
	suite.SetupStableCoinPrices()

	senderInitBalance := sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000_000)}
	swapFee := sdkmath.LegacyZeroDec()
	tokenIn := sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000_000_000)
	tokenOutMin := sdkmath.ZeroInt()
	tokenOut := sdk.NewInt64Coin("uusdt", 98928368576)
	swapRoute := []types.SwapAmountInRoute{
		{
			PoolId:        1,
			TokenOutDenom: "uusdt",
		},
	}
	expSenderBalance := sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 900000000000), sdk.NewInt64Coin("uusdt", 98928368576)}
	expSenderBalanceSplitSwap := sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 900000000000), sdk.NewInt64Coin("uusdt", 99934261482)}

	// bootstrap accounts
	sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	poolAddr := types.NewPoolAddress(uint64(1))
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	poolCoins := sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000_000_000_000), sdk.NewInt64Coin("uusdt", 1000_000_000_000)}

	// bootstrap balances
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, senderInitBalance)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, senderInitBalance)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolCoins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolCoins)
	suite.Require().NoError(err)
	// execute function
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
		Denom:     ptypes.BaseCurrency,
		Liquidity: sdkmath.NewInt(1000000_000_000),
	})
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
		Denom:     "uusdt",
		Liquidity: sdkmath.NewInt(1000000_000_000),
	})

	pool := types.Pool{
		PoolId:            1,
		Address:           poolAddr.String(),
		RebalanceTreasury: treasuryAddr.String(),
		PoolParams: types.PoolParams{
			SwapFee:   swapFee,
			FeeDenom:  ptypes.BaseCurrency,
			UseOracle: true,
		},
		TotalShares: sdk.Coin{},
		PoolAssets: []types.PoolAsset{
			{
				Token:                  poolCoins[0],
				Weight:                 sdkmath.NewInt(10),
				ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
			},
			{
				Token:                  poolCoins[1],
				Weight:                 sdkmath.NewInt(10),
				ExternalLiquidityRatio: sdkmath.LegacyNewDec(10),
			},
		},
		TotalWeight: sdkmath.ZeroInt(),
	}
	suite.app.AmmKeeper.SetPool(suite.ctx, pool)
	suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

	cacheCtx, _ := suite.ctx.CacheContext()
	msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
	resp, err := msgServer.SwapExactAmountIn(
		cacheCtx,
		&types.MsgSwapExactAmountIn{
			Sender:            sender.String(),
			Routes:            swapRoute,
			TokenIn:           tokenIn,
			TokenOutMinAmount: tokenOutMin,
		})
	suite.Require().NoError(err)
	suite.Require().Equal(tokenOut.Amount.String(), resp.TokenOutAmount.String())
	suite.app.AmmKeeper.EndBlocker(cacheCtx)
	suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

	// check balance change on sender
	balances := suite.app.BankKeeper.GetAllBalances(cacheCtx, sender)
	suite.Require().Equal(balances.String(), expSenderBalance.String())

	// execute 100x swap with split
	cacheCtx, _ = suite.ctx.CacheContext()
	for i := 0; i < 100; i++ {
		resp, err = msgServer.SwapExactAmountIn(
			cacheCtx,
			&types.MsgSwapExactAmountIn{
				Sender:            sender.String(),
				Routes:            swapRoute,
				TokenIn:           sdk.Coin{Denom: tokenIn.Denom, Amount: tokenIn.Amount.Quo(sdkmath.NewInt(100))},
				TokenOutMinAmount: tokenOutMin,
			})
		suite.Require().NoError(err)
		suite.Require().True(suite.VerifyPoolAssetWithBalance(1))
		fmt.Printf("outAmount%d: %s\n", i, resp.TokenOutAmount.String())
	}
	suite.app.AmmKeeper.EndBlocker(cacheCtx)
	// check balance change on sender after splitting swap to 100
	balances = suite.app.BankKeeper.GetAllBalances(cacheCtx, sender)
	suite.Require().Equal(balances.String(), expSenderBalanceSplitSwap.String())
	suite.Require().True(suite.VerifyPoolAssetWithBalance(1))
}
