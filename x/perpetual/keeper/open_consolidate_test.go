package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestOpenConsolidate() {
	testCases := []struct {
		name           string
		setup          func() (*types.MsgOpen, *types.MTP, *types.MTP)
		expectedErrMsg string
	}{
		{
			"Pool does not exist",
			func() (*types.MsgOpen, *types.MTP, *types.MTP) {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, ammPool.PoolId, positionCreator, position.Id)
				suite.Require().NoError(err)

				suite.app.AmmKeeper.RemovePool(suite.ctx, firstPool)

				return openPositionMsg, &mtp, &mtp
			},
			"perpetual pool does not exist",
		},
		{
			"Force Closed: Mtp health will be low for the safety factor",
			func() (*types.MsgOpen, *types.MTP, *types.MTP) {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				suite.SetPerpetualPool(1)
				_, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				amount := math.NewInt(400)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, firstPool, positionCreator, position.Id)
				suite.Require().NoError(err)

				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.SafetyFactor = math.LegacyMustNewDecFromStr("1.30")
				suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)

				return openPositionMsg, &mtp, &mtp
			},
			"",
		},
		{
			"Sucess: MTP consolidation",
			func() (*types.MsgOpen, *types.MTP, *types.MTP) {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				suite.SetPerpetualPool(1)
				_, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				amount := math.NewInt(400)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, firstPool, positionCreator, position.Id)
				suite.Require().NoError(err)

				return openPositionMsg, &mtp, &mtp
			},
			"",
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg, existingMtp, newMtp := tc.setup()
			_, err := suite.app.PerpetualKeeper.OpenConsolidate(suite.ctx, existingMtp, newMtp, msg, ptypes.ATOM, types.NewPerpetualFeesWithEmptyCoins())

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *PerpetualKeeperTestSuite) TestOpenConsolidateUsingOpen() {
	var initialPoolBankBalance sdk.Coins
	var initialAccountedPoolBalance sdk.Coins

	var finalPoolBankBalance sdk.Coins
	var finalAccountedPoolBalance sdk.Coins

	var ammPool ammtypes.Pool
	var msg types.MsgOpen

	testCases := []struct {
		name            string
		setup           func() *types.MsgOpen
		expectedErrMsg  string
		consolidatedMtp *types.MTP
		postValidate    func(msg *types.MsgOpen)
	}{
		{
			"Success: Consolidate two position with different leverage and take profit price",
			func() *types.MsgOpen {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool = suite.SetPerpetualPool(1)
				_, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				amount := math.NewInt(400)
				msg = types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				_, err = suite.app.PerpetualKeeper.Open(suite.ctx, &msg)
				suite.Require().NoError(err)

				msg.Leverage = math.LegacyNewDec(3)
				msg.TakeProfitPrice = math.LegacyMustNewDecFromStr("1.5")

				initialPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				initialAccountedPoolBalance = accountedPool.TotalTokens
				return &msg
			},
			"",
			&types.MTP{
				Collateral:      math.NewInt(800),
				Liabilities:     math.NewInt(650),
				Custody:         math.NewInt(4000),
				TakeProfitPrice: math.LegacyMustNewDecFromStr("1.5"),
			},
			func(msg *types.MsgOpen) {
				finalPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				finalAccountedPoolBalance = accountedPool.TotalTokens

				suite.Require().Equal(initialPoolBankBalance.Add(msg.Collateral), finalPoolBankBalance)
				atleastExpected := initialAccountedPoolBalance.Add(msg.Collateral).Add(sdk.NewCoin(msg.Collateral.Denom, msg.Leverage.Sub(math.LegacyOneDec()).MulInt(msg.Collateral.Amount).TruncateInt()))
				suite.Require().True(finalAccountedPoolBalance.AmountOf(ptypes.ATOM).GTE(atleastExpected.AmountOf(ptypes.ATOM)))
				suite.Require().True(finalAccountedPoolBalance.AmountOf(ptypes.BaseCurrency).LTE(atleastExpected.AmountOf(ptypes.BaseCurrency)))
			},
		},
		{
			"Success: add collateral for the existing position",
			func() *types.MsgOpen {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				suite.SetPerpetualPool(1)
				_, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				amount := math.NewInt(400)
				msg = types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				_, err = suite.app.PerpetualKeeper.Open(suite.ctx, &msg)
				suite.Require().NoError(err)

				// make new Position leverage 0 to add collateral
				msg.Leverage = math.LegacyNewDec(0)

				return &msg
			},
			"",
			&types.MTP{
				Collateral:      math.NewInt(800),
				Liabilities:     math.NewInt(406),
				Custody:         math.NewInt(2800),
				TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
			},
			func(msg *types.MsgOpen) {
				finalPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				finalAccountedPoolBalance = accountedPool.TotalTokens

				atleastExpected := initialAccountedPoolBalance.Add(msg.Collateral)
				suite.Require().Equal(initialPoolBankBalance.Add(msg.Collateral), finalPoolBankBalance)
				suite.Require().True(finalAccountedPoolBalance.AmountOf(ptypes.ATOM).GTE(atleastExpected.AmountOf(ptypes.ATOM)))
				suite.Require().True(finalAccountedPoolBalance.AmountOf(ptypes.BaseCurrency).LTE(atleastExpected.AmountOf(ptypes.BaseCurrency)))
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.setup()
			ammPool, _ = suite.app.AmmKeeper.GetPool(suite.ctx, ammPool.PoolId)
			position, err := suite.app.PerpetualKeeper.Open(suite.ctx, msg)
			suite.Require().NoError(err)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				consolidateMtp, mtpErr := suite.app.PerpetualKeeper.GetMTP(suite.ctx, uint64(1), sdk.MustAccAddressFromBech32(msg.Creator), position.Id)
				suite.Require().NoError(mtpErr)
				suite.Require().Equal(tc.consolidatedMtp.Collateral, consolidateMtp.Collateral)
				suite.Require().Equal(tc.consolidatedMtp.Liabilities, consolidateMtp.Liabilities)
				suite.Require().Equal(tc.consolidatedMtp.Custody, consolidateMtp.Custody)
				suite.Require().Equal(tc.consolidatedMtp.TakeProfitPrice, consolidateMtp.TakeProfitPrice)
			}
			tc.postValidate(msg)
		})
	}
}
