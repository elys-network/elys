package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestMsgServerJoinPool() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		poolInitBalance   sdk.Coins
		poolParams        types.PoolParams
		shareOutAmount    sdk.Int
		expSenderBalance  sdk.Coins
		expTotalLiquidity sdk.Coins
		expTokenIn        sdk.Coins
		expPass           bool
	}{
		{
			desc:              "successful non-oracle join pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uusdc", 100000), sdk.NewInt64Coin("uusdt", 100000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusdc", 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:                     sdk.ZeroDec(),
				ExitFee:                     sdk.ZeroDec(),
				UseOracle:                   false,
				WeightBreakingFeeMultiplier: sdk.ZeroDec(),
				SlippageReduction:           sdk.ZeroDec(),
				LpFeePortion:                sdk.ZeroDec(),
				StakingFeePortion:           sdk.ZeroDec(),
				WeightRecoveryFeePortion:    sdk.ZeroDec(),
				ThresholdWeightDifference:   sdk.ZeroDec(),
				FeeDenom:                    "uusdc",
			},
			shareOutAmount:   types.OneShare.Quo(sdk.NewInt(5)),
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{sdk.NewInt64Coin("uusdc", 100000), sdk.NewInt64Coin("uusdt", 100000)},
			expPass:          true,
		},
		{
			desc:              "not enough balance to join pool - non-oracle pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uusdc", 1000000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusdc", 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:                     sdk.ZeroDec(),
				ExitFee:                     sdk.ZeroDec(),
				UseOracle:                   false,
				WeightBreakingFeeMultiplier: sdk.ZeroDec(),
				SlippageReduction:           sdk.ZeroDec(),
				LpFeePortion:                sdk.ZeroDec(),
				StakingFeePortion:           sdk.ZeroDec(),
				WeightRecoveryFeePortion:    sdk.ZeroDec(),
				ThresholdWeightDifference:   sdk.ZeroDec(),
				FeeDenom:                    "uusdc",
			},
			shareOutAmount:   types.OneShare.Quo(sdk.NewInt(5)),
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{},
			expPass:          false,
		},
		{
			desc:              "oracle pool join - breaking weight on balanced pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uusdt", 1000000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusdc", 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:                     sdk.ZeroDec(),
				ExitFee:                     sdk.ZeroDec(),
				UseOracle:                   true,
				WeightBreakingFeeMultiplier: sdk.NewDecWithPrec(1, 0), // 1.00
				SlippageReduction:           sdk.ZeroDec(),
				LpFeePortion:                sdk.ZeroDec(),
				StakingFeePortion:           sdk.ZeroDec(),
				WeightRecoveryFeePortion:    sdk.ZeroDec(),
				ThresholdWeightDifference:   sdk.NewDecWithPrec(2, 1), // 20%
				FeeDenom:                    "uusdc",
			},
			shareOutAmount:   sdk.NewInt(833333333333333333), // weight breaking fee
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{sdk.NewInt64Coin("uusdt", 1000000)},
			expPass:          true,
		},
		{
			desc:              "oracle pool join - weight recovering on imbalanced pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uusdt", 1000000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusdc", 1500000), sdk.NewInt64Coin("uusdt", 500000)},
			poolParams: types.PoolParams{
				SwapFee:                     sdk.ZeroDec(),
				ExitFee:                     sdk.ZeroDec(),
				UseOracle:                   true,
				WeightBreakingFeeMultiplier: sdk.NewDecWithPrec(1, 0), // 1.00
				SlippageReduction:           sdk.ZeroDec(),
				LpFeePortion:                sdk.ZeroDec(),
				StakingFeePortion:           sdk.ZeroDec(),
				WeightRecoveryFeePortion:    sdk.ZeroDec(),
				ThresholdWeightDifference:   sdk.NewDecWithPrec(2, 1), // 20%
				FeeDenom:                    "uusdc",
			},
			shareOutAmount:   sdk.NewInt(1250000000000000000), // weight breaking fee
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{sdk.NewInt64Coin("uusdt", 1000000)},
			expPass:          true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			suite.SetupStableCoinPrices()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.poolInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, tc.poolInitBalance)
			suite.Require().NoError(err)

			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
				Denom:     tc.poolInitBalance[0].Denom,
				Liquidity: tc.poolInitBalance[0].Amount,
			})
			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
				Denom:     tc.poolInitBalance[1].Denom,
				Liquidity: tc.poolInitBalance[1].Amount,
			})

			// setup pool to join
			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams:        tc.poolParams,
				TotalShares:       sdk.NewCoin("amm/pool/1", sdk.NewInt(2).Mul(types.OneShare)),
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
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)

			// execute function
			msgServer := keeper.NewMsgServerImpl(suite.app.AmmKeeper)
			resp, err := msgServer.JoinPool(
				sdk.WrapSDKContext(suite.ctx),
				&types.MsgJoinPool{
					Sender:         sender.String(),
					PoolId:         1,
					MaxAmountsIn:   tc.senderInitBalance,
					ShareAmountOut: tc.shareOutAmount,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(resp.ShareAmountOut, tc.shareOutAmount)
				suite.Require().Equal(sdk.Coins(resp.TokenIn).String(), tc.expTokenIn.String())

				pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)
				suite.Require().Len(pools, 1)
				suite.Require().Equal(pools[0].PoolId, uint64(1))
				suite.Require().Equal(pools[0].PoolParams, tc.poolParams)
				suite.Require().Equal(pools[0].TotalShares.Amount.String(), pool.TotalShares.Amount.Add(tc.shareOutAmount).String())

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())

				// check lp token commitment
				commitments, found := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, sender.String())
				suite.Require().True(found)
				suite.Require().Len(commitments.CommittedTokens, 1)
				suite.Require().Equal(commitments.CommittedTokens[0].Denom, "amm/pool/1")
				suite.Require().Equal(commitments.CommittedTokens[0].Amount.String(), tc.shareOutAmount.String())
			}
		})
	}
}
