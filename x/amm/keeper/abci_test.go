package keeper_test

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
)

func (suite *KeeperTestSuite) TestExecuteSwapRequests() {
	sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		swapMsgs          []sdk.Msg
		expSwapOrder      []uint64
	}{
		{
			desc:              "single swap request",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uusdc",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uelys", 10000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{0},
		},
		{
			desc:              "two requests with opposite direction",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uusdc",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uelys", 10000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uelys",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uusdc", 8000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{1, 0},
		},
		{
			desc:              "three requests",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uusdc",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uelys", 11000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uelys",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uusdc", 8000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uelys",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uusdc", 1000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{1, 0, 2},
		},
		{
			desc:              "three requests combining different swap msg types",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uusdc",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uelys", 11000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
				&types.MsgSwapExactAmountOut{
					Sender: sender.String(),
					Routes: []types.SwapAmountOutRoute{
						{
							PoolId:       1,
							TokenInDenom: "uusdc",
						},
					},
					TokenOut:         sdk.NewInt64Coin("uelys", 8000),
					TokenInMaxAmount: sdk.NewInt(1000000),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uelys",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uusdc", 1000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{2, 1, 0},
		},
		{
			desc:              "three requests combining different swap msg types",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uelys", 1000000), sdk.NewInt64Coin("uusdc", 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uusdc",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uelys", 11000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
				&types.MsgSwapExactAmountOut{
					Sender: sender.String(),
					Routes: []types.SwapAmountOutRoute{
						{
							PoolId:       1,
							TokenInDenom: "uusdc",
						},
					},
					TokenOut:         sdk.NewInt64Coin("uelys", 8000),
					TokenInMaxAmount: sdk.NewInt(1000000),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "uelys",
						},
						{
							PoolId:        2,
							TokenOutDenom: "uusdt",
						},
					},
					TokenIn:           sdk.NewInt64Coin("uusdc", 1000),
					TokenOutMinAmount: sdk.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{2, 1, 0},
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
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
					SwapFee:  sdk.ZeroDec(),
					FeeDenom: "uusdc",
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
					SwapFee:  sdk.ZeroDec(),
					FeeDenom: "uusdc",
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
			for _, msg := range tc.swapMsgs {
				switch msg := msg.(type) {
				case *types.MsgSwapExactAmountIn:
					_, err := msgServer.SwapExactAmountIn(
						sdk.WrapSDKContext(suite.ctx),
						msg,
					)
					suite.Require().NoError(err)
				case *types.MsgSwapExactAmountOut:
					_, err := msgServer.SwapExactAmountOut(
						sdk.WrapSDKContext(suite.ctx),
						msg,
					)
					suite.Require().NoError(err)
				}
			}

			requests := suite.app.AmmKeeper.ExecuteSwapRequests(suite.ctx)
			for i, order := range tc.expSwapOrder {
				suite.Require().Equal(tc.swapMsgs[order], requests[i], fmt.Sprintf("%dth item not match", i))
			}
		})
	}
}
