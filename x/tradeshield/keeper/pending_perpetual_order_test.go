package keeper_test

import (
	"fmt"

	perpetualtypes "github.com/elys-network/elys/v7/x/perpetual/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) createNPendingPerpetualOrder(n int) []types.PerpetualOrder {
	items := make([]types.PerpetualOrder, n)
	for i := range items {
		items[i] = types.PerpetualOrder{
			OrderId:            0,
			OwnerAddress:       sdk.AccAddress([]byte(fmt.Sprintf("address%d", i))).String(),
			PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
			Position:           types.PerpetualPosition_LONG,
			TriggerPrice:       math.LegacyNewDec(1),
			Collateral:         sdk.Coin{Denom: "denom", Amount: math.NewInt(10)},
			Leverage:           math.LegacyNewDec(int64(i)),
			TakeProfitPrice:    math.LegacyNewDec(1),
			PositionId:         uint64(i),
			Status:             types.Status_PENDING,
			StopLossPrice:      math.LegacyNewDec(1),
		}
		items[i].OrderId = suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, items[i])
	}
	return items
}

func (suite *TradeshieldKeeperTestSuite) TestPendingPerpetualOrderGet() {
	items := suite.createNPendingPerpetualOrder(10)
	for _, item := range items {
		got, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(item.OwnerAddress), item.PoolId, item.OrderId)
		suite.Require().True(found)
		suite.Require().Equal(
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func (suite *TradeshieldKeeperTestSuite) TestPendingPerpetualOrderRemove() {
	items := suite.createNPendingPerpetualOrder(10)
	for _, item := range items {
		suite.app.TradeshieldKeeper.RemovePendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(item.OwnerAddress), item.PoolId, item.OrderId)
		_, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(item.OwnerAddress), item.PoolId, item.OrderId)
		suite.Require().False(found)
	}
}

func (suite *TradeshieldKeeperTestSuite) TestPendingPerpetualOrderGetAll() {
	items := suite.createNPendingPerpetualOrder(10)
	suite.Require().ElementsMatch(
		nullify.Fill(items),
		nullify.Fill(suite.app.TradeshieldKeeper.GetAllPendingPerpetualOrder(suite.ctx)),
	)
}

func (suite *TradeshieldKeeperTestSuite) TestPendingPerpetualOrderCount() {
	items := suite.createNPendingPerpetualOrder(10)
	count := uint64(len(items))
	suite.Require().Equal(count, suite.app.TradeshieldKeeper.GetPendingPerpetualOrderCount(suite.ctx)-1)
}

func (suite *TradeshieldKeeperTestSuite) TestExecuteLimitOpenOrder() {

	address := suite.AddAccounts(3, nil)
	_, _, _ = suite.SetPerpetualPool(1)

	perpetualOrder := types.PerpetualOrder{
		OwnerAddress:       address[2].String(),
		OrderId:            1,
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice:       math.LegacyMustNewDecFromStr("10"),
		Position:           types.PerpetualPosition_LONG,
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(10),
		TakeProfitPrice:    math.LegacyNewDec(50),
		StopLossPrice:      math.LegacyZeroDec(),
		PoolId:             1,
	}

	// Fund orderAddress
	orderAddress := perpetualOrder.GetOrderAddress()
	suite.AddAccounts(1, []sdk.AccAddress{orderAddress})

	suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, perpetualOrder)

	order, _ := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(perpetualOrder.OwnerAddress), perpetualOrder.PoolId, perpetualOrder.OrderId)

	err := suite.app.TradeshieldKeeper.ExecuteLimitOpenOrder(suite.ctx, order)
	suite.Require().NoError(err)

	_, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(perpetualOrder.OwnerAddress), perpetualOrder.PoolId, perpetualOrder.OrderId)
	suite.Require().False(found)
}

