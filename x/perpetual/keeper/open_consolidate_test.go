package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestOpenConsolidate() {
	testCases := []struct {
		name            string
		setup           func() (*types.MsgOpen, *types.MTP, *types.MTP)
		expectedErrMsg  string
		consolidatedMtp *types.MTP
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

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, positionCreator, position.Id)
				suite.Require().NoError(err)

				suite.app.AmmKeeper.RemovePool(suite.ctx, firstPool)

				return openPositionMsg, &mtp, &mtp
			},
			"perpetual pool does not exist",
			nil,
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
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, positionCreator, position.Id)
				suite.Require().NoError(err)

				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.SafetyFactor = math.LegacyMustNewDecFromStr("1.30")
				suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)

				return openPositionMsg, &mtp, &mtp
			},
			"mtp health would be too low for safety factor",
			nil,
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
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
				suite.Require().NoError(err)
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, positionCreator, position.Id)
				suite.Require().NoError(err)

				return openPositionMsg, &mtp, &mtp
			},
			"",
			nil,
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
