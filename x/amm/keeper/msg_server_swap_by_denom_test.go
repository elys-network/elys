package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/v5/x/amm/keeper"
	"github.com/elys-network/elys/v5/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/v5/x/oracle/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestMsgServerSwapByDenom() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		swapFee           sdkmath.LegacyDec
		tokenIn           sdk.Coin
		tokenOutMin       sdkmath.Int
		tokenOut          sdk.Coin
		expSenderBalance  sdk.Coins
		expPass           bool
		errMsg            string
	}{
		{
			desc:              "successful execution with positive swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyNewDecWithPrec(1, 2), // 1%
			tokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:       sdkmath.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin(ptypes.BaseCurrency, 9802),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009802)},
			expPass:           true,
		},
		{
			desc:              "successful execution with zero swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:       sdkmath.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009900)},
			expPass:           true,
		},
		{
			desc:              "multiple routes",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenIn:           sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			tokenOutMin:       sdkmath.ZeroInt(),
			tokenOut:          sdk.NewInt64Coin("uusdt", 9802),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990000), sdk.NewInt64Coin("uusdt", 9802)},
			expPass:           false,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// set asset profile
			suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
				BaseDenom: ptypes.Elys,
				Denom:     ptypes.Elys,
				Decimals:  6,
			})

			suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
				BaseDenom: ptypes.BaseCurrency,
				Denom:     ptypes.BaseCurrency,
				Decimals:  6,
			})

			// Set up oracle asset info
			suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
				Denom:   ptypes.Elys,
				Decimal: 6,
			})
			suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
				Denom:   ptypes.BaseCurrency,
				Decimal: 6,
			})
			suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
				Denom:   "uusda",
				Decimal: 6,
			})

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
			resp, err := msgServer.SwapByDenom(
				suite.ctx,
				&types.MsgSwapByDenom{
					Sender:    sender.String(),
					Amount:    tc.tokenIn,
					MinAmount: sdk.NewCoin(tc.tokenOut.Denom, tc.tokenOutMin),
					MaxAmount: sdk.Coin{},
					DenomIn:   tc.tokenIn.Denom,
					DenomOut:  tc.tokenOut.Denom,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(resp.Amount.String(), tc.tokenOut.String())
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

func (suite *AmmKeeperTestSuite) TestMsgServerSwapByDenomWithOutRoute() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		swapFee           sdkmath.LegacyDec
		tokenOut          sdk.Coin
		tokenOutMax       sdkmath.Int
		tokenDenomIn      string
		expTokenIn        sdk.Coin
		expSenderBalance  sdk.Coins
		expPass           bool
		errMsg            string
	}{
		{
			desc:              "successful execution with positive swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyNewDecWithPrec(1, 2), // 1%
			tokenOut:          sdk.NewInt64Coin(ptypes.Elys, 10000),
			expTokenIn:        sdk.NewInt64Coin(ptypes.BaseCurrency, 10204),
			tokenOutMax:       sdkmath.NewInt(1000000),
			tokenDenomIn:      ptypes.BaseCurrency,
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 989796)},
			expPass:           true,
		},
		{
			desc:              "successful execution with zero swap fee",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFee:           sdkmath.LegacyZeroDec(),
			tokenOut:          sdk.NewInt64Coin(ptypes.Elys, 10000),
			expTokenIn:        sdk.NewInt64Coin(ptypes.BaseCurrency, 10102),
			tokenOutMax:       sdkmath.NewInt(1000000),
			tokenDenomIn:      ptypes.BaseCurrency,
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 989898)},
			expPass:           true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// set asset profile
			suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
				BaseDenom: ptypes.Elys,
				Denom:     ptypes.Elys,
				Decimals:  6,
			})

			suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
				BaseDenom: ptypes.BaseCurrency,
				Denom:     ptypes.BaseCurrency,
				Decimals:  6,
			})

			// Set up oracle asset info
			suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
				Denom:   ptypes.Elys,
				Decimal: 6,
			})
			suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
				Denom:   ptypes.BaseCurrency,
				Decimal: 6,
			})
			suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
				Denom:   "uusda",
				Decimal: 6,
			})

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
			resp, err := msgServer.SwapByDenom(
				suite.ctx,
				&types.MsgSwapByDenom{
					Sender:    sender.String(),
					Amount:    tc.tokenOut,
					MinAmount: sdk.Coin{},
					MaxAmount: sdk.NewCoin(tc.tokenOut.Denom, tc.tokenOutMax),
					DenomIn:   tc.tokenDenomIn,
					DenomOut:  tc.tokenOut.Denom,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expTokenIn.String(), resp.Amount.String())
				suite.app.AmmKeeper.EndBlocker(suite.ctx)
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(2))

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(tc.expSenderBalance.String(), balances.String())
			}
		})
	}
}
