package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/v5/testutil/keeper"
	"github.com/elys-network/elys/v5/testutil/nullify"
	"github.com/elys-network/elys/v5/x/tier/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestPortfolioQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)

	msgs := createNPortfolio(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPortfolioRequest
		response *types.QueryGetPortfolioResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPortfolioRequest{
				User: msgs[0].Creator,
			},
			response: &types.QueryGetPortfolioResponse{TotalPortfolio: msgs[0].Portfolio.String()},
		},
		{
			desc: "Second",
			request: &types.QueryGetPortfolioRequest{
				User: msgs[1].Creator,
			},
			response: &types.QueryGetPortfolioResponse{TotalPortfolio: msgs[0].Portfolio.String()},
		},
		{
			// TODO: update, should be empty
			desc: "KeyNotFound",
			request: &types.QueryGetPortfolioRequest{
				User: "cosmos1f5wt5hm3yj2x5etpqjzwzaar6cjhr0hvfcwn5s",
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Portfolio(ctx, tc.request)
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

func TestPortfolioQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)

	msgs := createNPortfolio(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPortfolioRequest {
		return &types.QueryAllPortfolioRequest{
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
			resp, err := keeper.PortfolioAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Portfolio), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Portfolio),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PortfolioAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Portfolio), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Portfolio),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.PortfolioAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Portfolio),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.PortfolioAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
