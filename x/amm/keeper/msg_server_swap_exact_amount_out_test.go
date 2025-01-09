package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestMsgServerSwapExactAmountOut() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		swapFee           sdkmath.LegacyDec
		tokenIn           sdk.Coin
		tokenInMax        sdkmath.Int
		tokenOut          sdk.Coin
		swapRoutes        []types.SwapAmountOutRoute
		expSenderBalance  sdk.Coins
		expPass           bool
	}{
		{
			desc:              "successful execution with positive swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyNewDecWithPrec(1, 2), // 1%
			tokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10203),
			tokenInMax:        sdkmath.NewInt(10000000),
			tokenOut:          sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			swapRoutes: []types.SwapAmountOutRoute{
				{
					PoolId:       1,
					TokenInDenom: ptypes.Elys,
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 989797), sdk.NewInt64Coin(ptypes.BaseCurrency, 1010000)},
			expPass:          true,
		},
		{
			desc:              "successful execution with zero swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10101),
			tokenInMax:        sdkmath.NewInt(10000000),
			tokenOut:          sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			swapRoutes: []types.SwapAmountOutRoute{
				{
					PoolId:       1,
					TokenInDenom: ptypes.Elys,
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 989899), sdk.NewInt64Coin(ptypes.BaseCurrency, 1010000)},
			expPass:          true,
		},
		{
			desc:              "not available pool as route",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10102),
			tokenInMax:        sdkmath.NewInt(10000000),
			tokenOut:          sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			swapRoutes: []types.SwapAmountOutRoute{
				{
					PoolId:       3,
					TokenInDenom: ptypes.Elys,
				},
			},
			expSenderBalance: sdk.Coins{},
			expPass:          false,
		},
		{
			desc:              "multiple routes",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenIn:           sdk.NewInt64Coin(ptypes.BaseCurrency, 10204),
			tokenInMax:        sdkmath.NewInt(10000000),
			tokenOut:          sdk.NewInt64Coin("uusdt", 10000),
			swapRoutes: []types.SwapAmountOutRoute{
				{
					PoolId:       1,
					TokenInDenom: ptypes.BaseCurrency,
				},
				{
					PoolId:       2,
					TokenInDenom: ptypes.Elys,
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 989796), sdk.NewInt64Coin("uusdt", 10000)},
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
					SwapFee: tc.swapFee,
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
			resp, err := msgServer.SwapExactAmountOut(
				suite.ctx,
				&types.MsgSwapExactAmountOut{
					Sender:           sender.String(),
					Routes:           tc.swapRoutes,
					TokenOut:         tc.tokenOut,
					TokenInMaxAmount: tc.tokenInMax,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(resp.TokenInAmount.String(), tc.tokenIn.Amount.String())
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
