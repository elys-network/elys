package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func (suite *KeeperTestSuite) TestPriceFeederMsgServerUpdate() {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgSetPriceFeeder
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgSetPriceFeeder{
				Creator: creator,
				Feeder:  strconv.Itoa(0),
			},
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OracleKeeper, suite.ctx
			params := types.DefaultParams()
			params.ModuleAdmin = creator
			suite.app.OracleKeeper.SetParams(ctx, params)

			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			_, err := srv.SetPriceFeeder(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				rst, found := k.GetPriceFeeder(ctx,
					tc.request.Feeder,
				)
				suite.Require().True(found)
				suite.Require().Equal(tc.request.Feeder, rst.Feeder)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestPriceFeederMsgServerDelete() {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeletePriceFeeder
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeletePriceFeeder{Creator: creator,
				Feeder: strconv.Itoa(0),
			},
		},
	} {
		suite.Run(tc.desc, func() {
			k, ctx := suite.app.OracleKeeper, suite.ctx
			params := types.DefaultParams()
			params.ModuleAdmin = creator
			suite.app.OracleKeeper.SetParams(ctx, params)

			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			_, err := srv.SetPriceFeeder(wctx, &types.MsgSetPriceFeeder{
				Creator: creator,
				Feeder:  strconv.Itoa(0),
			})
			suite.Require().NoError(err)
			_, err = srv.DeletePriceFeeder(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				_, found := k.GetPriceFeeder(ctx, tc.request.Feeder)
				suite.Require().False(found)
			}
		})
	}
}
