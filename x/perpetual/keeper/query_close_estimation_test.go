package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/testutil/sample"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestQueryCloseEstimation_InvalidRequest() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.CloseEstimation(ctx, nil)

	suite.Require().ErrorContains(err, "invalid request")
}

func (suite *PerpetualKeeperTestSuite) TestQueryCloseEstimation_CloseAmountIsNegative() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.CloseEstimation(ctx, &types.QueryCloseEstimationRequest{
		Address:     sample.AccAddress(),
		PositionId:  0,
		CloseAmount: math.NewInt(-200),
	})

	suite.Require().ErrorContains(err, "invalid close_amount")
}

func (suite *PerpetualKeeperTestSuite) TestQueryCloseEstimation_ErrAccAddressFromBech32() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.CloseEstimation(ctx, &types.QueryCloseEstimationRequest{
		Address:     "error",
		PositionId:  0,
		CloseAmount: math.NewInt(200),
	})

	suite.Require().ErrorContains(err, "invalid bech32 string length 5")
}

func (suite *PerpetualKeeperTestSuite) TestQueryCloseEstimation_ErrorGetMTP() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.CloseEstimation(ctx, &types.QueryCloseEstimationRequest{
		Address:     sample.AccAddress(),
		PositionId:  1,
		CloseAmount: math.NewInt(200),
	})

	suite.Require().ErrorContains(err, "mtp not found")
}

func (suite *PerpetualKeeperTestSuite) TestQueryCloseEstimation_ErrorGetPerpetualPool() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]
	secondPositionCreator := addr[1]

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

	secondOpenPositionMsg := &types.MsgOpen{
		Creator:         secondPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	openResponse, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)
	_, err = suite.app.PerpetualKeeper.Open(ctx, secondOpenPositionMsg)
	suite.Require().NoError(err)

	suite.app.PerpetualKeeper.RemovePool(ctx, firstPool)

	_, err = k.CloseEstimation(ctx, &types.QueryCloseEstimationRequest{
		Address:     firstPositionCreator.String(),
		PositionId:  openResponse.Id,
		CloseAmount: math.NewInt(200),
	})

	suite.Require().ErrorContains(err, "perpetual pool 1 not found")
}

func (suite *PerpetualKeeperTestSuite) TestQueryCloseEstimation_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]
	secondPositionCreator := addr[1]

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

	secondOpenPositionMsg := &types.MsgOpen{
		Creator:         secondPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	openResponse, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)
	_, err = suite.app.PerpetualKeeper.Open(ctx, secondOpenPositionMsg)
	suite.Require().NoError(err)

	_, err = k.CloseEstimation(ctx, &types.QueryCloseEstimationRequest{
		Address:     firstPositionCreator.String(),
		PositionId:  openResponse.Id,
		CloseAmount: math.NewInt(200),
	})

	suite.Require().Nil(err)
}
