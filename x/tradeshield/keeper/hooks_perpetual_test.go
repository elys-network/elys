package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) TestPerpetualHooks_AfterPerpetualPositionClosed() {
	// Setup accounts and pool
	addresses := suite.AddAccounts(2, nil)
	owner := addresses[0]
	poolId := uint64(1)
	positionId := uint64(100)

	// Create pending perpetual orders that should be deleted
	limitCloseOrder := types.PerpetualOrder{
		OrderId:            uint64(1),
		OwnerAddress:       owner.String(),
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		Position:           types.PerpetualPosition_LONG,
		TriggerPrice:       math.LegacyNewDec(1),
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(1),
		TakeProfitPrice:    math.LegacyNewDec(1),
		PositionId:         positionId,
		Status:             types.Status_PENDING,
		StopLossPrice:      math.LegacyNewDec(1),
		PoolId:             poolId,
	}
	suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, limitCloseOrder)

	stopLossOrder := types.PerpetualOrder{
		OrderId:            uint64(2),
		OwnerAddress:       owner.String(),
		PerpetualOrderType: types.PerpetualOrderType_STOPLOSSPERP,
		Position:           types.PerpetualPosition_LONG,
		TriggerPrice:       math.LegacyNewDec(1),
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(1),
		TakeProfitPrice:    math.LegacyNewDec(1),
		PositionId:         positionId,
		Status:             types.Status_PENDING,
		StopLossPrice:      math.LegacyNewDec(1),
		PoolId:             poolId,
	}
	suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, stopLossOrder)

	// Create an order that should NOT be deleted (different position ID)
	otherPositionOrder := types.PerpetualOrder{
		OrderId:            uint64(3),
		OwnerAddress:       owner.String(),
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		Position:           types.PerpetualPosition_LONG,
		TriggerPrice:       math.LegacyNewDec(1),
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(1),
		TakeProfitPrice:    math.LegacyNewDec(1),
		PositionId:         positionId + 1, // Different position ID
		Status:             types.Status_PENDING,
		StopLossPrice:      math.LegacyNewDec(1),
		PoolId:             poolId,
	}
	suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, otherPositionOrder)

	// Create an order that should NOT be deleted (different position ID)
	openPositionOrder := types.PerpetualOrder{
		OrderId:            uint64(3),
		OwnerAddress:       owner.String(),
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		Position:           types.PerpetualPosition_LONG,
		TriggerPrice:       math.LegacyNewDec(1),
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(1),
		TakeProfitPrice:    math.LegacyNewDec(1),
		PositionId:         positionId, // Different position ID
		Status:             types.Status_PENDING,
		StopLossPrice:      math.LegacyNewDec(1),
		PoolId:             poolId,
	}
	suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, openPositionOrder)

	perpetualPool, _, ammPool := suite.SetPerpetualPool(poolId)

	// Execute the hook
	perpetualHooks := suite.app.TradeshieldKeeper.PerpetualHooks()
	err := perpetualHooks.AfterPerpetualPositionClosed(suite.ctx, ammPool, perpetualPool, owner, math.LegacyOneDec(), positionId)

	// Verify the result
	suite.Require().NoError(err)

	// Verify that orders with matching position ID and correct types are deleted
	status := types.Status_ALL
	orders, _, err := suite.app.TradeshieldKeeper.GetPendingPerpetualOrdersForAddress(suite.ctx, owner.String(), &status, nil)
	suite.Require().NoError(err)

	// Should have 2 orders remaining (other position and open position)
	suite.Require().Len(orders, 2)

	// Verify that the remaining order is the one with different position ID
	suite.Require().Equal(positionId+1, orders[0].PositionId)
	suite.Require().Equal(positionId, orders[1].PositionId)
}
