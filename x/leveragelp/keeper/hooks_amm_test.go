package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite KeeperTestSuite) TestCheckAmmPoolUsdcBalance() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	_, pool := suite.OpenPosition(addr)

	ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, pool.AmmPoolId)
	suite.Require().True(found)
	err := k.CheckAmmPoolUsdcBalance(suite.ctx, ammPool)
	suite.Require().NoError(err)

	usdcDenom := suite.app.StablestakeKeeper.GetParams(suite.ctx).DepositDenom
	suite.Require().Equal(ammPool.PoolAssets[0].Token.Denom, usdcDenom)

	// assume usdc amount reduced to 1000 and check
	ammPool.PoolAssets[0].Token.Amount = sdk.NewInt(1000)
	err = k.CheckAmmPoolUsdcBalance(suite.ctx, ammPool)
	suite.Require().Error(err)
}
