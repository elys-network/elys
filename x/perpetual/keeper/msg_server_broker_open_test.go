package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/sample"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestMsgServerBrokerOpen_ErrUnauthorised() {
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	_, err := msg.BrokerOpen(suite.ctx, &types.MsgBrokerOpen{
		Creator:         sample.AccAddress(),
		Owner:           sample.AccAddress(),
		PoolId:          1,
		StopLossPrice:   math.LegacyMustNewDecFromStr("2"),
		TradingAsset:    ptypes.ATOM,
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyNewDec(2),
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100)),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("1"),
	})
	suite.Require().ErrorIs(err, types.ErrUnauthorised)
}

func (suite *PerpetualKeeperTestSuite) TestMsgServerBrokerOpen_CreatorIsNotBrokerAddress() {
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	_, err := msg.BrokerOpen(suite.ctx, &types.MsgBrokerOpen{
		Creator:         sample.AccAddress(),
		Owner:           sample.AccAddress(),
		PoolId:          1,
		StopLossPrice:   math.LegacyMustNewDecFromStr("2"),
		TradingAsset:    ptypes.ATOM,
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyNewDec(2),
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100)),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("1"),
	})
	suite.Require().ErrorIs(err, types.ErrUnauthorised)
}

func (suite *PerpetualKeeperTestSuite) TestMsgServerBrokerOpen_Successful() {

	suite.SetupCoinPrices()
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)

	msgServer := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	params := paramtypes.DefaultGenesis().Params
	suite.app.ParameterKeeper.SetParams(suite.ctx, params)
	_, err = msgServer.BrokerOpen(suite.ctx, &types.MsgBrokerOpen{
		Creator:         params.BrokerAddress,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyNewDec(5),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(int64(2000))),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		Owner:           sample.AccAddress(),
		StopLossPrice:   math.LegacyZeroDec(),
		PoolId:          1,
	})
	suite.Require().Error(err)
}
