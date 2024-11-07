package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *PerpetualKeeperTestSuite) TestHookAmm_AfterJoinPool() {
	ctx := suite.ctx
	k := suite.app.PerpetualKeeper
	_, poolCreator, ammPool := suite.SetPerpetualPool(uint64(1))

	ammHooks := k.AmmHooks()
	err := ammHooks.AfterJoinPool(ctx, poolCreator, ammPool, sdk.NewCoins(), math.NewInt(200))
	suite.Require().Nil(err)
}

func (suite *PerpetualKeeperTestSuite) TestHookAmm_AfterExitPool() {
	ctx := suite.ctx
	k := suite.app.PerpetualKeeper

	_, poolCreator, ammPool := suite.SetPerpetualPool(uint64(1))

	ammHooks := k.AmmHooks()
	err := ammHooks.AfterExitPool(ctx, poolCreator, ammPool, math.NewInt(200), sdk.NewCoins())
	suite.Require().Nil(err)
}

func (suite *PerpetualKeeperTestSuite) TestHookAmm_AfterSWap() {
	ctx := suite.ctx
	k := suite.app.PerpetualKeeper

	_, poolCreator, ammPool := suite.SetPerpetualPool(uint64(1))

	ammHooks := k.AmmHooks()
	err := ammHooks.AfterSwap(ctx, poolCreator, ammPool, sdk.NewCoins(), sdk.NewCoins())
	suite.Require().Nil(err)
}
