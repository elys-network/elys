package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *KeeperTestSuite) TestSwapExactAmountIn() {
	for _, tc := range []struct {
		desc                string
		senderInitBalance   sdk.Coins
		poolInitBalance     sdk.Coins
		treasuryInitBalance sdk.Coins
		swapFeeIn           sdk.Dec
		swapFeeOut          sdk.Dec
		tokenIn             sdk.Coin
		tokenOutMin         math.Int
		tokenOut            sdk.Coin
		weightBalanceBonus  sdk.Dec
		isOraclePool        bool
		useNewRecipient     bool
		expSenderBalance    sdk.Coins
		expRecipientBalance sdk.Coins
		expPoolBalance      sdk.Coins
		expTreasuryBalance  sdk.Coins
		expPass             bool
	}{
		{
			desc:                "pool does not enough balance for out",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  sdk.ZeroDec(),
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{},
			expTreasuryBalance:  sdk.Coins{},
			expPass:             false,
		},
		{
			desc:                "sender does not have enough balance for in",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 100)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  sdk.ZeroDec(),
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{},
			expTreasuryBalance:  sdk.Coins{},
			expPass:             false,
		},
		{
			desc:                "successful execution with positive swap fee",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.NewDecWithPrec(1, 2), // 1%
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9704),
			weightBalanceBonus:  sdk.ZeroDec(),
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009704)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990198)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful execution with zero swap fee",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			weightBalanceBonus:  sdk.ZeroDec(),
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009900)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990100)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful execution with positive slippage on oracle pool",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9975),
			weightBalanceBonus:  sdk.ZeroDec(),
			isOraclePool:        true,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009975)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990025)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful weight bonus & huge amount rebalance treasury",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			weightBalanceBonus:  sdk.NewDecWithPrec(3, 1), // 30% bonus
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009900)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990100)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful weight bonus & lack of rebalance treasury",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			weightBalanceBonus:  sdk.NewDecWithPrec(3, 1), // 30% bonus
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009900)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990100)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			expPass:             true,
		},
		{
			desc:                "new recipient address",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			weightBalanceBonus:  sdk.NewDecWithPrec(3, 1), // 30% bonus
			isOraclePool:        false,
			useNewRecipient:     true,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expRecipientBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 9900)},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990100)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			expPass:             true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			suite.SetupStableCoinPrices()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			recipient := sender
			poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			if tc.useNewRecipient {
				recipient = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			}

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.poolInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, tc.poolInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.treasuryInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, treasuryAddr, tc.treasuryInitBalance)
			suite.Require().NoError(err)

			// execute function
			for _, coin := range tc.poolInitBalance {
				suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
					Denom:     coin.Denom,
					Liquidity: coin.Amount,
				})
			}
			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: types.PoolParams{
					UseOracle:                   tc.isOraclePool,
					ExternalLiquidityRatio:      sdk.NewDec(2),
					WeightBreakingFeeMultiplier: sdk.ZeroDec(),
					WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
					WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
					ThresholdWeightDifference:   sdk.ZeroDec(),
					SwapFee:                     tc.swapFeeIn,
					FeeDenom:                    ptypes.BaseCurrency,
				},
				TotalShares: sdk.Coin{},
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

			tokenOut, err := suite.app.AmmKeeper.SwapExactAmountIn(suite.ctx, sender, recipient, pool, tc.tokenIn, tc.tokenOut.Denom, tc.tokenOutMin, tc.swapFeeIn)
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tokenOut.String(), tc.tokenOut.Amount.String())

				// check pool balance increase/decrease
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
				suite.Require().Equal(balances.String(), tc.expPoolBalance.String())

				// check balance change on sender
				balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())

				if tc.useNewRecipient {
					// check balance change on recipient
					balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, recipient)
					suite.Require().Equal(balances.String(), tc.expRecipientBalance.String())
				}

				// check balance change on treasury
				balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, treasuryAddr)
				suite.Require().Equal(balances.String(), tc.expTreasuryBalance.String())

				// check pool object change
				savedPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, pool.PoolId)
				suite.Require().True(found)
				suite.Require().Equal(savedPool.PoolParams.SwapFee, pool.PoolParams.SwapFee)

				// check total liquidity change (increase + decrease)
				liquidity, found := suite.app.AmmKeeper.GetDenomLiquidity(suite.ctx, tc.tokenIn.Denom)
				suite.Require().True(found)
				suite.Require().Equal(liquidity.Liquidity, tc.expPoolBalance.AmountOf(tc.tokenIn.Denom))
				liquidity, found = suite.app.AmmKeeper.GetDenomLiquidity(suite.ctx, tc.tokenOut.Denom)
				suite.Require().True(found)
				suite.Require().True(liquidity.Liquidity.LTE(tc.expPoolBalance.AmountOf(tc.tokenOut.Denom)))

				track := suite.app.AmmKeeper.GetSlippageTrack(suite.ctx, 1, uint64(suite.ctx.BlockTime().Unix()))
				if tc.isOraclePool {
					suite.Require().Equal(track.Tracked.String(), "25uusdc")
				} else {
					suite.Require().Equal(track.Tracked.String(), "")
				}
			}
		})
	}
}
