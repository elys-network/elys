package keeper_test

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (suite KeeperTestSuite) TestBeginBlocker() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, _ := suite.OpenPosition(addr)
	_, found := suite.app.AmmKeeper.GetPool(suite.ctx, position.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.221000000000000000") // slippage enabled on amm
	suite.Require().Equal("1.250000000000000000", health.String()) // slippage disabled on amm

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 500))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	suite.app.StablestakeKeeper.UpdateInterestStackedByAddress(suite.ctx, sdk.AccAddress(position.GetPositionAddress()))
	health, err = k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.024543738200125865") // slippage enabled on amm
	suite.Require().Equal("1.048855698433009587", health.String()) // slippage disabled on amm

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
	health, err := k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.221000000000000000") // slippage enabled on amm
	suite.Require().Equal(health.String(), "1.250000000000000000") // slippage disabled on amm

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 500))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	suite.app.StablestakeKeeper.UpdateInterestStackedByAddress(suite.ctx, sdk.AccAddress(position.GetPositionAddress()))
	health, err = k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.024543738200125865") // slippage enabled on amm
	suite.Require().Equal("1.048855698433009587", health.String()) // slippage disabled on amm

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
	k.SetPosition(cacheCtx, position, sdk.NewInt(0))
	underStopLossPrice, earlyReturn := k.ClosePositionIfUnderStopLossPrice(cacheCtx, position, pool, ammPool)
	suite.Require().True(underStopLossPrice)
	suite.Require().False(earlyReturn)
	_, err = k.GetPosition(cacheCtx, position.Address, position.Id)
	suite.Require().Error(err)
}

func (suite KeeperTestSuite) TestLiquidatePositionSorted() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, _ := suite.OpenPosition(addr)

	// open positions with other addresses
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr4 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr5 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	usdcTokenTotal := sdk.NewInt64Coin("uusdc", 500000)
	usdcToken := sdk.NewInt64Coin("uusdc", 100000)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{usdcTokenTotal})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr2, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr3, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr4, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr5, sdk.Coins{usdcToken})
	suite.Require().NoError(err)

	stableMsgServer := stablestakekeeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
	_, err = stableMsgServer.Bond(sdk.WrapSDKContext(suite.ctx), &stablestaketypes.MsgBond{
		Creator: addr2.String(),
		Amount:  sdk.NewInt(100000),
	})
	suite.Require().NoError(err)

	position3, err := k.OpenLong(suite.ctx, &types.MsgOpen{
		Creator:          addr3.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdk.NewInt(2000),
		AmmPoolId:        1,
		Leverage:         sdk.NewDec(2),
	})
	suite.Require().NoError(err)

	position4, err := k.OpenLong(suite.ctx, &types.MsgOpen{
		Creator:          addr4.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdk.NewInt(2000),
		AmmPoolId:        1,
		Leverage:         sdk.NewDec(6),
	})
	suite.Require().NoError(err)

	position5, err := k.OpenLong(suite.ctx, &types.MsgOpen{
		Creator:          addr5.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdk.NewInt(2000),
		AmmPoolId:        1,
		Leverage:         sdk.NewDec(4),
	})
	suite.Require().NoError(err)

	_, found := suite.app.AmmKeeper.GetPool(suite.ctx, position3.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position3)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "2.000000000000000000") // slippage disabled on amm

	health, err = k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "1.250000000000000000") // slippage disabled on amm

	health, err = k.GetPositionHealth(suite.ctx, *position4)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "1.200000000000000000") // slippage disabled on amm

	health, err = k.GetPositionHealth(suite.ctx, *position5)
	suite.Require().NoError(err)
	suite.Require().Equal(health.String(), "1.333333333333333333") // slippage disabled on amm

	// Check order in list
	suite.app.LeveragelpKeeper.IteratePoolPosIdsLiquidationSorted(suite.ctx, position.AmmPoolId, func(posId types.AddressId) bool {
		position, _ := k.GetPosition(suite.ctx, posId.Address, posId.Id)
		health, _ := k.GetPositionHealth(suite.ctx, position)
		fmt.Printf("Address: %s, Id: %d, value: %s\n", position.Address, position.Id, health.String())
		return false
	})

	err = k.ProcessAddCollateral(suite.ctx, addr4.String(), position4.Id, sdk.NewInt(1000))
	suite.Require().NoError(err)

	// Check order in list
	suite.app.LeveragelpKeeper.IteratePoolPosIdsLiquidationSorted(suite.ctx, position.AmmPoolId, func(posId types.AddressId) bool {
		position, _ := k.GetPosition(suite.ctx, posId.Address, posId.Id)
		health, _ := k.GetPositionHealth(suite.ctx, position)
		fmt.Printf("Address: %s, Id: %d, value: %s\n", position.Address, position.Id, health.String())
		return false
	})

	// add more lev
	k.OpenConsolidate(suite.ctx, position5, &types.MsgOpen{
		Creator:          addr5.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdk.NewInt(1000),
		AmmPoolId:        1,
		Leverage:         sdk.NewDec(4),
	})
	suite.Require().NoError(err)

	// Check order in list
	suite.app.LeveragelpKeeper.IteratePoolPosIdsLiquidationSorted(suite.ctx, position.AmmPoolId, func(posId types.AddressId) bool {
		position, _ := k.GetPosition(suite.ctx, posId.Address, posId.Id)
		health, _ := k.GetPositionHealth(suite.ctx, position)
		fmt.Printf("Address: %s, Id: %d, value: %s\n", position.Address, position.Id, health.String())
		return false
	})

	// Partial close.
	var (
		msg = &types.MsgClose{
			Creator:  addr5.String(),
			Id:       position5.Id,
			LpAmount: position5.LeveragedLpAmount.Quo(sdk.NewInt(2)),
		}
	)
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour))

	_, _, err = k.CloseLong(suite.ctx, msg)
	suite.Require().NoError(err)

	// Check order in list
	suite.app.LeveragelpKeeper.IteratePoolPosIdsLiquidationSorted(suite.ctx, position.AmmPoolId, func(posId types.AddressId) bool {
		position, _ := k.GetPosition(suite.ctx, posId.Address, posId.Id)
		health, _ := k.GetPositionHealth(suite.ctx, position)
		fmt.Printf("Address: %s, Id: %d, value: %s\n", position.Address, position.Id, health.String())
		return false
	})
}

// Add stablestake update hook test
