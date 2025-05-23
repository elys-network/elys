package keeper_test

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/v5/x/amm/keeper"
	"github.com/elys-network/elys/v5/x/amm/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestExecuteSwapRequests() {
	sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		swapMsgs          []sdk.Msg
		expSwapOrder      []uint64
	}{
		{
			desc:              "single swap request",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.BaseCurrency,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{0},
		},
		{
			desc:              "two requests with opposite direction",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.BaseCurrency,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.Elys,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.BaseCurrency, 8000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{1, 0},
		},
		{
			desc:              "three requests",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.BaseCurrency,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.Elys, 11000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.Elys,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.BaseCurrency, 8000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.Elys,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.BaseCurrency, 1000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{1, 2, 0},
		},
		{
			desc:              "three requests combining different swap msg types",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.BaseCurrency,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.Elys, 11000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
				&types.MsgSwapExactAmountOut{
					Sender: sender.String(),
					Routes: []types.SwapAmountOutRoute{
						{
							PoolId:       1,
							TokenInDenom: ptypes.BaseCurrency,
						},
					},
					TokenOut:         sdk.NewInt64Coin(ptypes.Elys, 8000),
					TokenInMaxAmount: sdkmath.NewInt(1000000),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.Elys,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.BaseCurrency, 1000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{2, 1, 0},
		},
		{
			desc:              "three requests combining different swap msg types",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapMsgs: []sdk.Msg{
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.BaseCurrency,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.Elys, 11000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
				&types.MsgSwapExactAmountOut{
					Sender: sender.String(),
					Routes: []types.SwapAmountOutRoute{
						{
							PoolId:       1,
							TokenInDenom: ptypes.BaseCurrency,
						},
					},
					TokenOut:         sdk.NewInt64Coin(ptypes.Elys, 8000),
					TokenInMaxAmount: sdkmath.NewInt(1000000),
				},
				&types.MsgSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.Elys,
						},
						{
							PoolId:        2,
							TokenOutDenom: "uusdt",
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.BaseCurrency, 1000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				},
			},
			expSwapOrder: []uint64{2, 1, 0},
		},
	} {
		suite.Run(tc.desc, func() {
			suite.ResetSuite()

			// bootstrap accounts
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
				Address:           types.NewPoolAddress(uint64(1)).String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: types.PoolParams{
					SwapFee:  sdkmath.LegacyZeroDec(),
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
				Address:           types.NewPoolAddress(uint64(2)).String(),
				RebalanceTreasury: treasuryAddr2.String(),
				PoolParams: types.PoolParams{
					SwapFee:  sdkmath.LegacyZeroDec(),
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
			for _, msg := range tc.swapMsgs {
				switch msg := msg.(type) {
				case *types.MsgSwapExactAmountIn:
					_, err := msgServer.SwapExactAmountIn(
						suite.ctx,
						msg,
					)
					suite.Require().NoError(err)
				case *types.MsgSwapExactAmountOut:
					_, err := msgServer.SwapExactAmountOut(
						suite.ctx,
						msg,
					)
					suite.Require().NoError(err)
				}
			}

			requests := suite.app.AmmKeeper.ExecuteSwapRequests(suite.ctx)
			for i, order := range tc.expSwapOrder {
				suite.Require().Equal(tc.swapMsgs[order], requests[i], fmt.Sprintf("%dth item not match", i))
			}
			suite.Require().True(suite.VerifyPoolAssetWithBalance(1))
			suite.Require().True(suite.VerifyPoolAssetWithBalance(2))
		})
	}
}

func (suite *AmmKeeperTestSuite) TestClearOutdatedSlippageTrack() {
	now := time.Now()
	tracks := []types.OraclePoolSlippageTrack{
		{
			PoolId:    1,
			Timestamp: uint64(now.Unix() - 86400*8),
			Tracked:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000)},
		},
		{
			PoolId:    1,
			Timestamp: uint64(now.Unix() - 86400),
			Tracked:   sdk.Coins{sdk.NewInt64Coin("uelys", 2000)},
		},
		{
			PoolId:    2,
			Timestamp: uint64(now.Unix()),
			Tracked:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000)},
		},
	}
	for _, track := range tracks {
		suite.app.AmmKeeper.SetSlippageTrack(suite.ctx, track)
	}

	suite.ctx = suite.ctx.WithBlockTime(now)
	suite.app.AmmKeeper.ClearOutdatedSlippageTrack(suite.ctx)
	tracksStored := suite.app.AmmKeeper.AllSlippageTracks(suite.ctx)
	suite.Require().Len(tracksStored, 2)
}

