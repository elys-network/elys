package keeper_test

import (
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite KeeperTestSuite) TestBeginBlocker() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, _ := suite.OpenPosition(addr)
	ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, position.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position, ammPool)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.221000000000000000") // slippage enabled on amm
	suite.Require().Equal(health.String(), "1.250000000000000000") // slippage disabled on amm

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 70))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	health, err = k.GetPositionHealth(suite.ctx, *position, ammPool)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.024543738200125865") // slippage enabled on amm
	suite.Require().Equal(health.String(), "1.048877700860079715") // slippage disabled on amm

	params := k.GetParams(suite.ctx)
	params.SafetyFactor = sdk.NewDecWithPrec(11, 1)
	k.SetParams(suite.ctx, &params)
	k.BeginBlocker(suite.ctx)
	_, err = k.GetPosition(suite.ctx, position.Address, position.Id)
	suite.Require().Error(err)
}

func (suite KeeperTestSuite) TestLiquidatePositionIfUnhealthy() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, pool := suite.OpenPosition(addr)
	ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, position.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position, ammPool)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.221000000000000000") // slippage enabled on amm
	suite.Require().Equal(health.String(), "1.250000000000000000") // slippage disabled on amm

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 70))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	health, err = k.GetPositionHealth(suite.ctx, *position, ammPool)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.024543738200125865") // slippage enabled on amm
	suite.Require().Equal(health.String(), "1.048877700860079715") // slippage disabled on amm

	params := k.GetParams(suite.ctx)
	params.SafetyFactor = sdk.NewDecWithPrec(11, 1)
	k.SetParams(suite.ctx, &params)
	k.LiquidatePositionIfUnhealthy(suite.ctx, position, pool, ammPool)
	_, err = k.GetPosition(suite.ctx, position.Address, position.Id)
	suite.Require().Error(err)
}
