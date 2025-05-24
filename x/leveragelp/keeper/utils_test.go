package keeper_test

import (
	"errors"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v5/app"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
	"github.com/elys-network/elys/v5/x/leveragelp/types"
)

func (suite *KeeperTestSuite) TestCheckUserAuthorization() {
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

func (suite *KeeperTestSuite) TestCheckSameAssets() {
	app := suite.app
	k := app.LeveragelpKeeper
	addr := simapp.AddTestAddrs(app, suite.ctx, 1, math.NewInt(1000000))
	suite.SetupCoinPrices(suite.ctx)

	position := types.NewPosition(addr[0].String(), sdk.NewInt64Coin("USDC", 0), 1)
	k.SetPosition(suite.ctx, position)

	msg := &types.MsgOpen{
		Creator:          addr[0].String(),
		CollateralAsset:  "USDC",
		CollateralAmount: math.NewInt(100),
		AmmPoolId:        1,
		Leverage:         math.LegacyNewDec(1),
	}

	// Expect no error
	position, _ = k.CheckSamePosition(suite.ctx, msg)
	suite.Require().NotNil(position)
}

func (suite *KeeperTestSuite) TestCheckPoolHealth() {
	k := suite.app.LeveragelpKeeper
	poolId := uint64(1)

	// PoolNotFound
	err := k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().True(errors.Is(err, types.ErrPoolDoesNotExist))

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	_, _, _ = suite.OpenPosition(addr)

	// PoolHealthTooLow
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.Pool{
		AmmPoolId: 1,
		Health:    math.LegacyNewDec(5).Quo(math.LegacyNewDec(100)),
	})
	err = k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().Error(errors.New("pool health too low to open new positions"))

	// PoolIsHealthy
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.Pool{
		AmmPoolId:          1,
		LeveragedLpAmount:  math.NewInt(100),
		Health:             math.LegacyNewDec(15),
		MaxLeveragelpRatio: math.LegacyNewDec(5),
	})
	err = k.CheckPoolHealth(suite.ctx, poolId)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestCheckMaxOpenPositions() {
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
	_ = k.CheckMaxOpenPositions(suite.ctx)
	suite.Require().Error(types.ErrMaxOpenPositions)

	// OpenPositionsExceedMax
	k.SetOpenPositionCount(suite.ctx, 11)
	_ = k.CheckMaxOpenPositions(suite.ctx)
	suite.Require().Error(types.ErrMaxOpenPositions)
}

func (suite *KeeperTestSuite) TestGetAmmPool() {
	k := suite.app.LeveragelpKeeper

	poolId := uint64(42)

	// PoolNotFound
	_, err := k.GetAmmPool(suite.ctx, poolId)
	suite.Require().True(errors.Is(err, types.ErrPoolDoesNotExist))

	// PoolFound
	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId:  poolId,
		Address: ammtypes.NewPoolAddress(poolId).String(),
	})
	_, err = k.GetAmmPool(suite.ctx, poolId)
	suite.Require().NoError(err)
}
