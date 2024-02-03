package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *KeeperTestSuite) TestRouteExactAmountIn() {
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
		expSenderBalance    sdk.Coins
		expPoolBalance      sdk.Coins
		expTreasuryBalance  sdk.Coins
		expPass             bool
	}{
		{
			desc:                "pool does not enough balance for out",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  sdk.ZeroDec(),
			expSenderBalance:    sdk.Coins{},
			expPoolBalance:      sdk.Coins{},
			expTreasuryBalance:  sdk.Coins{},
			expPass:             false,
		},
		{
			desc:                "sender does not have enough balance for in",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 100)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  sdk.ZeroDec(),
			expSenderBalance:    sdk.Coins{},
			expPoolBalance:      sdk.Coins{},
			expTreasuryBalance:  sdk.Coins{},
			expPass:             false,
		},
		{
			desc:                "successful execution with positive swap fee",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.NewDecWithPrec(1, 2), // 1%
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9704),
			weightBalanceBonus:  sdk.ZeroDec(),
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009704)},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990198)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful execution with zero swap fee",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			weightBalanceBonus:  sdk.ZeroDec(),
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009900)},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990100)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful weight bonus & huge amount rebalance treasury",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			weightBalanceBonus:  sdk.NewDecWithPrec(3, 1), // 30% bonus
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009900)},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990100)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful weight bonus & lack of rebalance treasury",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			swapFeeIn:           sdk.ZeroDec(),
			swapFeeOut:          sdk.ZeroDec(),
			tokenIn:             sdk.NewInt64Coin(ptypes.Elys, 10000),
			tokenOutMin:         sdk.ZeroInt(),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9900),
			weightBalanceBonus:  sdk.NewDecWithPrec(3, 1), // 30% bonus
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009900)},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990100)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			expPass:             true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

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
					SwapFee:                  tc.swapFeeIn,
					FeeDenom:                 ptypes.BaseCurrency,
					WeightRecoveryFeePortion: sdk.NewDecWithPrec(10, 2), // 10%
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
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)

			// TODO: add multiple route case
			// TODO: add invalid route case
			// TODO: add Elys token involved route case
			tokenOut, _, _, err := suite.app.AmmKeeper.RouteExactAmountIn(
				suite.ctx,
				sender,
				sender,
				[]types.SwapAmountInRoute{
					{
						PoolId:        1,
						TokenOutDenom: tc.tokenOut.Denom,
					},
				},
				tc.tokenIn, tc.tokenOutMin,
				sdk.ZeroDec(),
			)
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
			}
		})
	}
}
