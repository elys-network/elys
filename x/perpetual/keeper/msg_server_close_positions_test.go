package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/sample"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestClosePositions() {

	// Define test cases
	testCases := []struct {
		name                   string
		setup                  func() *types.MsgClosePositions
		expectedErrMsg         string
		expectedTotalPositions int
	}{
		{
			"BaseCurrency not found in asset profile",
			func() *types.MsgClosePositions {
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
				return &types.MsgClosePositions{
					Creator: sample.AccAddress(),
				}
			},
			"",
			0,
		},
		{
			"Pool Not found for the close positions requests, just continue",
			func() *types.MsgClosePositions {
				suite.ResetSuite()
				firstPool := uint64(1)
				secondPool := uint64(2)

				suite.SetPerpetualPool(firstPool)
				suite.SetPerpetualPool(secondPool)

				amount := math.NewInt(400)

				addr := suite.AddAccounts(2, nil)
				firstPositionCreator := addr[0]
				secondPositionCreator := addr[1]

				firstOpenPositionMsg := &types.MsgOpen{
					Creator:         firstPositionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          firstPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				firstPosition, err := suite.app.PerpetualKeeper.Open(suite.ctx, firstOpenPositionMsg)
				suite.Require().NoError(err)

				secondOpenPositionMsg := &types.MsgOpen{
					Creator:         secondPositionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          secondPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				secondPosition, err := suite.app.PerpetualKeeper.Open(suite.ctx, secondOpenPositionMsg)
				suite.Require().NoError(err)

				suite.app.PerpetualKeeper.RemovePool(suite.ctx, firstPool)
				suite.app.AmmKeeper.RemovePool(suite.ctx, secondPool)

				return &types.MsgClosePositions{
					Creator: sample.AccAddress(),
					Liquidate: []types.PositionRequest{
						{
							Address: firstPositionCreator.String(),
							Id:      firstPosition.Id,
						},
						{
							Address: secondPositionCreator.String(),
							Id:      secondPosition.Id,
						},
						{
							Address: sample.AccAddress(),
							Id:      2000,
						},
					},
					StopLoss: []types.PositionRequest{
						{
							Address: firstPositionCreator.String(),
							Id:      firstPosition.Id,
						},
						{
							Address: sample.AccAddress(),
							Id:      2000,
						},
					},
					TakeProfit: []types.PositionRequest{
						{
							Address: firstPositionCreator.String(),
							Id:      firstPosition.Id,
						},
						{
							Address: sample.AccAddress(),
							Id:      2000,
						},
					},
				}
			},
			"",
			2,
		},
		{
			"Liquidate unhealthy position",
			func() *types.MsgClosePositions {
				suite.ResetSuite()
				firstPool := uint64(1)

				suite.SetPerpetualPool(firstPool)

				amount := math.NewInt(400)

				addr := suite.AddAccounts(1, nil)
				firstPositionCreator := addr[0]
				tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				firstOpenPositionMsg := &types.MsgOpen{
					Creator:         firstPositionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_LONG,
					PoolId:          firstPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				firstPosition, err := suite.app.PerpetualKeeper.Open(suite.ctx, firstOpenPositionMsg)
				suite.Require().NoError(err)

				// Increase unpaid liability to reduce the MTP health
				mtp, _ := suite.app.PerpetualKeeper.GetMTP(suite.ctx, firstPositionCreator, firstPosition.Id)
				mtp.BorrowInterestUnpaidLiability = math.NewInt(389)
				suite.app.PerpetualKeeper.SetMTP(suite.ctx, &mtp)

				return &types.MsgClosePositions{
					Creator: sample.AccAddress(),
					Liquidate: []types.PositionRequest{
						{
							Address: firstPositionCreator.String(),
							Id:      firstPosition.Id,
						},
					},
					StopLoss:   []types.PositionRequest{},
					TakeProfit: []types.PositionRequest{},
				}
			},
			"",
			0,
		},
		{
			"Close at stop Loss",
			func() *types.MsgClosePositions {
				suite.ResetSuite()
				firstPool := uint64(1)

				suite.SetPerpetualPool(firstPool)

				amount := math.NewInt(400)

				addr := suite.AddAccounts(1, nil)
				firstPositionCreator := addr[0]
				tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				firstOpenPositionMsg := &types.MsgOpen{
					Creator:         firstPositionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_LONG,
					PoolId:          firstPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyMustNewDecFromStr("2.00"),
				}

				firstPosition, err := suite.app.PerpetualKeeper.Open(suite.ctx, firstOpenPositionMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyMustNewDecFromStr("2.00"),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgClosePositions{
					Creator:   sample.AccAddress(),
					Liquidate: []types.PositionRequest{},
					StopLoss: []types.PositionRequest{
						{
							Address: firstPositionCreator.String(),
							Id:      firstPosition.Id,
						},
					},
					TakeProfit: []types.PositionRequest{},
				}
			},
			"",
			0,
		},
		{
			"Close at Take Profit Price",
			func() *types.MsgClosePositions {
				suite.ResetSuite()
				firstPool := uint64(1)

				suite.SetPerpetualPool(firstPool)

				amount := math.NewInt(400)

				addr := suite.AddAccounts(1, nil)
				firstPositionCreator := addr[0]
				tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)

				firstOpenPositionMsg := &types.MsgOpen{
					Creator:         firstPositionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_LONG,
					PoolId:          firstPool,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyMustNewDecFromStr("2.00"),
				}

				firstPosition, err := suite.app.PerpetualKeeper.Open(suite.ctx, firstOpenPositionMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     tradingAssetPrice.MulInt64(4),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgClosePositions{
					Creator:   sample.AccAddress(),
					Liquidate: []types.PositionRequest{},
					StopLoss:  []types.PositionRequest{},
					TakeProfit: []types.PositionRequest{
						{
							Address: firstPositionCreator.String(),
							Id:      firstPosition.Id,
						},
					},
				}
			},
			"",
			0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.setup()
			msgSrvr := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
			_, err := msgSrvr.ClosePositions(suite.ctx, msg)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				totalMTPs := suite.app.PerpetualKeeper.GetAllMTPs(suite.ctx)
				suite.Require().Equal(tc.expectedTotalPositions, len(totalMTPs))
			}
		})
	}
}
