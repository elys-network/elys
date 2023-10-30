package keeper_test

import (
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite KeeperTestSuite) TestBeginBlocker() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	mtp, _ := suite.OpenPosition(addr)
	ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, mtp.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetMTPHealth(suite.ctx, *mtp, ammPool)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "1.221000000000000000")

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 70))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	health, err = k.GetMTPHealth(suite.ctx, *mtp, ammPool)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "1.024543738200125865")

	params := k.GetParams(suite.ctx)
	params.SafetyFactor = sdk.NewDecWithPrec(11, 1)
	k.SetParams(suite.ctx, &params)
	k.BeginBlocker(suite.ctx)
	_, err = k.GetMTP(suite.ctx, mtp.Address, mtp.Id)
	suite.Require().Error(err)
}

func (suite KeeperTestSuite) TestLiquidatePositionIfUnhealthy() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	mtp, pool := suite.OpenPosition(addr)
	ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, mtp.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetMTPHealth(suite.ctx, *mtp, ammPool)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "1.221000000000000000")

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 70))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	health, err = k.GetMTPHealth(suite.ctx, *mtp, ammPool)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "1.024543738200125865")

	params := k.GetParams(suite.ctx)
	params.SafetyFactor = sdk.NewDecWithPrec(11, 1)
	k.SetParams(suite.ctx, &params)
	k.LiquidatePositionIfUnhealthy(suite.ctx, mtp, pool, ammPool)
	_, err = k.GetMTP(suite.ctx, mtp.Address, mtp.Id)
	suite.Require().Error(err)
}
