package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/testutil/sample"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/keeper"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestAddCollateral() {

	var initialPoolBankBalance sdk.Coins
	var initialAccountedPoolBalance sdk.Coins

	var finalPoolBankBalance sdk.Coins
	var finalAccountedPoolBalance sdk.Coins

	var ammPool ammtypes.Pool
	// Define test cases
	testCases := []struct {
		name           string
		setup          func() *types.MsgAddCollateral
		expectedErrMsg string
		postValidate   func(msg *types.MsgAddCollateral)
	}{
		{
			"invalid address",
			func() *types.MsgAddCollateral {
				return &types.MsgAddCollateral{
					Creator:       "invalid",
					Id:            uint64(10),
					AddCollateral: sdk.NewCoin("uusdc", math.NewInt(12000)),
					PoolId:        1,
				}
			},
			"decoding bech32 failed: invalid bech32 string length 7",
			func(msg *types.MsgAddCollateral) {

			},
		},
		{
			"mtp not found",
			func() *types.MsgAddCollateral {
				return &types.MsgAddCollateral{
					Creator:       sample.AccAddress(),
					Id:            uint64(10),
					AddCollateral: sdk.NewCoin("uusdc", math.NewInt(12000)),
					PoolId:        1,
				}
			},
			"mtp not found",
			func(msg *types.MsgAddCollateral) {

			},
		},
		{
			"asset profile not found",
			func() *types.MsgAddCollateral {
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool = suite.SetPerpetualPool(1)
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
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
				return &types.MsgAddCollateral{
					Creator:       positionCreator.String(),
					Id:            position.Id,
					AddCollateral: sdk.NewCoin("uusdc", math.NewInt(12000)),
					PoolId:        ammPool.PoolId,
				}
			},
			"asset uusdc not found",
			func(msg *types.MsgAddCollateral) {

			},
		},
		{
			"invalid collateral denom",
			func() *types.MsgAddCollateral {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool = suite.SetPerpetualPool(1)
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

				return &types.MsgAddCollateral{
					Creator:       positionCreator.String(),
					Id:            position.Id,
					AddCollateral: sdk.NewCoin(ptypes.ATOM, openPositionMsg.Collateral.Amount.QuoRaw(2)),
					PoolId:        ammPool.PoolId,
				}
			},
			"denom not same as collateral asset",
			func(msg *types.MsgAddCollateral) {

			},
		},
		{
			"Success: adding collateral LONG and collateral usdc",
			func() *types.MsgAddCollateral {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool = suite.SetPerpetualPool(1)
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

				initialPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				initialAccountedPoolBalance = accountedPool.TotalTokens

				return &types.MsgAddCollateral{
					Creator:       positionCreator.String(),
					Id:            position.Id,
					AddCollateral: sdk.NewCoin(ptypes.BaseCurrency, openPositionMsg.Collateral.Amount.QuoRaw(2)),
					PoolId:        ammPool.PoolId,
				}
			},
			"",
			func(msg *types.MsgAddCollateral) {
				finalPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				finalAccountedPoolBalance = accountedPool.TotalTokens

				suite.Require().Equal(initialAccountedPoolBalance, finalAccountedPoolBalance)
				suite.Require().Equal(initialPoolBankBalance.Add(msg.AddCollateral), finalPoolBankBalance)
			},
		},
		{
			"Success: adding collateral LONG and collateral atom",
			func() *types.MsgAddCollateral {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool = suite.SetPerpetualPool(1)
				tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					Collateral:      sdk.NewCoin(ptypes.ATOM, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				initialPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				initialAccountedPoolBalance = accountedPool.TotalTokens

				return &types.MsgAddCollateral{
					Creator:       positionCreator.String(),
					Id:            position.Id,
					AddCollateral: sdk.NewCoin(ptypes.ATOM, openPositionMsg.Collateral.Amount.QuoRaw(2)),
					PoolId:        ammPool.PoolId,
				}
			},
			"",
			func(msg *types.MsgAddCollateral) {
				finalPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				finalAccountedPoolBalance = accountedPool.TotalTokens

				// This is position expansion
				suite.Require().Equal(initialPoolBankBalance.Add(msg.AddCollateral), finalPoolBankBalance)
				atleastExpected := initialAccountedPoolBalance.Add(msg.AddCollateral)
				suite.Require().True(finalAccountedPoolBalance.AmountOf(ptypes.ATOM).LTE(initialAccountedPoolBalance.AmountOf(ptypes.ATOM)))
				suite.Require().True(finalAccountedPoolBalance.AmountOf(ptypes.BaseCurrency).GTE(atleastExpected.AmountOf(ptypes.BaseCurrency)))
			},
		},
		{
			"Success: adding collateral SHORT",
			func() *types.MsgAddCollateral {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool = suite.SetPerpetualPool(1)
				tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_SHORT,
					PoolId:          ammPool.PoolId,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.QuoInt64(4),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				initialPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				initialAccountedPoolBalance = accountedPool.TotalTokens

				return &types.MsgAddCollateral{
					Creator:       positionCreator.String(),
					Id:            position.Id,
					AddCollateral: sdk.NewCoin(ptypes.BaseCurrency, openPositionMsg.Collateral.Amount.QuoRaw(2)),
					PoolId:        ammPool.PoolId,
				}
			},
			"",
			func(msg *types.MsgAddCollateral) {
				finalPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				finalAccountedPoolBalance = accountedPool.TotalTokens

				// This is position expansion
				suite.Require().Equal(initialPoolBankBalance.Add(msg.AddCollateral), finalPoolBankBalance)

				suite.Require().True(finalAccountedPoolBalance.AmountOf(ptypes.ATOM).GTE(initialAccountedPoolBalance.AmountOf(ptypes.ATOM)))
				suite.Require().True(finalAccountedPoolBalance.AmountOf(ptypes.BaseCurrency).LTE(initialAccountedPoolBalance.AmountOf(ptypes.BaseCurrency)))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.setup()
			server := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
			_, err := server.AddCollateral(suite.ctx, msg)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postValidate(msg)
		})
	}
}
