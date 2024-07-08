package keeper_test

import (
	"time"

	"cosmossdk.io/math"
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

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 500))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	health, err = k.GetPositionHealth(suite.ctx, *position, ammPool)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.024543738200125865") // slippage enabled on amm
	suite.Require().Equal(health.String(), "1.025220422390814025") // slippage disabled on amm

	params := k.GetParams(suite.ctx)
	params.SafetyFactor = sdk.NewDecWithPrec(11, 1)
	k.SetParams(suite.ctx, &params)
	k.BeginBlocker(suite.ctx)
	_, err = k.GetPosition(suite.ctx, position.Address, position.Id)
	suite.Require().Error(err)

	sortDec := math.LegacyNewDecFromInt(sdk.NewInt(100)).QuoInt(sdk.NewInt(10))
	suite.Require().Equal(sortDec.String(), "1.250000000000000000")
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

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 500))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	health, err = k.GetPositionHealth(suite.ctx, *position, ammPool)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.024543738200125865") // slippage enabled on amm
	suite.Require().Equal(health.String(), "1.025220422390814025") // slippage disabled on amm

	cacheCtx, _ := suite.ctx.CacheContext()
	params := k.GetParams(cacheCtx)
	params.SafetyFactor = sdk.NewDecWithPrec(11, 1)
	k.SetParams(cacheCtx, &params)
	isHealthy, earlyReturn := k.LiquidatePositionIfUnhealthy(cacheCtx, position, pool, ammPool)
	suite.Require().False(isHealthy)
	suite.Require().False(earlyReturn)
	_, err = k.GetPosition(cacheCtx, position.Address, position.Id)
	suite.Require().Error(err)

	cacheCtx, _ = suite.ctx.CacheContext()
	position.StopLossPrice = math.LegacyNewDec(100000)
	k.SetPosition(cacheCtx, position)
	underStopLossPrice, earlyReturn := k.ClosePositionIfUnderStopLossPrice(cacheCtx, position, pool, ammPool)
	suite.Require().True(underStopLossPrice)
	suite.Require().False(earlyReturn)
	_, err = k.GetPosition(cacheCtx, position.Address, position.Id)
	suite.Require().Error(err)
}

// Test sorted liquidity flow
// Add values
// Edge cases
