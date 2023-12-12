package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

func (suite *KeeperTestSuite) TestPriceMsgServerCreate() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	srv := keeper.NewMsgServerImpl(k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	suite.app.OracleKeeper.SetPriceFeeder(ctx, types.PriceFeeder{
		Feeder:   creator,
		IsActive: true,
	})
	for i := 0; i < 5; i++ {
		expected := &types.MsgFeedPrice{
			Provider: creator,
			Asset:    strconv.Itoa(i),
			Source:   "elys",
		}
		_, err := srv.FeedPrice(wctx, expected)
		suite.Require().NoError(err)
		rst, found := k.GetPrice(ctx, expected.Asset, expected.Source, uint64(ctx.BlockTime().Unix()))
		suite.Require().True(found)
		suite.Require().Equal(expected.Provider, rst.Provider)
	}
}
