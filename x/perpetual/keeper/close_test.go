package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/testutil/sample"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
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
				tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
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
			"asset info uusdc not found",
			math.NewInt(0),
		},
		{
			"Success: Close amount greater than position size, CloseRatio > 1", // CloseRatio gets changed to 1
			func() *types.MsgClose {
				suite.ResetSuite()

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
			math.NewInt(204),
		},
		{
			"Close with price greater than open price and less than take profit price",
			func() *types.MsgClose {
				suite.ResetSuite()

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
					Amount:  math.NewInt(107),
				}
			},
			"",
			math.NewInt(31), // less than at the same price
		},
		{
			"Close at take profit price",
			func() *types.MsgClose {
				suite.ResetSuite()

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
					Amount:  math.NewInt(699),
				}
			},
			"",
			math.NewInt(91),
		},
		{
			"Close at stopLoss price",
			func() *types.MsgClose {
				suite.ResetSuite()

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
					Amount:  math.NewInt(699),
				}
			},
			"",
			math.NewInt(502),
		},
		{
			"Success: close long position,at same price as open price",
			func() *types.MsgClose {
				suite.ResetSuite()

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
					Amount:  math.NewInt(699),
				}
			},
			"",
			math.NewInt(204),
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
			math.NewInt(4498),
		},
		// TODO: Edge case when custody becomes low, this is throwing error, instead it should be closed
		// FIX this: error updating mtp health: unable to swap (EstimateSwapGivenOut) for out 1uatom and in denom uusdc: amount too low
		// {
		// 	"Force Close with too much unpaid Liability making custody amount 0",
		// 	func() *types.MsgClose {
		// 		suite.ResetSuite()

		// 		addr := suite.AddAccounts(1, nil)
		// 		positionCreator := addr[0]
		// 		_, _, ammPool := suite.SetPerpetualPool(1)
		// 		tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
		// 		suite.Require().NoError(err)
		// 		openPositionMsg := &types.MsgOpen{
		// 			Creator:         positionCreator.String(),
		// 			Leverage:        math.LegacyNewDec(2),
		// 			Position:        types.Position_LONG,
		// 			PoolId:          ammPool.PoolId,
		// 			TradingAsset:    ptypes.ATOM,
		// 			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
		// 			TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		// 			StopLossPrice:   math.LegacyZeroDec(),
		// 		}
		// 		position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
		// 		suite.Require().NoError(err)

		// 		// Increase unpaid liability
		// 		mtp, _ := suite.app.PerpetualKeeper.GetMTP(suite.ctx, positionCreator, position.Id)
		// 		mtp.BorrowInterestUnpaidLiability = math.NewInt(1995)
		// 		suite.app.PerpetualKeeper.SetMTP(suite.ctx, &mtp)
		// 		suite.T().Log("MTP: ", mtp)

		// 		return &types.MsgClose{
		// 			Creator: positionCreator.String(),
		// 			Id:      position.Id,
		// 			Amount:  math.NewInt(900),
		// 		}
		// 	},
		// 	"",
		// 	math.NewInt(206),
		// },
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
				suite.Require().Equal(res.Amount.String(), tc.repayAmount.String())
			}
		})
	}
}
