package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/v5/x/amm/keeper"
	"github.com/elys-network/elys/v5/x/amm/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *AmmKeeperTestSuite) TestMsgServerJoinPool() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		poolInitBalance   sdk.Coins
		poolParams        types.PoolParams
		shareOutAmount    math.Int
		expSenderBalance  sdk.Coins
		expTotalLiquidity sdk.Coins
		expTokenIn        sdk.Coins
		expPass           bool
	}{
		{
			desc:              "successful non-oracle join pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100000), sdk.NewInt64Coin("uusdt", 100000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:   math.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			shareOutAmount:   types.OneShare.Quo(math.NewInt(5)),
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100000), sdk.NewInt64Coin("uusdt", 100000)},
			expPass:          true,
		},
		{
			desc:              "not enough balance to join pool - non-oracle pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:   math.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			shareOutAmount:   types.OneShare.Quo(math.NewInt(5)),
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{},
			expPass:          false,
		},
		{
			desc:              "oracle pool join - breaking weight on balanced pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uusdt", 1000000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:   math.LegacyZeroDec(),
				UseOracle: true,
				FeeDenom:  ptypes.BaseCurrency,
			},
			// shareOutAmount:   math.NewInt(694444166666666666), // weight breaking fee - slippage enable
			shareOutAmount:   math.NewInt(832154489169426942), // weight breaking fee - slippage enable
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{sdk.NewInt64Coin("uusdt", 1000000)},
			expPass:          true,
		},
		{
			desc:              "oracle pool join - weight recovering on imbalanced pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin("uusdt", 1000000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1500000), sdk.NewInt64Coin("uusdt", 500000)},
			poolParams: types.PoolParams{
				SwapFee:   math.LegacyZeroDec(),
				UseOracle: true,
				FeeDenom:  ptypes.BaseCurrency,
			},
			// shareOutAmount:   math.NewInt(805987500000000000), // weight recovery direction - slippage enable
			shareOutAmount:   math.NewInt(644790000000000000), // weight recovery direction - slippage enable
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{sdk.NewInt64Coin("uusdt", 1000000)},
			expPass:          true,
		},
		{
			desc:              "oracle pool join - zero slippage add liquidity",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1500000), sdk.NewInt64Coin("uusdt", 500000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1500000), sdk.NewInt64Coin("uusdt", 500000)},
			poolParams: types.PoolParams{
				SwapFee:   math.LegacyZeroDec(),
				UseOracle: true,
				FeeDenom:  ptypes.BaseCurrency,
			},
			shareOutAmount:   math.NewInt(2000000000000000000),
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1500000), sdk.NewInt64Coin("uusdt", 500000)},
			expPass:          true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			suite.SetupStableCoinPrices()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolAddr := types.NewPoolAddress(uint64(1))
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
				TotalShares:       sdk.NewCoin("amm/pool/1", math.NewInt(2).Mul(types.OneShare)),
				PoolAssets: []types.PoolAsset{
					{
						Token:  tc.poolInitBalance[0],
						Weight: math.NewInt(10),
					},
					{
						Token:  tc.poolInitBalance[1],
						Weight: math.NewInt(10),
					},
				},
				TotalWeight: math.ZeroInt(),
			}
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)
			suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

			// execute function
			msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
			resp, err := msgServer.JoinPool(
				suite.ctx,
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
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

				pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)
				suite.Require().Len(pools, 1)
				suite.Require().Equal(pools[0].PoolId, uint64(1))
				suite.Require().Equal(pools[0].PoolParams, tc.poolParams)
				suite.Require().Equal(pools[0].TotalShares.Amount.String(), pool.TotalShares.Amount.Add(tc.shareOutAmount).String())

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())

				// check lp token commitment
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, sender)
				suite.Require().Len(commitments.CommittedTokens, 1)
				suite.Require().Equal(commitments.CommittedTokens[0].Denom, "amm/pool/1")
				suite.Require().Equal(commitments.CommittedTokens[0].Amount.String(), tc.shareOutAmount.String())
			}
		})
	}
}

