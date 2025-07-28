package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestSetAmmPool() {
	p := types.AmmPool{
		Id:               1,
		TotalLiabilities: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)},
	}

	suite.app.StablestakeKeeper.SetAmmPool(suite.ctx, p)
	r := suite.app.StablestakeKeeper.GetAmmPool(suite.ctx, 1)
	suite.Assert().Equal(p, r)

	r = suite.app.StablestakeKeeper.GetAmmPool(suite.ctx, 2)
	suite.Assert().Equal(types.AmmPool{
		Id:               2,
		TotalLiabilities: sdk.Coins{},
	}, r)

	all := suite.app.StablestakeKeeper.GetAllAmmPools(suite.ctx)
	suite.Assert().Equal([]types.AmmPool{p}, all)

	suite.app.StablestakeKeeper.AddPoolLiabilities(suite.ctx, p.Id, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	r = suite.app.StablestakeKeeper.GetAmmPool(suite.ctx, 1)
	p.AddLiabilities(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	suite.Assert().Equal(p, r)
}
