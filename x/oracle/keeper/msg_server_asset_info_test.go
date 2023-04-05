package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func (suite *KeeperTestSuite) TestMsgServerSetAssetInfo() {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgSetAssetInfo
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgSetAssetInfo{Creator: creator,
				Denom: "denom" + strconv.Itoa(0),
			},
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OracleKeeper, suite.ctx
			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			_, err := srv.SetAssetInfo(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				rst, found := k.GetAssetInfo(ctx, tc.request.Denom)
				suite.Require().True(found)
				suite.Require().Equal(tc.request.Denom, rst.Denom)
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

			_, err := srv.SetAssetInfo(wctx, &types.MsgSetAssetInfo{
				Creator: creator,
				Denom:   "denom" + strconv.Itoa(0),
				Display: "DENOM" + strconv.Itoa(0),
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
