package keeper_test

import (
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *KeeperTestSuite) TestMsgServerExitPool() {
	for _, tc := range []struct {
		desc              string
		poolInitBalance   sdk.Coins
		poolParams        types.PoolParams
		shareInAmount     sdk.Int
		tokenOutDenom     string
		minAmountsOut     sdk.Coins
		expSenderBalance  sdk.Coins
		expTotalLiquidity sdk.Coins
		expPass           bool
	}{
		{
			desc:            "successful non-oracle exit pool",
			poolInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:                     sdk.ZeroDec(),
				ExitFee:                     sdk.ZeroDec(),
				UseOracle:                   false,
				WeightBreakingFeeMultiplier: sdk.ZeroDec(),
				WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
				ExternalLiquidityRatio:      sdk.NewDec(1),
				WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
				ThresholdWeightDifference:   sdk.ZeroDec(),
				FeeDenom:                    ptypes.BaseCurrency,
			},
			shareInAmount:    types.OneShare.Quo(sdk.NewInt(5)),
			tokenOutDenom:    "",
			minAmountsOut:    sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100000), sdk.NewInt64Coin("uusdt", 100000)},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100000), sdk.NewInt64Coin("uusdt", 100000)},
			expPass:          true,
		},
		{
			desc:            "not enough balance to exit pool - non-oracle pool",
			poolInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:                     sdk.ZeroDec(),
				ExitFee:                     sdk.ZeroDec(),
				UseOracle:                   false,
				WeightBreakingFeeMultiplier: sdk.ZeroDec(),
				WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
				ExternalLiquidityRatio:      sdk.NewDec(1),
				WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
				ThresholdWeightDifference:   sdk.ZeroDec(),
				FeeDenom:                    ptypes.BaseCurrency,
			},
			shareInAmount:    types.OneShare.Quo(sdk.NewInt(5)),
			tokenOutDenom:    "",
			minAmountsOut:    sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expSenderBalance: sdk.Coins{},
			expPass:          false,
		},
		{
			desc:            "oracle pool exit - breaking weight on balanced pool",
			poolInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:                     sdk.ZeroDec(),
				ExitFee:                     sdk.ZeroDec(),
				UseOracle:                   true,
				WeightBreakingFeeMultiplier: sdk.NewDecWithPrec(1, 2),  // 0.01
				WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
				ExternalLiquidityRatio:      sdk.NewDec(1),
				WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
				ThresholdWeightDifference:   sdk.NewDecWithPrec(2, 1),  // 20%
				FeeDenom:                    ptypes.BaseCurrency,
			},
			shareInAmount: types.OneShare.Quo(sdk.NewInt(10)),
			tokenOutDenom: "uusdt",
			minAmountsOut: sdk.Coins{sdk.NewInt64Coin("uusdt", 98699)},
			// expSenderBalance: sdk.Coins{sdk.NewInt64Coin("uusdt", 95114)}, // slippage enabled
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin("uusdt", 98699)}, // slippage disabled
			expPass:          true,
		},
		{
			desc:            "oracle pool exit - weight recovering on imbalanced pool",
			poolInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1500000), sdk.NewInt64Coin("uusdt", 500000)},
			poolParams: types.PoolParams{
				SwapFee:                     sdk.ZeroDec(),
				ExitFee:                     sdk.ZeroDec(),
				UseOracle:                   true,
				WeightBreakingFeeMultiplier: sdk.NewDecWithPrec(1, 2),  // 0.01
				WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
				ExternalLiquidityRatio:      sdk.NewDec(1),
				WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
				ThresholdWeightDifference:   sdk.NewDecWithPrec(2, 1),  // 20%
				FeeDenom:                    ptypes.BaseCurrency,
			},
			shareInAmount: types.OneShare.Quo(sdk.NewInt(10)),
			tokenOutDenom: ptypes.BaseCurrency,
			minAmountsOut: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100000)},
			// expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 99197)}, // slippage enabled
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100000)}, // slippage disabled
			expPass:          true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			suite.SetupStableCoinPrices()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.poolInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.poolInitBalance)
			suite.Require().NoError(err)

			// execute function
			msgServer := keeper.NewMsgServerImpl(suite.app.AmmKeeper)
			poolAssets := []types.PoolAsset{
				{
					Token:  tc.poolInitBalance[0],
					Weight: sdk.NewInt(10),
				},
				{
					Token:  tc.poolInitBalance[1],
					Weight: sdk.NewInt(10),
				},
			}
			_, err = msgServer.CreatePool(
				sdk.WrapSDKContext(suite.ctx),
				&types.MsgCreatePool{
					Sender:     sender.String(),
					PoolParams: &tc.poolParams,
					PoolAssets: poolAssets,
				})
			suite.Require().NoError(err)
			pool := suite.app.AmmKeeper.GetAllPool(suite.ctx)[0]
			suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour))
			resp, err := msgServer.ExitPool(
				sdk.WrapSDKContext(suite.ctx),
				&types.MsgExitPool{
					Sender:        sender.String(),
					PoolId:        pool.PoolId,
					MinAmountsOut: tc.minAmountsOut,
					ShareAmountIn: tc.shareInAmount,
					TokenOutDenom: tc.tokenOutDenom,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(sdk.Coins(resp.TokenOut).String(), tc.minAmountsOut.String())

				pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)
				suite.Require().Len(pools, 1)
				suite.Require().Equal(pools[0].PoolId, uint64(1))
				suite.Require().Equal(pools[0].PoolParams, tc.poolParams)
				suite.Require().Equal(pools[0].TotalShares.Amount.String(), pool.TotalShares.Amount.Sub(tc.shareInAmount).String())

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())

				// check lp token commitment
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, sender.String())
				suite.Require().Len(commitments.CommittedTokens, 1)
				suite.Require().Equal(commitments.CommittedTokens[0].Denom, "amm/pool/1")
				suite.Require().Equal(commitments.CommittedTokens[0].Amount.String(), pool.TotalShares.Amount.Sub(tc.shareInAmount).String())
			}
		})
	}
}

// TODO: check combined scenario - $500 JUNO sell, $500 JUNO buy (weight breaking)
// TODO: check combined scenario - $500 JUNO sell, $500 JUNO buy (weight not breaking)
// TODO: run simulation test with a lot of traffic, and see pool status after the execution
// TODO: Check maximum weight breaking fee applied
// TODO: Check maximum weight recovery bonus applied
// TODO: Check weight recovery treasury empty case
// TODO: handle case bonus pool does not have enough amount
// TODO: check fee distribution
// TODO: write table driven data on spec folder for various cases to show the comparison with Osmosis