// TODO: Update it when close is supported from tradeshield
func (suite *TradeshieldKeeperTestSuite) TestExecuteLimitCloseOrder() {
	address := suite.AddAccounts(3, nil)
	_, _, _ = suite.SetPerpetualPool(1)

	perpetualOrder := types.PerpetualOrder{
		OwnerAddress:       address[2].String(),
		OrderId:            1,
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice:       math.LegacyMustNewDecFromStr("10"),
		Position:           types.PerpetualPosition_LONG,
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(10),
		TakeProfitPrice:    math.LegacyNewDec(50),
		StopLossPrice:      math.LegacyZeroDec(),
		PoolId:             1,
	}

	// Fund orderAddress
	orderAddress := perpetualOrder.GetOrderAddress()
	suite.AddAccounts(1, []sdk.AccAddress{orderAddress})

	suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, perpetualOrder)
	order, _ := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(perpetualOrder.OwnerAddress), perpetualOrder.PoolId, perpetualOrder.OrderId)
	err := suite.app.TradeshieldKeeper.ExecuteLimitOpenOrder(suite.ctx, order)
	suite.Require().NoError(err)

	perpetualOrder.PerpetualOrderType = types.PerpetualOrderType_LIMITCLOSE
	perpetualOrder.TriggerPrice = math.LegacyZeroDec()
	perpetualOrder.PositionId = 1
	orderId := suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, perpetualOrder)

	order, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(perpetualOrder.OwnerAddress), perpetualOrder.PoolId, orderId)
	suite.Require().True(found)

	err = suite.app.TradeshieldKeeper.ExecuteLimitCloseOrder(suite.ctx, order)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), perpetualtypes.ErrInvalidAmount.Error())

	_, found = suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(perpetualOrder.OwnerAddress), perpetualOrder.PoolId, orderId)
	suite.Require().True(found)
}

func (suite *TradeshieldKeeperTestSuite) TestExecuteMarketOpenOrder() {
	address := suite.AddAccounts(3, nil)
	_, _, _ = suite.SetPerpetualPool(1)

	suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, types.PerpetualOrder{
		OwnerAddress:       address[2].String(),
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice:       math.LegacyMustNewDecFromStr("10"),
		Position:           types.PerpetualPosition_LONG,
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(10),
		TakeProfitPrice:    math.LegacyNewDec(10),
		StopLossPrice:      math.LegacyZeroDec(),
		PoolId:             1,
	})

	order, _ := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(address[2].String()), 1, 1)

	err := suite.app.TradeshieldKeeper.ExecuteMarketOpenOrder(suite.ctx, order)
	suite.Require().NoError(err)

	_, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(address[2].String()), 1, 1)
	suite.Require().False(found)
}

// TODO: Update tests when close is supported from tradeshield
func (suite *TradeshieldKeeperTestSuite) TestExecuteMarketCloseOrder() {
	address := suite.AddAccounts(3, nil)
	_, _, _ = suite.SetPerpetualPool(1)

	perpetualOrder := types.PerpetualOrder{
		OwnerAddress:       address[2].String(),
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice:       math.LegacyMustNewDecFromStr("10"),
		Position:           types.PerpetualPosition_LONG,
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(10),
		TakeProfitPrice:    math.LegacyNewDec(10),
		StopLossPrice:      math.LegacyZeroDec(),
		PoolId:             1,
	}
	suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, perpetualOrder)
	order, _ := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(address[2].String()), 1, 1)
	err := suite.app.TradeshieldKeeper.ExecuteMarketOpenOrder(suite.ctx, order)
	suite.Require().NoError(err)

	perpetualOrder.PerpetualOrderType = types.PerpetualOrderType_LIMITCLOSE
	perpetualOrder.TriggerPrice = math.LegacyZeroDec()
	perpetualOrder.PositionId = 1
	orderId := suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, perpetualOrder)

	order, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(address[2].String()), 1, orderId)
	suite.Require().True(found)

	err = suite.app.TradeshieldKeeper.ExecuteMarketCloseOrder(suite.ctx, order)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), perpetualtypes.ErrInvalidAmount.Error())

	_, found = suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(address[2].String()), 1, orderId)
	suite.Require().True(found)
}
