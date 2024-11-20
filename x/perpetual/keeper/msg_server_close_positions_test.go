package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/sample"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestClosePositions_GetEntryError() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx
	msg := keeper.NewMsgServerImpl(*k)

	suite.app.AssetprofileKeeper.RemoveEntry(ctx, ptypes.BaseCurrency)

	_, err := msg.ClosePositions(ctx, &types.MsgClosePositions{Creator: sample.AccAddress()})
	suite.Require().Nil(err)
}

func (suite *PerpetualKeeperTestSuite) TestMsgServerClose_HandleLiquidateErrors() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx
	msg := keeper.NewMsgServerImpl(*k)

	firstPool := uint64(1)
	secondPool := uint64(2)

	suite.SetPerpetualPool(firstPool)
	suite.SetPerpetualPool(secondPool)

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

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)

	secondOpenPositionMsg := &types.MsgOpen{
		Creator:         secondPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          secondPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	secondPosition, err := suite.app.PerpetualKeeper.Open(ctx, secondOpenPositionMsg)
	suite.Require().NoError(err)

	k.RemovePool(ctx, firstPool)
	suite.app.AmmKeeper.RemovePool(ctx, secondPool)

	_, err = msg.ClosePositions(ctx, &types.MsgClosePositions{
		Creator: sample.AccAddress(),
		Liquidate: []types.PositionRequest{
			{
				Address: firstPositionCreator.String(),
				Id:      firstPosition.Id,
			},
			{
				Address: secondPositionCreator.String(),
				Id:      secondPosition.Id,
			},
			{
				Address: sample.AccAddress(),
				Id:      2000,
			},
		},
		StopLoss: []types.PositionRequest{
			{
				Address: firstPositionCreator.String(),
				Id:      firstPosition.Id,
			},
			{
				Address: sample.AccAddress(),
				Id:      2000,
			},
		},
		TakeProfit: []types.PositionRequest{
			{
				Address: firstPositionCreator.String(),
				Id:      firstPosition.Id,
			},
			{
				Address: sample.AccAddress(),
				Id:      2000,
			},
		},
	})
	suite.Require().Nil(err)

}

func (suite *PerpetualKeeperTestSuite) TestMsgServerClose_HandleLiquidateCheck() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx
	msg := keeper.NewMsgServerImpl(*k)

	firstPool := uint64(1)

	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)

	addr := suite.AddAccounts(1, nil)
	firstPositionCreator := addr[0]

	suite.app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.ATOM,
		Denom:     ptypes.ATOM,
		Decimals:  6,
	})

	tradingAssetPrice, err := k.GetAssetPrice(suite.ctx, ptypes.ATOM)
	suite.Require().NoError(err)

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          firstPool,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)

	_, err = msg.ClosePositions(ctx, &types.MsgClosePositions{
		Creator: sample.AccAddress(),
		Liquidate: []types.PositionRequest{
			{
				Address: firstPositionCreator.String(),
				Id:      firstPosition.Id,
			},
		},
		StopLoss: []types.PositionRequest{{
			Address: firstPositionCreator.String(),
			Id:      firstPosition.Id,
		}},
		TakeProfit: []types.PositionRequest{{
			Address: firstPositionCreator.String(),
			Id:      firstPosition.Id,
		}},
	})
	suite.Require().Nil(err)

}
