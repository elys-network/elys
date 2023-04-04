package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func (suite *KeeperTestSuite) TestPriceFeederMsgServerCreate() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	srv := keeper.NewMsgServerImpl(k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreatePriceFeeder{Creator: creator,
			Feeder: strconv.Itoa(i),
		}
		_, err := srv.CreatePriceFeeder(wctx, expected)
		suite.Require().NoError(err)
		rst, found := k.GetPriceFeeder(ctx, expected.Feeder)
		suite.Require().True(found)
		suite.Require().Equal(expected.Feeder, rst.Feeder)
	}
}

func (suite *KeeperTestSuite) TestPriceFeederMsgServerUpdate() {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdatePriceFeeder
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdatePriceFeeder{Creator: creator,
				Feeder: strconv.Itoa(0),
			},
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OracleKeeper, suite.ctx
			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreatePriceFeeder{Creator: creator,
				Feeder: strconv.Itoa(0),
			}
			_, err := srv.CreatePriceFeeder(wctx, expected)
			suite.Require().NoError(err)

			_, err = srv.UpdatePriceFeeder(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				rst, found := k.GetPriceFeeder(ctx,
					expected.Feeder,
				)
				suite.Require().True(found)
				suite.Require().Equal(expected.Feeder, rst.Feeder)
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
			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreatePriceFeeder(wctx, &types.MsgCreatePriceFeeder{Creator: creator,
				Feeder: strconv.Itoa(0),
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
