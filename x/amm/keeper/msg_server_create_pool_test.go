package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *KeeperTestSuite) TestMsgServerCreatePool() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		poolParams        types.PoolParams
		poolAssets        []types.PoolAsset
		expSenderBalance  sdk.Coins
		expTotalLiquidity sdk.Coins
		expLpCommitment   sdk.Coin
		expPass           bool
	}{
		{
			desc:              "zero tvl pool creation",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 1000000)},
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
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdk.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.Elys, 1000000),
					Weight: sdk.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.NewCoin("amm/pool/1", sdk.NewInt(100).Mul(types.OneShare)),
			expPass:          true,
		},
		{
			desc:              "positive tvl pool creation",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
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
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdk.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000),
					Weight: sdk.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.NewCoin("amm/pool/1", sdk.NewInt(2).Mul(types.OneShare)),
			expPass:          true,
		},
		{
			desc:              "not enough balance to create pool",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000)},
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
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdk.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000),
					Weight: sdk.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.Coin{},
			expPass:          false,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			suite.SetupStableCoinPrices()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.senderInitBalance)
			suite.Require().NoError(err)

			// execute function
			msgServer := keeper.NewMsgServerImpl(suite.app.AmmKeeper)
			resp, err := msgServer.CreatePool(
				sdk.WrapSDKContext(suite.ctx),
				&types.MsgCreatePool{
					Sender:     sender.String(),
					PoolParams: &tc.poolParams,
					PoolAssets: tc.poolAssets,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(resp.PoolID, uint64(1))

				pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)
				suite.Require().Len(pools, 1)
				suite.Require().Equal(pools[0].PoolId, uint64(1))
				suite.Require().Equal(pools[0].PoolParams, tc.poolParams)
				suite.Require().Equal(pools[0].TotalShares.Amount.String(), tc.expLpCommitment.Amount.String())

				totalWeight := sdk.ZeroInt()
				for _, poolAsset := range tc.poolAssets {
					totalWeight = totalWeight.Add(poolAsset.Weight)
				}
				suite.Require().Equal(pools[0].TotalWeight.String(), totalWeight.MulRaw(types.GuaranteedWeightPrecision).String())

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())

				// check lp token commitment
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, sender.String())
				suite.Require().Len(commitments.CommittedTokens, 1)
				suite.Require().Equal(commitments.CommittedTokens[0].Denom, tc.expLpCommitment.Denom)
				suite.Require().Equal(commitments.CommittedTokens[0].Amount.String(), tc.expLpCommitment.Amount.String())
			}
		})
	}
}
