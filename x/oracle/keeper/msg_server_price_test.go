package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
		expected := &types.MsgCreatePrice{
			Provider: creator,
			Asset:    strconv.Itoa(i),
		}
		_, err := srv.CreatePrice(wctx, expected)
		suite.Require().NoError(err)
		rst, found := k.GetPrice(ctx, expected.Asset)
		suite.Require().True(found)
		suite.Require().Equal(expected.Provider, rst.Provider)
	}
}

func (suite *KeeperTestSuite) TestPriceMsgServerUpdate() {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdatePrice
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdatePrice{
				Provider: creator,
				Asset:    strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdatePrice{
				Provider: "B",
				Asset:    strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdatePrice{
				Provider: creator,
				Asset:    strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OracleKeeper, suite.ctx
			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreatePrice{
				Provider: creator,
				Asset:    strconv.Itoa(0),
			}
			_, err := srv.CreatePrice(wctx, expected)
			suite.Require().NoError(err)

			_, err = srv.UpdatePrice(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				rst, found := k.GetPrice(ctx, expected.Asset)
				suite.Require().True(found)
				suite.Require().Equal(expected.Provider, rst.Provider)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestPriceMsgServerDelete() {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeletePrice
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeletePrice{Creator: creator,
				Asset: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeletePrice{Creator: "B",
				Asset: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeletePrice{Creator: creator,
				Asset: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OracleKeeper, suite.ctx
			srv := keeper.NewMsgServerImpl(k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreatePrice(wctx, &types.MsgCreatePrice{
				Provider: creator,
				Asset:    strconv.Itoa(0),
			})
			suite.Require().NoError(err)
			_, err = srv.DeletePrice(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				_, found := k.GetPrice(ctx, tc.request.Asset)
				suite.Require().False(found)
			}
		})
	}
}
