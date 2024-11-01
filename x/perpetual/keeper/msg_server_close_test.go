package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/testutil/sample"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestMsgServerClose() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx
	msg := keeper.NewMsgServerImpl(*k)
	_, err := msg.Close(ctx, &types.MsgClose{
		Creator: sample.AccAddress(),
		Id:      2,
		Amount:  math.NewInt(int64(2)),
	})

	suite.app.AssetprofileKeeper.RemoveEntry(ctx, ptypes.BaseCurrency)

	suite.Require().Error(err)
}
