package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
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
				tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, positionCreator, position.Id)
				suite.Require().NoError(err)

				suite.app.AmmKeeper.RemovePool(suite.ctx, firstPool)

				return openPositionMsg, &mtp, &mtp
			},
			"perpetual pool does not exist",
		},
		{
			"Mtp health will be low for the safety factor",
			func() (*types.MsgOpen, *types.MTP, *types.MTP) {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				suite.SetPerpetualPool(1)
				_, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				amount := math.NewInt(400)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, positionCreator, position.Id)
				suite.Require().NoError(err)

				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.SafetyFactor = math.LegacyMustNewDecFromStr("1.30")
				suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)

				return openPositionMsg, &mtp, &mtp
			},
			"mtp health would be too low for safety factor",
		},
		{
			"Sucess: MTP consolidation",
			func() (*types.MsgOpen, *types.MTP, *types.MTP) {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				suite.SetPerpetualPool(1)
				_, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				amount := math.NewInt(400)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, positionCreator, position.Id)
				suite.Require().NoError(err)

				return openPositionMsg, &mtp, &mtp
			},
			"",
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg, existingMtp, newMtp := tc.setup()
			_, err := suite.app.PerpetualKeeper.OpenConsolidate(suite.ctx, existingMtp, newMtp, msg, ptypes.BaseCurrency)

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
	testCases := []struct {
		name            string
		setup           func() *types.MsgOpen
		expectedErrMsg  string
		consolidatedMtp *types.MTP
	}{
		{
			"Sucess: Consolidate two position with different leverage and take profit price",
			func() *types.MsgOpen {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				suite.SetPerpetualPool(1)
				_, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				amount := math.NewInt(400)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				_, err = suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				openPositionMsg.Leverage = math.LegacyNewDec(3)
				openPositionMsg.TakeProfitPrice = math.LegacyMustNewDecFromStr("0.5")

				return openPositionMsg
			},
			"",
			&types.MTP{
				Collateral:      math.NewInt(800),
				Liabilities:     math.NewInt(649),
				Custody:         math.NewInt(4000),
				TakeProfitPrice: math.LegacyMustNewDecFromStr("0.692857142857142857"),
			},
		},
		{
			"Sucess: add collateral for the existing position",
			func() *types.MsgOpen {
				suite.ResetSuite()

				firstPool := uint64(1)
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				suite.SetPerpetualPool(1)
				_, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				amount := math.NewInt(400)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				_, err = suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				// make new Positon leverage 0 to add collateral
				openPositionMsg.Leverage = math.LegacyNewDec(0)

				return openPositionMsg
			},
			"",
			&types.MTP{
				Collateral:      math.NewInt(800),
				Liabilities:     math.NewInt(406),
				Custody:         math.NewInt(2800),
				TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.setup()
			position, err := suite.app.PerpetualKeeper.Open(suite.ctx, msg)
			suite.Require().NoError(err)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				consolidateMtp, mtpErr := suite.app.PerpetualKeeper.GetMTP(suite.ctx, sdk.MustAccAddressFromBech32(msg.Creator), position.Id)
				suite.Require().NoError(mtpErr)
				suite.Require().Equal(tc.consolidatedMtp.Collateral, consolidateMtp.Collateral)
				suite.Require().Equal(tc.consolidatedMtp.Liabilities, consolidateMtp.Liabilities)
				suite.Require().Equal(tc.consolidatedMtp.Custody, consolidateMtp.Custody)
				suite.Require().Equal(tc.consolidatedMtp.TakeProfitPrice, consolidateMtp.TakeProfitPrice)
			}
		})
	}
}
