package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestMsgServerSwapExactAmountIn() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		swapFee           sdk.Dec
		tokenIn           sdk.Coin
		tokenOutMin       sdk.Int
		tokenOut          sdk.Coin
		swapRoute         []types.SwapAmountInRoute
		expSenderBalance  sdk.Coins
		expPass           bool
	}{
		{
			desc:              "successful execution with positive swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapFee:           sdk.NewDecWithPrec(1, 2), // 1%
			tokenIn:           sdk.NewInt64Coin("uelys", 10000),
			tokenOutMin:       sdk.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin("uusdc", 9704),
			swapRoute: []types.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: "uusdc",
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 990000), sdk.NewInt64Coin("uusdc", 1009704)},
			expPass:          true,
		},
		{
			desc:              "successful execution with zero swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapFee:           sdk.ZeroDec(),
			tokenIn:           sdk.NewInt64Coin("uelys", 10000),
			tokenOutMin:       sdk.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin("uusdc", 9900),
			swapRoute: []types.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: "uusdc",
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 990000), sdk.NewInt64Coin("uusdc", 1009900)},
			expPass:          true,
		},
		{
			desc:              "not available pool as route",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapFee:           sdk.ZeroDec(),
			tokenIn:           sdk.NewInt64Coin("uelys", 10000),
			tokenOutMin:       sdk.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin("uusdc", 9900),
			swapRoute: []types.SwapAmountInRoute{
				{
					PoolId:        3,
					TokenOutDenom: "uusdc",
				},
			},
			expSenderBalance: sdk.Coins{},
			expPass:          false,
		},
		{
			desc:              "multiple routes",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapFee:           sdk.ZeroDec(),
			tokenIn:           sdk.NewInt64Coin("uusdc", 10000),
			tokenOutMin:       sdk.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin("uusdt", 9802),
			swapRoute: []types.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: "uelys",
				},
				{
					PoolId:        2,
					TokenOutDenom: "uusdt",
				},
			},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 990000), sdk.NewInt64Coin("uusdt", 9802)},
			expPass:          true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolAddr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolCoins := sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)}
			pool2Coins := sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdt", 1000000)}

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
				Denom:     "uelys",
				Liquidity: sdk.NewInt(2000000),
			})
			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
				Denom:     "uusdc",
				Liquidity: sdk.NewInt(1000000),
			})
			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
				Denom:     "uusdt",
				Liquidity: sdk.NewInt(1000000),
			})

			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: types.PoolParams{
					SwapFee: tc.swapFee,
				},
				TotalShares: sdk.Coin{},
				PoolAssets: []types.PoolAsset{
					{
						Token:  poolCoins[0],
						Weight: sdk.NewInt(10),
					},
					{
						Token:  poolCoins[1],
						Weight: sdk.NewInt(10),
					},
				},
				TotalWeight: sdk.ZeroInt(),
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
						Weight: sdk.NewInt(10),
					},
					{
						Token:  pool2Coins[1],
						Weight: sdk.NewInt(10),
					},
				},
				TotalWeight: sdk.ZeroInt(),
			}
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)
			suite.app.AmmKeeper.SetPool(suite.ctx, pool2)

			msgServer := keeper.NewMsgServerImpl(suite.app.AmmKeeper)
			resp, err := msgServer.SwapExactAmountIn(
				sdk.WrapSDKContext(suite.ctx),
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

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())
			}
		})
	}
}
