package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/testutil/sample"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestClose_ErrorGetMtp() {
	_, err := suite.app.PerpetualKeeper.Close(suite.ctx, &types.MsgClose{
		Creator: sample.AccAddress(),
		Id:      uint64(10),
		Amount:  math.NewInt(12000),
	})

	suite.Require().ErrorIs(err, types.ErrMTPDoesNotExist)
}
func (suite *PerpetualKeeperTestSuite) TestClose_ErrorGetEntry() {
	suite.SetupCoinPrices()

	addr := suite.AddAccounts(2, nil)
	amount := sdk.NewInt(1000)
	positionCreator := addr[1]
	poolId := uint64(1)
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
	suite.Require().NoError(err)

	poolCreator := addr[0]
	_ = suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	openPositionMsg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   sdk.ZeroDec(),
	}

	position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
	suite.Require().NoError(err)

	suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)

	_, err = suite.app.PerpetualKeeper.Close(suite.ctx, &types.MsgClose{
		Creator: positionCreator.String(),
		Id:      position.Id,
		Amount:  math.NewInt(500),
	})

	suite.Require().ErrorIs(err, assetprofiletypes.ErrAssetProfileNotFound)
}
func (suite *PerpetualKeeperTestSuite) TestClose_Long_NotEnoughLiquidity() {
	suite.SetupCoinPrices()

	addr := suite.AddAccounts(2, nil)
	amount := sdk.NewInt(1000)
	positionCreator := addr[1]
	poolId := uint64(1)
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
	suite.Require().NoError(err)

	poolCreator := addr[0]
	_ = suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	openPositionMsg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   sdk.ZeroDec(),
	}

	position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
	suite.Require().NoError(err)

	_, err = suite.app.PerpetualKeeper.Close(suite.ctx, &types.MsgClose{
		Creator: positionCreator.String(),
		Id:      position.Id,
		Amount:  math.NewInt(500),
	})

	suite.Require().EqualError(err, "not enough liquidity")
}
func (suite *PerpetualKeeperTestSuite) TestClose_Short() {
	suite.SetupCoinPrices()

	addr := suite.AddAccounts(2, nil)
	amount := sdk.NewInt(1000)
	positionCreator := addr[1]
	poolId := uint64(1)
	//tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
	//suite.Require().NoError(err)

	poolCreator := addr[0]
	_ = suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	openPositionMsg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   sdk.ZeroDec(),
	}

	position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
	suite.Require().NoError(err)

	_, err = suite.app.PerpetualKeeper.Close(suite.ctx, &types.MsgClose{
		Creator: positionCreator.String(),
		Id:      position.Id,
		Amount:  math.NewInt(900),
	})

	suite.Require().Nil(err)
}
func (suite *PerpetualKeeperTestSuite) TestClose_Long() {
	suite.SetupCoinPrices()

	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     ptypes.Elys,
		Liquidity: sdk.NewInt(2000000),
	})
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     ptypes.BaseCurrency,
		Liquidity: sdk.NewInt(1000000),
	})
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     ptypes.ATOM,
		Liquidity: sdk.NewInt(1000000),
	})

	addr := suite.AddAccounts(3, nil)
	amount := sdk.NewInt(1000)
	positionCreator := addr[1]
	poolId := uint64(1)
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
	suite.Require().NoError(err)

	poolCreator := addr[0]
	ammPool := suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	openPositionMsg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   sdk.ZeroDec(),
	}

	tokensIn := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000_000_000)), sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000_000_000)))
	suite.AddLiquidity(ammPool, addr[2], tokensIn)
	position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
	suite.Require().NoError(err)

	_, err = suite.app.PerpetualKeeper.Close(suite.ctx, &types.MsgClose{
		Creator: positionCreator.String(),
		Id:      position.Id,
		Amount:  math.NewInt(500),
	})

	suite.Require().Nil(err)
}
func (suite *PerpetualKeeperTestSuite) TestClose_Short_NotEnoughLiquidity() {
	suite.SetupCoinPrices()

	addr := suite.AddAccounts(2, nil)
	amount := sdk.NewInt(1000)
	positionCreator := addr[1]
	poolId := uint64(1)

	poolCreator := addr[0]
	_ = suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	openPositionMsg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   sdk.ZeroDec(),
	}

	position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
	suite.Require().NoError(err)

	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     ptypes.BaseCurrency,
		Liquidity: sdk.NewInt(0),
	})

	_, err = suite.app.PerpetualKeeper.Close(suite.ctx, &types.MsgClose{
		Creator: positionCreator.String(),
		Id:      position.Id,
		Amount:  math.NewInt(900),
	})

	suite.Require().EqualError(err, "not enough liquidity")
}
func (suite *PerpetualKeeperTestSuite) TestClose_NotEnoughLiquidity() {
	suite.SetupCoinPrices()

	addr := suite.AddAccounts(3, nil)
	amount := sdk.NewInt(1000)
	positionCreator := addr[1]
	poolId := uint64(1)
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
	suite.Require().NoError(err)

	poolCreator := addr[0]
	ammPool := suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	openPositionMsg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   sdk.ZeroDec(),
	}

	tokensIn := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000_000_000)), sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000_000_000)))
	suite.AddLiquidity(ammPool, addr[2], tokensIn)
	position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
	suite.Require().NoError(err)

	_, err = suite.app.PerpetualKeeper.Close(suite.ctx, &types.MsgClose{
		Creator: positionCreator.String(),
		Id:      position.Id,
		Amount:  math.NewInt(500),
	})

	suite.Require().EqualError(err, "not enough liquidity")
}