func (suite *AmmKeeperTestSuite) TestMsgServerJoinPoolExploitScenario() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		poolInitBalance   sdk.Coins
		poolParams        types.PoolParams
		shareOutAmount    math.Int
		expSenderBalance  sdk.Coins
		expTotalLiquidity sdk.Coins
		expTokenIn        sdk.Coins
		expPass           bool
	}{
		{
			desc:              "Exploit scenario for Join Pool - unfair liquidity extraction",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.ATOM, 100_000_000_000000)},
			poolInitBalance:   sdk.Coins{sdk.NewInt64Coin(ptypes.ATOM, 100_000_000_000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000_000_000000)},
			poolParams: types.PoolParams{
				SwapFee:   math.LegacyZeroDec(),
				UseOracle: true,
				FeeDenom:  ptypes.BaseCurrency,
			},
			shareOutAmount:   math.NewInt(2_000000000000000000),
			expSenderBalance: sdk.Coins{},
			expTokenIn:       sdk.Coins{sdk.NewInt64Coin(ptypes.ATOM, 1_000000)},
			expPass:          false,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			suite.SetupCoinPrices()

			// Step 1: Bootstrap accounts
			// Create sender, pool, and treasury accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolAddr := types.NewPoolAddress(1)
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// Step 2: Bootstrap balances
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

			// Step 3: Setup initial pool with 50:50 weight
			pool := types.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams:        tc.poolParams,
				TotalShares:       sdk.NewCoin("amm/pool/1", math.NewInt(2).Mul(types.OneShare)),
				PoolAssets: []types.PoolAsset{
					{
						Token:                  tc.poolInitBalance[0],
						Weight:                 math.NewInt(1),
						ExternalLiquidityRatio: math.LegacyNewDec(1),
					},
					{
						Token:                  tc.poolInitBalance[1],
						Weight:                 math.NewInt(1),
						ExternalLiquidityRatio: math.LegacyNewDec(1),
					},
				},
				TotalWeight: math.ZeroInt(),
			}
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)
			suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

			// Step 4: Simulate market price movement - adjust weights to 10:1
			pool.PoolAssets[0].Weight = math.NewInt(10)
			pool.PoolAssets[1].Weight = math.NewInt(1)
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)

			// Step 5: New LP adds single-sided liquidity
			msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
			resp, err := msgServer.JoinPool(
				suite.ctx,
				&types.MsgJoinPool{
					Sender:         sender.String(),
					PoolId:         1,
					MaxAmountsIn:   tc.senderInitBalance,
					ShareAmountOut: tc.shareOutAmount,
				})

			suite.Require().NoError(err)
			suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

			snapshot := suite.app.AmmKeeper.GetPoolWithAccountedBalance(suite.ctx, pool.PoolId)
			initialWeightIn := types.GetDenomOracleAssetWeight(suite.ctx, suite.app.OracleKeeper, snapshot.PoolAssets, tc.senderInitBalance[0].Denom)
			initialWeightOut := osmomath.OneBigDec().Sub(initialWeightIn)

			// Step 6: Validate if exploit was successful (It should fail)
			// Calculate expected number of shares without weight balance bonus
			totalShares := pool.TotalShares.Amount
			joinValueWithSlippage, _, _ := pool.CalcJoinValueWithSlippage(suite.ctx, snapshot, suite.app.OracleKeeper, tc.senderInitBalance[0], initialWeightOut, types.DefaultParams())
			tvl, _ := pool.TVL(suite.ctx, suite.app.OracleKeeper, suite.app.AccountedPoolKeeper)
			expectedNumShares := osmomath.BigDecFromSDKInt(totalShares).
				Mul(joinValueWithSlippage).Quo(tvl).Dec().RoundInt()

			// Number of shares must be lesser or equal to expected
			suite.Require().GreaterOrEqual(expectedNumShares.String(), resp.ShareAmountOut.String(), "Exploit detected: Sender received more shares than expected")
		})
	}
}
