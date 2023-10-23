package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite KeeperTestSuite) TestCheckUserAuthorization() {
	// Create an instance of Keeper with the mock checker
	k := suite.app.LeveragelpKeeper
	msg := &types.MsgOpen{Creator: "whitelistedUser"}

	params := k.GetParams(suite.ctx)
	params.WhitelistingEnabled = true
	k.SetParams(suite.ctx, &params)
	k.WhitelistAddress(suite.ctx, msg.Creator)
	err := k.CheckUserAuthorization(suite.ctx, msg)
	suite.Require().NoError(err)

	k.DewhitelistAddress(suite.ctx, msg.Creator)
	err = k.CheckUserAuthorization(suite.ctx, msg)
	suite.Require().Error(err)

	params.WhitelistingEnabled = false
	k.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
}

func (suite KeeperTestSuite) TestCheckSameAssets() {
	app := suite.app
	k := app.LeveragelpKeeper

	mtp := types.NewMTP("creator", ptypes.BaseCurrency, sdk.NewDec(5), 1)
	k.SetMTP(suite.ctx, mtp)

	msg := &types.MsgOpen{
		Creator:          "creator",
		CollateralAsset:  ptypes.ATOM,
		CollateralAmount: sdk.NewInt(100),
		AmmPoolId:        1,
		Leverage:         sdk.NewDec(1),
	}

	// Expect no error
	mtp = k.CheckSamePosition(suite.ctx, msg)
	suite.Require().NotNil(mtp)
}

func (suite KeeperTestSuite) TestCheckPoolHealth() {
	k := suite.app.LeveragelpKeeper
	poolId := uint64(1)

	// PoolNotFound
	err := k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().True(errors.Is(err, types.ErrInvalidBorrowingAsset))

	// PoolDisabledOrClosed
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.Pool{
		AmmPoolId: 1,
		Enabled:   false,
	})
	err = k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().True(errors.Is(err, types.ErrInvalidBorrowingAsset))

	// PoolHealthTooLow
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.Pool{
		AmmPoolId: 1,
		Enabled:   false,
		Health:    sdk.NewDec(5),
	})
	err = k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().True(errors.Is(err, types.ErrInvalidPosition))

	// PoolIsHealthy
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.Pool{
		AmmPoolId: 1,
		Enabled:   false,
		Health:    sdk.NewDec(15),
	})
	err = k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().NoError(err)
}

func (suite KeeperTestSuite) TestCheckMaxOpenPositions(t *testing.T) {
	k := suite.app.LeveragelpKeeper

	// OpenPositionsBelowMax
	err := k.CheckMaxOpenPositions(suite.ctx)
	suite.Require().NoError(err)

	//  Expect an error about max open positions
	k.SetOpenMTPCount(suite.ctx, 10)
	params := k.GetParams(suite.ctx)
	params.MaxOpenPositions = 10
	k.SetParams(suite.ctx, &params)
	err = k.CheckMaxOpenPositions(suite.ctx)
	suite.Require().True(errors.Is(err, types.ErrMaxOpenPositions))

	// OpenPositionsExceedMax
	k.SetOpenMTPCount(suite.ctx, 11)
	err = k.CheckMaxOpenPositions(suite.ctx)
	suite.Require().True(errors.Is(err, types.ErrMaxOpenPositions))
}

func (suite KeeperTestSuite) TestGetAmmPool() {
	k := suite.app.LeveragelpKeeper

	ctx := sdk.Context{} // mock or setup a context
	poolId := uint64(42)

	// PoolNotFound
	_, err := k.GetAmmPool(ctx, poolId)
	suite.Require().True(errors.Is(err, types.ErrPoolDoesNotExist))

	// PoolFound
	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId: poolId,
	})
	_, err = k.GetAmmPool(ctx, poolId)
	suite.Require().NoError(err)
}
