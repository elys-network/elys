package keeper_test

import (
	"strconv"

	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/elys-network/elys/v4/x/oracle/keeper"
	"github.com/elys-network/elys/v4/x/oracle/types"
)

func (suite *KeeperTestSuite) TestPriceMsgServerCreate() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	srv := keeper.NewMsgServerImpl(k)
	creator := authtypes.NewModuleAddress("A").String()
	suite.app.OracleKeeper.SetPriceFeeder(ctx, types.PriceFeeder{
		Feeder:   creator,
		IsActive: true,
	})
	for i := 0; i < 5; i++ {
		expected := &types.MsgFeedPrice{
			Provider: creator,
			FeedPrice: types.FeedPrice{
				Asset:  strconv.Itoa(i),
				Source: "elys",
				Price:  math.LegacyOneDec(),
			},
		}
		_, err := srv.FeedPrice(ctx, expected)
		suite.Require().NoError(err)
		rst, found := k.GetPrice(ctx, expected.FeedPrice.Asset, expected.FeedPrice.Source, uint64(ctx.BlockTime().Unix()))
		suite.Require().True(found)
		suite.Require().Equal(expected.Provider, rst.Provider)
	}
}
