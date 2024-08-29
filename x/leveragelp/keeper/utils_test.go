package keeper_test

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (suite KeeperTestSuite) TestCheckUserAuthorization() {
	// Create an instance of Keeper with the mock checker
	k := suite.app.LeveragelpKeeper
	pk := ed25519.GenPrivKey().PubKey()
	creator := sdk.AccAddress(pk.Address())
	msg := &types.MsgOpen{Creator: creator.String()}

	params := k.GetParams(suite.ctx)
	params.WhitelistingEnabled = true
	k.SetParams(suite.ctx, &params)
	k.WhitelistAddress(suite.ctx, creator)
	err := k.CheckUserAuthorization(suite.ctx, msg)
	suite.Require().NoError(err)

	k.DewhitelistAddress(suite.ctx, creator)
	err = k.CheckUserAuthorization(suite.ctx, msg)
	suite.Require().Error(err)

	params.WhitelistingEnabled = false
	k.SetParams(suite.ctx, &params)
	err = k.CheckUserAuthorization(suite.ctx, msg)
	suite.Require().NoError(err)
}

func (suite KeeperTestSuite) TestCheckSameAssets() {
	app := suite.app
	k := app.LeveragelpKeeper
	addr := simapp.AddTestAddrs(app, suite.ctx, 1, sdk.NewInt(1000000))
	suite.SetupCoinPrices(suite.ctx)

	position := types.NewPosition(addr[0].String(), sdk.NewInt64Coin("USDC", 0), sdk.NewDec(5), 1)
	k.SetPosition(suite.ctx, position)

	msg := &types.MsgOpen{
		Creator:          addr[0].String(),
		CollateralAsset:  "USDC",
		CollateralAmount: sdk.NewInt(100),
		AmmPoolId:        1,
		Leverage:         sdk.NewDec(1),
	}

	// Expect no error
	position, _ = k.CheckSamePosition(suite.ctx, msg)
	suite.Require().NotNil(position)
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
	suite.Require().Error(err)

	// PoolHealthTooLow
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.Pool{
		AmmPoolId: 1,
		Enabled:   false,
		Health:    sdk.NewDec(5),
	})
	err = k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().Error(err)

	// PoolIsHealthy
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.Pool{
		AmmPoolId: 1,
		Enabled:   true,
		Health:    sdk.NewDec(15),
		Closed:    false,
	})
	err = k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().NoError(err)
}

func (suite KeeperTestSuite) TestCheckMaxOpenPositions() {
	k := suite.app.LeveragelpKeeper

	params := k.GetParams(suite.ctx)
	params.MaxOpenPositions = 10
	k.SetParams(suite.ctx, &params)

	// OpenPositionsBelowMax
	k.SetOpenPositionCount(suite.ctx, 0)
	err := k.CheckMaxOpenPositions(suite.ctx)
	suite.Require().NoError(err)

	//  Expect an error about max open positions
	k.SetOpenPositionCount(suite.ctx, 10)
	err = k.CheckMaxOpenPositions(suite.ctx)
	suite.Require().Error(types.ErrMaxOpenPositions)

	// OpenPositionsExceedMax
	k.SetOpenPositionCount(suite.ctx, 11)
	err = k.CheckMaxOpenPositions(suite.ctx)
	suite.Require().Error(types.ErrMaxOpenPositions)
}

func (suite KeeperTestSuite) TestGetAmmPool() {
	k := suite.app.LeveragelpKeeper

	poolId := uint64(42)

	// PoolNotFound
	_, err := k.GetAmmPool(suite.ctx, poolId)
	suite.Require().True(errors.Is(err, types.ErrPoolDoesNotExist))

	// PoolFound
	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId: poolId,
	})
	_, err = k.GetAmmPool(suite.ctx, poolId)
	suite.Require().NoError(err)
}
