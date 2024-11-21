package keeper_test

import (
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/nullify"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) createNPendingSpotOrder(n int) []types.SpotOrder {
	items := make([]types.SpotOrder, n)
	for i := range items {
		items[i].OrderId = suite.app.TradeshieldKeeper.AppendPendingSpotOrder(suite.ctx, items[i])
	}
	return items
}

func (suite *TradeshieldKeeperTestSuite) TestPendingSpotOrderGet() {
	items := suite.createNPendingSpotOrder(10)
	for _, item := range items {
		got, found := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, item.OrderId)
		item.OrderPrice.Rate = sdkmath.LegacyZeroDec()
		suite.Require().True(found)
		suite.Require().Equal(
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func (suite *TradeshieldKeeperTestSuite) TestPendingSpotOrderRemove() {
	items := suite.createNPendingSpotOrder(10)
	for _, item := range items {
		suite.app.TradeshieldKeeper.RemovePendingSpotOrder(suite.ctx, item.OrderId)
		_, found := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, item.OrderId)
		suite.Require().False(found)
	}
}

func (suite *TradeshieldKeeperTestSuite) TestPendingSpotOrderGetAll() {
	_ = suite.createNPendingSpotOrder(10)

	for _, item := range suite.app.TradeshieldKeeper.GetAllPendingSpotOrder(suite.ctx) {
		got, found := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, item.OrderId)
		item.OrderPrice.Rate = sdkmath.LegacyZeroDec()
		suite.Require().True(found)
		suite.Require().Equal(
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func (suite *TradeshieldKeeperTestSuite) TestPendingSpotOrderCount() {
	items := suite.createNPendingSpotOrder(10)
	count := uint64(len(items))
	suite.Require().Equal(count, suite.app.TradeshieldKeeper.GetPendingSpotOrderCount(suite.ctx)-1)
}

// TestExecuteStopLossOrder
func (suite *TradeshieldKeeperTestSuite) TestExecuteStopLossOrder() {
	address := suite.AddAccounts(1, nil)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
	})
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.BaseCurrency,
		Denom:       ptypes.BaseCurrency,
		Decimals:    6,
		DisplayName: "USDC",
	})
	suite.SetupCoinPrices()

	_ = suite.CreateNewAmmPool(address[0], true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, math.NewInt(100000000000).MulRaw(10), math.NewInt(100000000000).MulRaw(10))

	// Set to main storage
	suite.app.TradeshieldKeeper.AppendPendingSpotOrder(suite.ctx, types.SpotOrder{
		OwnerAddress: address[0].String(),
		OrderId:      0,
		OrderType:    types.SpotOrderType_STOPLOSS,
		OrderPrice: types.OrderPrice{
			BaseDenom:  "uusdc",
			QuoteDenom: "uatom",
			Rate:       sdkmath.LegacyNewDec(1),
		},
		OrderTargetDenom: "uatom",
		OrderAmount:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000000000)),
	})

	order, _ := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, 1)

	err := suite.app.TradeshieldKeeper.ExecuteStopLossOrder(suite.ctx, order)
	suite.Require().NoError(err)

	// Should remove from pending order list
	res := suite.app.TradeshieldKeeper.GetAllPendingSpotOrder(suite.ctx)
	suite.Assert().Equal(res, []types.SpotOrder(nil))

	// Should remove from main storage
	_, found := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, 1)
	suite.Assert().False(found)
}

