package keeper_test

func (suite *PerpetualKeeperTestSuite) TestResetStore() {
	suite.Require().Nil(suite.app.PerpetualKeeper.ResetStore(suite.ctx))
}
