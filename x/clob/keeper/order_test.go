package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (suite *KeeperTestSuite) TestBuyOrderBook() {
	order1 := types.PerpetualOrder{
		OrderId: types.OrderId{
			MarketId:  MarketId,
			OrderType: types.OrderType_ORDER_TYPE_LIMIT_BUY,
			PriceTick: types.PriceTick(math.LegacyMustNewDecFromStr("1.0").MulInt64(types.PriceMultiplier).TruncateInt64()),
			Counter:   1,
		},
		Owner:        authtypes.NewModuleAddress("1").String(),
		SubAccountId: MarketId,
		Amount:       math.LegacyNewDec(100),
		Filled:       math.LegacyZeroDec(),
	}

	order2 := order1
	order2.OrderId.PriceTick = types.PriceTick(math.LegacyMustNewDecFromStr("1.5").MulInt64(types.PriceMultiplier).TruncateInt64())
	order2.OrderId.Counter = 2
	order2.Owner = authtypes.NewModuleAddress("2").String()
	order2.Amount = math.LegacyNewDec(50)

	order3 := order1
	order3.OrderId.PriceTick = types.PriceTick(math.LegacyMustNewDecFromStr("1.5").MulInt64(types.PriceMultiplier).TruncateInt64())
	order3.OrderId.Counter = 3
	order3.Owner = authtypes.NewModuleAddress("3").String()
	order3.Amount = math.LegacyNewDec(75)

	order4 := order1
	order4.OrderId.PriceTick = types.PriceTick(math.LegacyMustNewDecFromStr("2").MulInt64(types.PriceMultiplier).TruncateInt64())
	order4.OrderId.Counter = 4
	order4.Owner = authtypes.NewModuleAddress("4").String()
	order4.Amount = math.LegacyNewDec(90)

	suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, order2)
	suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, order4)
	suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, order1)
	suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, order3)

	expectedList := []types.PerpetualOrder{order4, order2, order3, order1}

	iterator := suite.app.ClobKeeper.GetBuyOrderIterator(suite.ctx, MarketId)
	defer iterator.Close()

	index := 0
	for ; iterator.Valid(); iterator.Next() {
		expectedOrderId := types.NewOrderId(expectedList[index].GetMarketId(), expectedList[index].GetOrderType(), expectedList[index].GetPriceTick(), expectedList[index].GetCounter())

		prefix := append(sdk.Uint64ToBigEndian(MarketId), []byte("/")...)
		prefix = append(prefix, types.TrueByte)
		prefix = append(prefix, []byte("/")...)
		iteratorKeyWithMarketId := append(prefix, iterator.Key()...)

		suite.Equal(expectedOrderId.KeyWithoutPrefix(), iteratorKeyWithMarketId)

		index++
	}
}

func (suite *KeeperTestSuite) TestSellOrderBook() {
	order1 := types.PerpetualOrder{
		OrderId: types.OrderId{
			MarketId:  MarketId,
			OrderType: types.OrderType_ORDER_TYPE_LIMIT_SELL,
			PriceTick: types.PriceTick(math.LegacyMustNewDecFromStr("1.0").MulInt64(types.PriceMultiplier).TruncateInt64()),
			Counter:   1,
		},
		Owner:        authtypes.NewModuleAddress("1").String(),
		SubAccountId: MarketId,
		Amount:       math.LegacyNewDec(100),
		Filled:       math.LegacyZeroDec(),
	}

	order2 := order1
	order2.OrderId.PriceTick = types.PriceTick(math.LegacyMustNewDecFromStr("1.5").MulInt64(types.PriceMultiplier).TruncateInt64())
	order2.OrderId.Counter = 2
	order2.Owner = authtypes.NewModuleAddress("2").String()
	order2.Amount = math.LegacyNewDec(50)

	order3 := order1
	order3.OrderId.PriceTick = types.PriceTick(math.LegacyMustNewDecFromStr("1.5").MulInt64(types.PriceMultiplier).TruncateInt64())
	order3.OrderId.Counter = 3
	order3.Owner = authtypes.NewModuleAddress("3").String()
	order3.Amount = math.LegacyNewDec(75)

	order4 := order1
	order4.OrderId.PriceTick = types.PriceTick(math.LegacyMustNewDecFromStr("2").MulInt64(types.PriceMultiplier).TruncateInt64())
	order4.OrderId.Counter = 4
	order4.Owner = authtypes.NewModuleAddress("4").String()
	order4.Amount = math.LegacyNewDec(90)

	suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, order2)
	suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, order4)
	suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, order1)
	suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, order3)

	expectedList := []types.PerpetualOrder{order1, order2, order3, order4}

	iterator := suite.app.ClobKeeper.GetSellOrderIterator(suite.ctx, MarketId)
	defer iterator.Close()

	index := 0
	for ; iterator.Valid(); iterator.Next() {
		expectedOrderId := types.NewOrderId(expectedList[index].GetMarketId(), expectedList[index].GetOrderType(), expectedList[index].GetPriceTick(), expectedList[index].GetCounter())

		prefix := append(sdk.Uint64ToBigEndian(MarketId), []byte("/")...)
		prefix = append(prefix, types.FalseByte)
		prefix = append(prefix, []byte("/")...)
		iteratorKeyWithMarketId := append(prefix, iterator.Key()...)

		suite.Equal(expectedOrderId.KeyWithoutPrefix(), iteratorKeyWithMarketId)

		index++
	}
}

func (suite *KeeperTestSuite) TestRequiredBalanceForOrder() {
	suite.SetupTest()
	suite.CreateMarketWithFees(BaseDenom)
	order1 := types.PerpetualOrder{
		OrderId: types.OrderId{
			MarketId:  MarketId,
			OrderType: types.OrderType_ORDER_TYPE_LIMIT_BUY,
			PriceTick: types.PriceTick(math.LegacyMustNewDecFromStr("10.5").MulInt64(types.PriceMultiplier).TruncateInt64()),
			Counter:   1,
		},
		Owner:        authtypes.NewModuleAddress("1").String(),
		SubAccountId: MarketId,
		Amount:       math.LegacyNewDec(100),
		Filled:       math.LegacyZeroDec(),
	}

	market, err := suite.app.ClobKeeper.GetPerpetualMarket(suite.ctx, MarketId)
	suite.Require().NoError(err)

	got, err := suite.app.ClobKeeper.RequiredBalanceForOrder(suite.ctx, order1)
	suite.Require().NoError(err)

	value := order1.GetPrice().Mul(order1.Amount)
	// Assuming taker fees >= maker fees
	fees := math.LegacyMaxDec(market.TakerFeeRate, market.MakerFeeRate).Mul(value)
	margin := market.InitialMarginRatio.Mul(value)
	suite.Require().Equal(got, sdk.NewCoin(market.QuoteDenom, fees.Add(margin).RoundInt()))
}
