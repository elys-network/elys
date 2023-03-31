package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func (suite *KeeperTestSuite) TestAssetInfoMsgServerCreate() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	srv := keeper.NewMsgServerImpl(k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateAssetInfo{Creator: creator,
			Denom: "denom" + strconv.Itoa(i),
		}
		_, err := srv.CreateAssetInfo(wctx, expected)
		suite.Require().NoError(err)
		rst, found := k.GetAssetInfo(ctx, expected.Denom)
		suite.Require().True(found)
		suite.Require().Equal(expected.Denom, rst.Denom)
	}
}

func (suite *KeeperTestSuite) TestAssetInfoMsgServerUpdate() {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateAssetInfo
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateAssetInfo{Creator: creator,
				Denom: "denom" + strconv.Itoa(0),
			},
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OracleKeeper, suite.ctx
			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateAssetInfo{Creator: creator,
				Denom: "denom" + strconv.Itoa(0),
			}
			_, err := srv.CreateAssetInfo(wctx, expected)
			suite.Require().NoError(err)

			_, err = srv.UpdateAssetInfo(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				rst, found := k.GetAssetInfo(ctx,
					expected.Denom,
				)
				suite.Require().True(found)
				suite.Require().Equal(expected.Denom, rst.Denom)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAssetInfoMsgServerDelete() {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteAssetInfo
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteAssetInfo{Creator: creator,
				Denom: "denom" + strconv.Itoa(0),
			},
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OracleKeeper, suite.ctx
			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateAssetInfo(wctx, &types.MsgCreateAssetInfo{Creator: creator,
				Denom: "denom" + strconv.Itoa(0),
			})
			suite.Require().NoError(err)
			_, err = srv.DeleteAssetInfo(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				_, found := k.GetAssetInfo(ctx, tc.request.Denom)
				suite.Require().False(found)
			}
		})
	}
}
