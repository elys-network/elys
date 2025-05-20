package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/testutil/nullify"
	"github.com/elys-network/elys/v4/x/accountedpool/types"
)

func TestAccountedPoolQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.AccountedPoolKeeper(t)
	msgs := createNAccountedPool(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetAccountedPoolRequest
		response *types.QueryGetAccountedPoolResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetAccountedPoolRequest{
				PoolId: msgs[0].PoolId,
			},
			response: &types.QueryGetAccountedPoolResponse{AccountedPool: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetAccountedPoolRequest{
				PoolId: msgs[1].PoolId,
			},
			response: &types.QueryGetAccountedPoolResponse{AccountedPool: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetAccountedPoolRequest{
				PoolId: (uint64)(100000),
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
			response, err := keeper.AccountedPool(ctx, tc.request)
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

func TestAccountedPoolQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.AccountedPoolKeeper(t)
	msgs := createNAccountedPool(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAccountedPoolRequest {
		return &types.QueryAllAccountedPoolRequest{
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
			resp, err := keeper.AccountedPoolAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AccountedPool), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.AccountedPool),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.AccountedPoolAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AccountedPool), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.AccountedPool),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AccountedPoolAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.AccountedPool),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AccountedPoolAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
