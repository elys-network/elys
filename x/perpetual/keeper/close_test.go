package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/sample"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestClose() {

	// Define test cases
	testCases := []struct {
		name           string
		setup          func() *types.MsgClose
		expectedErrMsg string
		repayAmount    math.Int
	}{
		{
			"mtp not found",
			func() *types.MsgClose {
				return &types.MsgClose{
					Creator: sample.AccAddress(),
					Id:      uint64(10),
					Amount:  math.NewInt(12000),
				}
			},
			"mtp not found",
			math.NewInt(0),
		},
		{
			"asset profile not found",
			func() *types.MsgClose {
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
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(500),
				}
			},
			"asset uusdc not found",
			math.NewInt(0),
		},
		{
			"Success: Close amount greater than position size, CloseRatio > 1", // CloseRatio gets changed to 1
			func() *types.MsgClose {
				suite.ResetSuite()

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

				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(50000000000000), // same as with amount 399
				}
			},
			"",
			math.NewInt(200),
		},
		{
			"Close with price greater than open price and less than take profit price",
			func() *types.MsgClose {
				suite.ResetSuite()

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

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyMustNewDecFromStr("10.0"),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(399),
				}
			},
			"",
			math.NewInt(100), // less than at the same price
		},
		{
			"Close at take profit price",
			func() *types.MsgClose {
				suite.ResetSuite()

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

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     tradingAssetPrice.MulInt64(4),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(399),
				}
			},
			"",
			math.NewInt(50),
		},
		{
			"Close at stopLoss price",
			func() *types.MsgClose {
				suite.ResetSuite()

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
					StopLossPrice:   math.LegacyMustNewDecFromStr("2.0"),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyMustNewDecFromStr("2.0"),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(399),
				}
			},
			"",
			math.NewInt(500),
		},
		{
			"Sucess: close long position,at same price as open price",
			func() *types.MsgClose {
				suite.ResetSuite()

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
				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(399),
				}
			},
			"",
			math.NewInt(200),
		},
		{
			"Success: close short position at same price as open price",
			func() *types.MsgClose {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(900),
				}
			},
			"",
			math.NewInt(4499),
		},
		{
			"Close with too much unpaid Liability to make custody amount 0",
			func() *types.MsgClose {
				suite.ResetSuite()

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

				// Increase unpaid liability
				mtp, _ := suite.app.PerpetualKeeper.GetMTP(suite.ctx, positionCreator, position.Id)
				mtp.BorrowInterestUnpaidLiability = math.NewInt(1995)
				suite.app.PerpetualKeeper.SetMTP(suite.ctx, &mtp)
				suite.T().Log("MTP: ", mtp)

				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(399),
				}
			},
			"error handling funding fee",
			math.NewInt(0),
		},
		{
			"Close short with Not Enough liquidity",
			func() *types.MsgClose {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
					Denom:     ptypes.BaseCurrency,
					Liquidity: math.NewInt(0),
				})

				return &types.MsgClose{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Amount:  math.NewInt(900),
				}
			},
			"not enough liquidity",
			math.NewInt(0),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.setup()

			res, err := suite.app.PerpetualKeeper.Close(suite.ctx, msg)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.repayAmount, res.Amount)
			}
		})
	}
}
