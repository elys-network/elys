package keeper_test

// import (
// 	"strconv"
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
// 	"github.com/stretchr/testify/require"

// 	keepertest "github.com/elys-network/elys/testutil/keeper"
// 	"github.com/elys-network/elys/x/oracle/keeper"
// 	"github.com/elys-network/elys/x/oracle/types"
// )

// // Prevent strconv unused error
// var _ = strconv.IntSize

// func TestAssetInfoMsgServerCreate(t *testing.T) {
// 	k, ctx := keepertest.OracleKeeper(t)
// 	srv := keeper.NewMsgServerImpl(*k)
// 	wctx := sdk.WrapSDKContext(ctx)
// 	creator := "A"
// 	for i := 0; i < 5; i++ {
// 		expected := &types.MsgCreateAssetInfo{Creator: creator,
// 			Denom: "denom" + strconv.Itoa(i),
// 		}
// 		_, err := srv.CreateAssetInfo(wctx, expected)
// 		require.NoError(t, err)
// 		rst, found := k.GetAssetInfo(ctx, expected.Denom)
// 		require.True(t, found)
// 		require.Equal(t, expected.Denom, rst.Denom)
// 	}
// }

// func TestAssetInfoMsgServerUpdate(t *testing.T) {
// 	creator := "A"

// 	for _, tc := range []struct {
// 		desc    string
// 		request *types.MsgUpdateAssetInfo
// 		err     error
// 	}{
// 		{
// 			desc: "Completed",
// 			request: &types.MsgUpdateAssetInfo{Creator: creator,
// 				Denom: "denom" + strconv.Itoa(0),
// 			},
// 		},
// 		{
// 			desc: "Unauthorized",
// 			request: &types.MsgUpdateAssetInfo{Creator: "B",
// 				Denom: "denom" + strconv.Itoa(0),
// 			},
// 			err: sdkerrors.ErrUnauthorized,
// 		},
// 		{
// 			desc: "KeyNotFound",
// 			request: &types.MsgUpdateAssetInfo{Creator: creator,
// 				Denom: "denom" + strconv.Itoa(100000),
// 			},
// 			err: sdkerrors.ErrKeyNotFound,
// 		},
// 	} {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			k, ctx := keepertest.OracleKeeper(t)
// 			srv := keeper.NewMsgServerImpl(*k)
// 			wctx := sdk.WrapSDKContext(ctx)
// 			expected := &types.MsgCreateAssetInfo{Creator: creator,
// 				Denom: "denom" + strconv.Itoa(0),
// 			}
// 			_, err := srv.CreateAssetInfo(wctx, expected)
// 			require.NoError(t, err)

// 			_, err = srv.UpdateAssetInfo(wctx, tc.request)
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 			} else {
// 				require.NoError(t, err)
// 				rst, found := k.GetAssetInfo(ctx,
// 					expected.Denom,
// 				)
// 				require.True(t, found)
// 				require.Equal(t, expected.Denom, rst.Denom)
// 			}
// 		})
// 	}
// }

// func TestAssetInfoMsgServerDelete(t *testing.T) {
// 	creator := "A"

// 	for _, tc := range []struct {
// 		desc    string
// 		request *types.MsgDeleteAssetInfo
// 		err     error
// 	}{
// 		{
// 			desc: "Completed",
// 			request: &types.MsgDeleteAssetInfo{Creator: creator,
// 				Denom: "denom" + strconv.Itoa(0),
// 			},
// 		},
// 		{
// 			desc: "Unauthorized",
// 			request: &types.MsgDeleteAssetInfo{Creator: "B",
// 				Denom: "denom" + strconv.Itoa(0),
// 			},
// 			err: sdkerrors.ErrUnauthorized,
// 		},
// 		{
// 			desc: "KeyNotFound",
// 			request: &types.MsgDeleteAssetInfo{Creator: creator,
// 				Denom: "denom" + strconv.Itoa(100000),
// 			},
// 			err: sdkerrors.ErrKeyNotFound,
// 		},
// 	} {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			k, ctx := keepertest.OracleKeeper(t)
// 			srv := keeper.NewMsgServerImpl(*k)
// 			wctx := sdk.WrapSDKContext(ctx)

// 			_, err := srv.CreateAssetInfo(wctx, &types.MsgCreateAssetInfo{Creator: creator,
// 				Denom: "denom" + strconv.Itoa(0),
// 			})
// 			require.NoError(t, err)
// 			_, err = srv.DeleteAssetInfo(wctx, tc.request)
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 			} else {
// 				require.NoError(t, err)
// 				_, found := k.GetAssetInfo(ctx, tc.request.Denom)
// 				require.False(t, found)
// 			}
// 		})
// 	}
// }
