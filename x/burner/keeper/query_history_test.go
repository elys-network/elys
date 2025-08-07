package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/burner/types"
)

func TestHistoryQuerySingle(t *testing.T) {
	keeper, ctx, _ := keepertest.BurnerKeeper(t)
	msgs := createNHistory(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetHistoryRequest
		response *types.QueryGetHistoryResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetHistoryRequest{
				Block: msgs[0].Block,
			},
			response: &types.QueryGetHistoryResponse{History: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetHistoryRequest{
				Block: msgs[1].Block,
			},
			response: &types.QueryGetHistoryResponse{History: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetHistoryRequest{
				Block: 100000,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.History(ctx, tc.request)
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

func TestHistoryQueryPaginated(t *testing.T) {
	keeper, ctx, _ := keepertest.BurnerKeeper(t)
	msgs := createNHistory(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllHistoryRequest {
		return &types.QueryAllHistoryRequest{
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
			resp, err := keeper.HistoryAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.History), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.History),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.HistoryAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.History), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.History),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.HistoryAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.History),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.HistoryAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
