package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestOpenConsolidate_ErrPoolDoesNotExist() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)

	mtp, err := k.GetMTP(ctx, firstPositionCreator, firstPosition.Id)
	suite.Require().NoError(err)

	suite.app.AmmKeeper.RemovePool(ctx, firstPool)

	positionCreator := addr[1]
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)

	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   math.LegacyZeroDec()}

	_, err = k.OpenConsolidate(
		ctx,
		&mtp,
		&mtp,
		msg,
		ptypes.BaseCurrency,
	)

	suite.Require().ErrorIs(err, types.ErrPoolDoesNotExist)
}

func (suite *PerpetualKeeperTestSuite) TestOpenConsolidate_ErrMTPUnhealthy() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)

	mtp, err := k.GetMTP(ctx, firstPositionCreator, firstPosition.Id)
	suite.Require().NoError(err)

	positionCreator := addr[1]
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)

	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   math.LegacyZeroDec()}

	params := k.GetParams(ctx)
	params.SafetyFactor = math.LegacyMustNewDecFromStr("1.30")
	k.SetParams(ctx, &params)

	_, err = k.OpenConsolidate(
		ctx,
		&mtp,
		&mtp,
		msg,
		ptypes.BaseCurrency,
	)

	suite.Require().ErrorIs(err, types.ErrMTPUnhealthy)
}

func (suite *PerpetualKeeperTestSuite) TestOpenConsolidate_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)

	mtp, err := k.GetMTP(ctx, firstPositionCreator, firstPosition.Id)
	suite.Require().NoError(err)

	positionCreator := addr[1]
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)

	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   math.LegacyZeroDec()}

	resp, err := k.OpenConsolidate(
		ctx,
		&mtp,
		&mtp,
		msg,
		ptypes.BaseCurrency,
	)

	suite.Require().Nil(err)
	suite.Require().Equal(resp.Id, firstPosition.Id)
}
