package keeper_test

func (suite *PerpetualKeeperTestSuite) TestGetEpochPosition() {
	k := suite.app.PerpetualKeeper

	ctx := suite.ctx
	ctx = ctx.WithBlockHeight(123)
	currentBlock := k.GetEpochPosition(ctx, 100)
	blockWant := int64(23)
	suite.Require().Equal(currentBlock, blockWant)

	// epoch length is 0
	currentBlock = k.GetEpochPosition(ctx, 0)
	blockWant = int64(0)
	suite.Require().Equal(currentBlock, blockWant)
}
