package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func (suite *KeeperTestSuite) TestPriceMsgServerCreate() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	srv := keeper.NewMsgServerImpl(k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgFeedPrice{
			Provider: creator,
			Asset:    strconv.Itoa(i),
			Source:   "binance",
		}
		_, err := srv.FeedPrice(wctx, expected)
		suite.Require().NoError(err)
		rst, found := k.GetPrice(ctx, expected.Asset, expected.Source, uint64(ctx.BlockTime().Unix()))
		suite.Require().True(found)
		suite.Require().Equal(expected.Provider, rst.Provider)
	}
}
