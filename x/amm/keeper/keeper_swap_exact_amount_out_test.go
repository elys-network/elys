package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/v4/x/amm/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *AmmKeeperTestSuite) TestSwapExactAmountOut() {
	for _, tc := range []struct {
		desc                string
		senderInitBalance   sdk.Coins
		poolInitBalance     sdk.Coins
		treasuryInitBalance sdk.Coins
		swapFeeOut          osmomath.BigDec
		tokenIn             sdk.Coin
		tokenInMax          sdkmath.Int
		tokenOut            sdk.Coin
		weightBalanceBonus  osmomath.BigDec
		isOraclePool        bool
		useNewRecipient     bool
		expSenderBalance    sdk.Coins
		expRecipientBalance sdk.Coins
		expPoolBalance      sdk.Coins
		expTreasuryBalance  sdk.Coins
		expPass             bool
		errMsg              string
	}{
		{
			desc:                "tokenIn is same as tokenOut",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusdc", 10000),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  osmomath.ZeroBigDec(),
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{},
			expTreasuryBalance:  sdk.Coins{},
			expPass:             false,
			errMsg:              "cannot trade the same denomination in and out",
		},
		{
			desc:                "tokenIn is 0 corrosponsing to tokenOut",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 1000),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 0),
			weightBalanceBonus:  osmomath.ZeroBigDec(),
			isOraclePool:        true,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{},
			expTreasuryBalance:  sdk.Coins{},
			expPass:             false,
			errMsg:              "amount too low",
		},
		{
			desc:                "MaxTokenIn is less than required tokenIn amount",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenInMax:          sdkmath.NewInt(10),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 9802),
			weightBalanceBonus:  osmomath.ZeroBigDec(),
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 990000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1009802)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010000), sdk.NewInt64Coin(ptypes.BaseCurrency, 990100)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             false,
			errMsg:              "calculated amount is larger than max amount",
		},
		{
			desc:                "pool does not enough balance for out",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  osmomath.ZeroBigDec(),
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
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10000),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  osmomath.ZeroBigDec(),
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
			swapFeeOut:          osmomath.NewBigDecWithPrec(1, 2), // 1%
			tokenIn:             sdk.NewInt64Coin("uusda", 10204),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  osmomath.ZeroBigDec(),
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 989796), sdk.NewInt64Coin(ptypes.BaseCurrency, 1010000)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010204), sdk.NewInt64Coin(ptypes.BaseCurrency, 989901)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful execution with zero swap fee",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10102),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  osmomath.ZeroBigDec(),
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 989898), sdk.NewInt64Coin(ptypes.BaseCurrency, 1010000)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010102), sdk.NewInt64Coin(ptypes.BaseCurrency, 990000)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful execution with positive slippage on oracle pool",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 222546),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 200000),
			weightBalanceBonus:  osmomath.ZeroBigDec(),
			isOraclePool:        true,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 777454), sdk.NewInt64Coin(ptypes.BaseCurrency, 1200000)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1222402), sdk.NewInt64Coin(ptypes.BaseCurrency, 800000)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000144), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful weight bonus & huge amount rebalance treasury",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10102),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			isOraclePool:        false,
			useNewRecipient:     false,
			weightBalanceBonus:  osmomath.NewBigDecWithPrec(3, 1), // 30% bonus
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 989898), sdk.NewInt64Coin(ptypes.BaseCurrency, 1010000)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010102), sdk.NewInt64Coin(ptypes.BaseCurrency, 990000)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:             true,
		},
		{
			desc:                "successful weight bonus & lack of rebalance treasury",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10102),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  osmomath.NewBigDecWithPrec(3, 1), // 30% bonus
			isOraclePool:        false,
			useNewRecipient:     false,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 989898), sdk.NewInt64Coin(ptypes.BaseCurrency, 1010000)},
			expRecipientBalance: sdk.Coins{},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010102), sdk.NewInt64Coin(ptypes.BaseCurrency, 990000)},
			expTreasuryBalance:  sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			expPass:             true,
		},
		{
			desc:                "new recipient address",
			senderInitBalance:   sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:     sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			treasuryInitBalance: sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100)},
			swapFeeOut:          osmomath.ZeroBigDec(),
			tokenIn:             sdk.NewInt64Coin("uusda", 10102),
			tokenInMax:          sdkmath.NewInt(10000000),
			tokenOut:            sdk.NewInt64Coin(ptypes.BaseCurrency, 10000),
			weightBalanceBonus:  osmomath.NewBigDecWithPrec(3, 1), // 30% bonus
			isOraclePool:        false,
			useNewRecipient:     true,
			expSenderBalance:    sdk.Coins{sdk.NewInt64Coin("uusda", 989898), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expRecipientBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 10000)},
			expPoolBalance:      sdk.Coins{sdk.NewInt64Coin("uusda", 1010102), sdk.NewInt64Coin(ptypes.BaseCurrency, 990000)},
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
			poolAddr := types.NewPoolAddress(uint64(1))
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
					UseOracle: tc.isOraclePool,
					SwapFee:   tc.swapFeeOut.Dec(),
					FeeDenom:  ptypes.BaseCurrency,
				},
				TotalShares: sdk.NewCoin(types.GetPoolShareDenom(1), sdkmath.ZeroInt()),
				PoolAssets: []types.PoolAsset{
					{
						Token:                  tc.poolInitBalance[0],
						Weight:                 sdkmath.NewInt(10),
						ExternalLiquidityRatio: sdkmath.LegacyNewDec(2),
					},
					{
						Token:                  tc.poolInitBalance[1],
						Weight:                 sdkmath.NewInt(10),
						ExternalLiquidityRatio: sdkmath.LegacyNewDec(2),
					},
				},
				TotalWeight: sdkmath.ZeroInt(),
			}
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)
			suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

			tokenInAmount, _, err := suite.app.AmmKeeper.InternalSwapExactAmountOut(suite.ctx, sender, recipient, pool, tc.tokenIn.Denom, tc.tokenInMax, tc.tokenOut, tc.swapFeeOut, osmomath.ZeroBigDec())
			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tokenInAmount.String(), tc.tokenIn.Amount.String())
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))
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
				suite.Require().Equal(liquidity.Liquidity, tc.expPoolBalance.AmountOf(tc.tokenOut.Denom))

				track := suite.app.AmmKeeper.GetSlippageTrack(suite.ctx, 1, uint64(suite.ctx.BlockTime().Unix()))
				if tc.isOraclePool {
					suite.Require().Equal(track.Tracked.String(), "22224uusda")
				} else {
					suite.Require().Equal(track.Tracked.String(), "")
				}
			}
		})
	}
}
