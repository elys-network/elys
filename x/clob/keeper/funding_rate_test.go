package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/clob/types"
)

func (suite *KeeperTestSuite) TestUpdateFundingRate() {
	suite.ResetSuite()

	baseDenom := "uatom"
	markets := suite.CreateMarket(baseDenom)
	market := markets[0]

	p1 := []types.Trade{
		{
			MarketId: 1,
			Price:    math.LegacyNewDecWithPrec(101, 1),
			Quantity: math.NewInt(100),
		},
		{
			MarketId: 1,
			Price:    math.LegacyNewDecWithPrec(103, 1),
			Quantity: math.NewInt(300),
		},
	}
	p2 := []types.Trade{
		{
			MarketId: 1,
			Price:    math.LegacyNewDecWithPrec(105, 1),
			Quantity: math.NewInt(200)},
		{
			MarketId: 1,
			Price:    math.LegacyNewDecWithPrec(110, 1),
			Quantity: math.NewInt(300)},
	}
	suite.IncreaseHeight(1)
	for _, p := range p1 {
		suite.app.ClobKeeper.SetTwapPrices(suite.ctx, p) // cumulativePrice = 0
	}
	suite.IncreaseHeight(1)
	for _, p := range p2 {
		suite.app.ClobKeeper.SetTwapPrices(suite.ctx, p) // cumulativePrice = 10.25*5 = 51.25
	}
	twapPrice := math.LegacyMustNewDecFromStr("10.25")
	suite.Require().Equal(twapPrice, suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, market.Id))

	testCases := []struct {
		name           string
		expectedErrMsg string
		result         math.LegacyDec
		pre            func()
		post           func()
	}{
		{
			"oracle price not found",
			"asset price not found for denom (uatom)",
			math.LegacyDec{},
			func() {
				suite.app.OracleKeeper.RemovePrice(suite.ctx, "ATOM", "test", uint64(suite.ctx.BlockTime().Unix()-15))
			},
			func() {
			},
		},
		{
			"funding rate increases within limit",
			"",
			math.LegacyMustNewDecFromStr("0.005025125628140704"),
			func() {
				oraclePrice := twapPrice.Sub(twapPrice.Mul(market.MaxFundingRateChange.Quo(math.LegacyNewDec(2)))) //twap price(1 - MaxFundingRateChange/2) = 10.19875
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{oraclePrice, math.LegacyNewDec(1)})
			},
			func() {
			},
		},
		{
			"funding rate decreases within limit",
			"",
			math.LegacyMustNewDecFromStr("-0.002493765586034913"),
			func() {
				oraclePrice := twapPrice.Add(twapPrice.Mul(market.MaxFundingRateChange.Quo(math.LegacyNewDec(4)))) //twap price(1 + MaxFundingRateChange/4) = 10.275625
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{oraclePrice, math.LegacyNewDec(1)})
			},
			func() {
			},
		},
		{
			"funding rate increases exceeding limit",
			"",
			math.LegacyMustNewDecFromStr("-0.002493765586034913").Add(market.MaxFundingRateChange),
			func() {
				oraclePrice := twapPrice.Sub(twapPrice.Mul(market.MaxFundingRateChange.MulInt64(2))) //twap price(1 - 2*MaxFundingRateChange)
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{oraclePrice, math.LegacyNewDec(1)})
			},
			func() {
			},
		},
		{
			"funding rate decreases exceeding limit",
			"",
			math.LegacyMustNewDecFromStr("-0.002493765586034913"),
			func() {
				oraclePrice := twapPrice.Add(twapPrice.Mul(market.MaxFundingRateChange.MulInt64(2))) //twap price(1 - 2*MaxFundingRateChange)
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{oraclePrice, math.LegacyNewDec(1)})
			},
			func() {
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.IncreaseHeight(1)
			tc.pre()
			err := suite.app.ClobKeeper.UpdateFundingRate(suite.ctx, market)
			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				fundingRate := suite.app.ClobKeeper.GetFundingRate(suite.ctx, market.Id)
				suite.Require().Equal(tc.result, fundingRate.Rate)
				suite.Require().Equal(uint64(suite.ctx.BlockHeight()), fundingRate.Block)
			}
			tc.post()
		})
	}

	markets = suite.CreateMarket("uosmo")
	suite.SetPrice([]string{"OSMO"}, []math.LegacyDec{math.LegacyNewDec(2)})
	err := suite.app.ClobKeeper.UpdateFundingRate(suite.ctx, markets[0])
	suite.Require().NoError(err)

	all := suite.app.ClobKeeper.GetAllFundingRate(suite.ctx)
	suite.Require().Len(all, 2)

	suite.Require().Equal(math.LegacyMustNewDecFromStr("-0.002493765586034913"), all[0].Rate)
	suite.Require().Equal(uint64(suite.ctx.BlockHeight()), all[0].Block)
	suite.Require().Equal(uint64(1), all[0].MarketId)

	suite.Require().Equal(math.LegacyMustNewDecFromStr("-0.01"), all[1].Rate)
	suite.Require().Equal(uint64(suite.ctx.BlockHeight()), all[1].Block)
	suite.Require().Equal(uint64(2), all[1].MarketId)

}
