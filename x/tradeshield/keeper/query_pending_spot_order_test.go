package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func TestPendingSpotOrderQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPendingSpotOrder(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPendingSpotOrderRequest
		response *types.QueryGetPendingSpotOrderResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetPendingSpotOrderRequest{Id: msgs[0].Id},
			response: &types.QueryGetPendingSpotOrderResponse{PendingSpotOrder: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetPendingSpotOrderRequest{Id: msgs[1].Id},
			response: &types.QueryGetPendingSpotOrderResponse{PendingSpotOrder: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetPendingSpotOrderRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.PendingSpotOrder(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestPendingSpotOrderQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPendingSpotOrder(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPendingSpotOrderRequest {
		return &types.QueryAllPendingSpotOrderRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PendingSpotOrderAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PendingSpotOrder), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.PendingSpotOrder),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PendingSpotOrderAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PendingSpotOrder), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.PendingSpotOrder),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.PendingSpotOrderAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.PendingSpotOrder),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.PendingSpotOrderAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
