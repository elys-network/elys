package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/testutil/sample"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestMsgServerBrokerClose_ErrUnauthorised() {
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	_, err := msg.BrokerClose(suite.ctx, &types.MsgBrokerClose{
		Creator: sample.AccAddress(),
		Owner:   sample.AccAddress(),
		Id:      1,
		Amount:  math.NewInt(1),
	})
	suite.Require().ErrorIs(err, types.ErrUnauthorised)
}

func (suite *PerpetualKeeperTestSuite) TestMsgServerBrokerClose_CreatorIsNotBrokerAddress() {
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	_, err := msg.BrokerClose(suite.ctx, &types.MsgBrokerClose{
		Creator: sample.AccAddress(),
		Owner:   sample.AccAddress(),
		Id:      1,
		Amount:  math.NewInt(1),
	})
	suite.Require().ErrorIs(err, types.ErrUnauthorised)
}

func (suite *PerpetualKeeperTestSuite) TestMsgServerBrokerClose_Successful() {

	suite.SetupCoinPrices()

	msgServer := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	params := paramtypes.DefaultGenesis().Params
	suite.app.ParameterKeeper.SetParams(suite.ctx, params)
	_, err := msgServer.BrokerClose(suite.ctx, &types.MsgBrokerClose{
		Creator: params.BrokerAddress,
		Owner:   sample.AccAddress(),
		Id:      uint64(1),
		Amount:  math.NewInt(int64(200)),
	})
	suite.Require().Error(err)
}