func (suite *AmmKeeperTestSuite) TestAbci() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"first pool id with swap amount out",
			func() {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)[0]

				msg := &types.MsgSwapExactAmountOut{
					Sender: addr.String(),
					Routes: []types.SwapAmountOutRoute{
						{
							PoolId:       1,
							TokenInDenom: ptypes.BaseCurrency,
						},
					},
					TokenOut:         sdk.NewInt64Coin(ptypes.Elys, 8000),
					TokenInMaxAmount: sdkmath.NewInt(1000000),
				}

				poolId := suite.app.AmmKeeper.FirstPoolId(msg)
				suite.Require().Equal(uint64(1), poolId)
			},
			func() {},
		},
		{
			"first pool id with a msg that is neither swap amount in nor swap amount out",
			func() {
				suite.ResetSuite()

				msg := sdk.Msg(nil)

				poolId := suite.app.AmmKeeper.FirstPoolId(msg)
				suite.Require().Equal(uint64(0), poolId)
			},
			func() {},
		},
		{
			"apply swap request with invalid address in swap exact amount in msg",
			func() {
				suite.ResetSuite()

				msg := &types.MsgSwapExactAmountIn{
					Sender: "invalid",
				}

				err := suite.app.AmmKeeper.ApplySwapRequest(suite.ctx, msg)
				suite.Require().Error(err)
			},
			func() {},
		},
		{
			"apply swap request with invalid denom in swap exact amount in msg",
			func() {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)[0]

				msg := &types.MsgSwapExactAmountIn{
					Sender: addr.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: "invalid",
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				}

				err := suite.app.AmmKeeper.ApplySwapRequest(suite.ctx, msg)
				suite.Require().Error(err)
			},
			func() {},
		},
		{
			"apply swap request with invalid address in swap exact amount in msg",
			func() {
				suite.ResetSuite()

				msg := &types.MsgSwapExactAmountOut{
					Sender: "invalid",
				}

				err := suite.app.AmmKeeper.ApplySwapRequest(suite.ctx, msg)
				suite.Require().Error(err)
			},
			func() {},
		},
		{
			"apply swap request with invalid denom in swap exact amount in msg",
			func() {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)[0]

				msg := &types.MsgSwapExactAmountOut{
					Sender: addr.String(),
					Routes: []types.SwapAmountOutRoute{
						{
							PoolId:       1,
							TokenInDenom: "invalid",
						},
					},
					TokenOut:         sdk.NewInt64Coin(ptypes.Elys, 10000),
					TokenInMaxAmount: sdkmath.ZeroInt(),
				}

				err := suite.app.AmmKeeper.ApplySwapRequest(suite.ctx, msg)
				suite.Require().Error(err)
			},
			func() {},
		},
		{
			"apply swap request with invalid swap msg type",
			func() {
				suite.ResetSuite()

				msg := sdk.Msg(nil)

				err := suite.app.AmmKeeper.ApplySwapRequest(suite.ctx, msg)
				suite.Require().Error(err)
			},
			func() {},
		},
		{
			"get stacked slippage when get pool returns not found",
			func() {
				suite.ResetSuite()

				poolId := uint64(2)
				ratio := suite.app.AmmKeeper.GetStackedSlippage(suite.ctx, poolId)
				suite.Require().Equal(sdkmath.LegacyZeroDec(), ratio.Dec())
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
