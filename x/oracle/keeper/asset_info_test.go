package keeper_test

import (
	"github.com/elys-network/elys/v6/x/oracle/types"
)

func (suite *KeeperTestSuite) TestAssetInfoSetGetDelete() {
	assetInfos := []types.AssetInfo{
		{
			Denom:      "satoshi",
			Display:    "BTC",
			BandTicker: "BTC",
			ElysTicker: "",
		},
		{
			Denom:      "wei",
			Display:    "ETH",
			BandTicker: "ETH",
			ElysTicker: "",
		},
	}

	for _, info := range assetInfos {
		suite.app.OracleKeeper.SetAssetInfo(suite.ctx, info)
	}

	for _, info := range assetInfos {
		b, found := suite.app.OracleKeeper.GetAssetInfo(suite.ctx, info.Denom)
		suite.Require().True(found)
		suite.Require().Equal(b, info)

	}

	allAssetInfos := suite.app.OracleKeeper.GetAllAssetInfo(suite.ctx)
	suite.Require().Len(allAssetInfos, 2)

	suite.app.OracleKeeper.RemoveAssetInfo(suite.ctx, assetInfos[0].Denom)

	allAssetInfos = suite.app.OracleKeeper.GetAllAssetInfo(suite.ctx)
	suite.Require().Len(allAssetInfos, 1)

	_, found := suite.app.OracleKeeper.GetAssetInfo(suite.ctx, assetInfos[0].Denom)
	suite.Require().False(found)
}
