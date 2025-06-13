package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"sort"

	"github.com/elys-network/elys/v6/x/oracle/types"
)

// This test ensures old price data is automatically removed
func (suite *KeeperTestSuite) TestEndBlock() {
	btcPriceData := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(23),
			Source:    "elys",
			Timestamp: 20,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(12),
			Source:    "band",
			Timestamp: 10,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(76),
			Source:    "elys",
			Timestamp: 40,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(55),
			Source:    "band",
			Timestamp: 30,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(89),
			Source:    "random",
			Timestamp: 25,
		},
	}
	ethPriceData := []types.Price{
		{
			Asset:     "ETH",
			Price:     sdkmath.LegacyNewDec(23),
			Source:    "elys",
			Timestamp: 34,
		},
		{
			Asset:     "ETH",
			Price:     sdkmath.LegacyNewDec(12),
			Source:    "band",
			Timestamp: 23,
		},
		{
			Asset:     "ETH",
			Price:     sdkmath.LegacyNewDec(76),
			Source:    "elys",
			Timestamp: 89,
		},
		{
			Asset:     "ETH",
			Price:     sdkmath.LegacyNewDec(55),
			Source:    "random",
			Timestamp: 54,
		},
		{
			Asset:     "ETH",
			Price:     sdkmath.LegacyNewDec(89),
			Source:    "band",
			Timestamp: 76,
		},
	}
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, types.AssetInfo{
		Denom:      "satoshi",
		Display:    "BTC",
		BandTicker: "BTC",
		ElysTicker: "BTC",
		Decimal:    8,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, types.AssetInfo{
		Denom:      "wei",
		Display:    "ETH",
		BandTicker: "ETH",
		ElysTicker: "ETH",
		Decimal:    18,
	})
	for _, price := range btcPriceData {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}
	for _, price := range ethPriceData {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}

	suite.app.OracleKeeper.EndBlock(suite.ctx)

	sort.Slice(btcPriceData, func(i, j int) bool {
		return btcPriceData[i].Timestamp < btcPriceData[j].Timestamp
	})
	sort.Slice(ethPriceData, func(i, j int) bool {
		return ethPriceData[i].Timestamp < ethPriceData[j].Timestamp
	})

	allPrices := suite.app.OracleKeeper.GetAllAssetPrice(suite.ctx, "BTC")
	suite.Require().Equal(1, len(allPrices))
	suite.Require().Equal(btcPriceData[len(btcPriceData)-1], allPrices[0])

	allPrices = suite.app.OracleKeeper.GetAllAssetPrice(suite.ctx, "ETH")
	suite.Require().Equal(1, len(allPrices))
	suite.Require().Equal(ethPriceData[len(ethPriceData)-1], allPrices[0])
}