// TestExecuteLimitSellOrder
func (suite *TradeshieldKeeperTestSuite) TestExecuteLimitSellOrder() {
	address := suite.AddAccounts(1, nil)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
	})
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.BaseCurrency,
		Denom:       ptypes.BaseCurrency,
		Decimals:    6,
		DisplayName: "USDC",
	})
	suite.SetupCoinPrices()

	_ = suite.CreateNewAmmPool(address[0], true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, math.NewInt(100000000000).MulRaw(10), math.NewInt(100000000000).MulRaw(10))

	// Set to main storage
	suite.app.TradeshieldKeeper.AppendPendingSpotOrder(suite.ctx, types.SpotOrder{
		OwnerAddress: address[0].String(),
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITSELL,
		OrderPrice: types.OrderPrice{
			BaseDenom:  "uusdc",
			QuoteDenom: "uatom",
			Rate:       sdkmath.LegacyNewDec(0),
		},
		OrderTargetDenom: "uatom",
		OrderAmount:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000000)),
	})

	order, _ := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, 1)

	err := suite.app.TradeshieldKeeper.ExecuteLimitSellOrder(suite.ctx, order)
	suite.Require().NoError(err)

	// Should remove from pending order list
	res := suite.app.TradeshieldKeeper.GetAllPendingSpotOrder(suite.ctx)
	suite.Assert().Equal(res, []types.SpotOrder(nil))

	// Should remove from main storage
	_, found := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, 1)
	suite.Assert().False(found)
}

// TestExecuteLimitBuyOrder
func (suite *TradeshieldKeeperTestSuite) TestExecuteLimitBuyOrder() {
	address := suite.AddAccounts(1, nil)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
	})
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.BaseCurrency,
		Denom:       ptypes.BaseCurrency,
		Decimals:    6,
		DisplayName: "USDC",
	})
	suite.SetupCoinPrices()

	_ = suite.CreateNewAmmPool(address[0], true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, math.NewInt(100000000000).MulRaw(10), math.NewInt(100000000000).MulRaw(10))

	// Set to main storage
	suite.app.TradeshieldKeeper.AppendPendingSpotOrder(suite.ctx, types.SpotOrder{
		OwnerAddress: address[0].String(),
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITBUY,
		OrderPrice: types.OrderPrice{
			BaseDenom:  "uusdc",
			QuoteDenom: "uatom",
			Rate:       sdkmath.LegacyNewDec(1),
		},
		OrderTargetDenom: "uatom",
		OrderAmount:      sdk.NewCoin("uusdc", sdkmath.NewInt(100000)),
	})

	order, _ := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, 1)

	err := suite.app.TradeshieldKeeper.ExecuteLimitBuyOrder(suite.ctx, order)
	suite.Require().NoError(err)

	// Should remove from pending order list
	res := suite.app.TradeshieldKeeper.GetAllPendingSpotOrder(suite.ctx)
	suite.Assert().Equal(res, []types.SpotOrder(nil))

	// Should remove from main storage
	_, found := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, 1)
	suite.Assert().False(found)
}

// TestExecuteMarketBuyOrder
func (suite *TradeshieldKeeperTestSuite) TestExecuteMarketBuyOrder() {
	address := suite.AddAccounts(1, nil)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
	})
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.BaseCurrency,
		Denom:       ptypes.BaseCurrency,
		Decimals:    6,
		DisplayName: "USDC",
	})
	suite.SetupCoinPrices()

	_ = suite.CreateNewAmmPool(address[0], true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, math.NewInt(100000000000).MulRaw(10), math.NewInt(100000000000).MulRaw(10))

	order := types.SpotOrder{
		OwnerAddress: address[0].String(),
		OrderId:      0,
		OrderType:    types.SpotOrderType_MARKETBUY,
		OrderPrice: types.OrderPrice{
			BaseDenom:  "uusdc",
			QuoteDenom: "uatom",
			Rate:       sdkmath.LegacyNewDec(1),
		},
		OrderTargetDenom: "uatom",
		OrderAmount:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
	}

	err := suite.app.TradeshieldKeeper.ExecuteMarketBuyOrder(suite.ctx, order)
	suite.Require().NoError(err)

	// Should remove from pending order list
	res := suite.app.TradeshieldKeeper.GetAllPendingSpotOrder(suite.ctx)
	suite.Assert().Equal(res, []types.SpotOrder(nil))

	// Should remove from main storage
	_, found := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, 1)
	suite.Assert().False(found)
}
