package keeper_test

import (
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestMsgServerExitPool() {
	for _, tc := range []struct {
		desc              string
		poolInitBalance   sdk.Coins
		poolParams        types.PoolParams
		shareInAmount     sdkmath.Int
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
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			shareInAmount:    types.OneShare.Quo(sdkmath.NewInt(5)),
			tokenOutDenom:    "",
			minAmountsOut:    sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100000), sdk.NewInt64Coin("uusdt", 100000)},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100000), sdk.NewInt64Coin("uusdt", 100000)},
			expPass:          true,
		},
		{
			desc:            "not enough balance to exit pool - non-oracle pool",
			poolInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			shareInAmount:    types.OneShare.Quo(sdkmath.NewInt(5)),
			tokenOutDenom:    "",
			minAmountsOut:    sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expSenderBalance: sdk.Coins{},
			expPass:          false,
		},
		{
			desc:            "oracle pool exit - breaking weight on balanced pool",
			poolInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000), sdk.NewInt64Coin("uusdt", 1000000)},
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: true,
				FeeDenom:  ptypes.BaseCurrency,
			},
			shareInAmount:    types.OneShare.Quo(sdkmath.NewInt(10)),
			tokenOutDenom:    "uusdt",
			minAmountsOut:    sdk.Coins{sdk.NewInt64Coin("uusdt", 97619)},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin("uusdt", 97619)}, // slippage enabled
			expPass:          true,
		},
		{
			desc:            "oracle pool exit - weight recovering on imbalanced pool",
			poolInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1500000), sdk.NewInt64Coin("uusdt", 500000)},
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: true,
				FeeDenom:  ptypes.BaseCurrency,
			},
			shareInAmount:    types.OneShare.Quo(sdkmath.NewInt(10)),
			tokenOutDenom:    ptypes.BaseCurrency,
			minAmountsOut:    sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 99197)},
			expSenderBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 99197)}, // slippage enabled
			expPass:          true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			suite.SetupCoinPrices()
			suite.SetAmmParams()
			suite.SetupAssetProfile()

			// bootstrap accounts
			sender := authtypes.NewModuleAddress(govtypes.ModuleName)
			params := suite.app.AmmKeeper.GetParams(suite.ctx)
			// bootstrap balances
			poolCreationFee := sdk.NewCoin(ptypes.Elys, params.PoolCreationFee)
			coins := tc.poolInitBalance.Add(poolCreationFee)
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, coins)
			suite.Require().NoError(err)

			// execute function
			msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
			poolAssets := []types.PoolAsset{
				{
					Token:                  tc.poolInitBalance[0],
					Weight:                 sdkmath.NewInt(10),
					ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
				},
				{
					Token:                  tc.poolInitBalance[1],
					Weight:                 sdkmath.NewInt(10),
					ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
				},
			}
			res, err := msgServer.CreatePool(
				suite.ctx,
				&types.MsgCreatePool{
					Sender:     sender.String(),
					PoolParams: tc.poolParams,
					PoolAssets: poolAssets,
				})
			suite.Require().NoError(err)
			suite.Require().True(suite.VerifyPoolAssetWithBalance(res.PoolID))
			pool := suite.app.AmmKeeper.GetAllPool(suite.ctx)[0]
			suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour))
			resp, err := msgServer.ExitPool(
				suite.ctx,
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
				suite.Require().True(suite.VerifyPoolAssetWithBalance(pool.PoolId))

				pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)
				suite.Require().Len(pools, 1)
				suite.Require().Equal(pools[0].PoolId, uint64(1))
				suite.Require().Equal(pools[0].PoolParams, tc.poolParams)
				suite.Require().Equal(pools[0].TotalShares.Amount.String(), pool.TotalShares.Amount.Sub(tc.shareInAmount).String())

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())

				// check lp token commitment
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, sender)
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
