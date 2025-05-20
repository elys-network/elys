package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/testutil/nullify"
	"github.com/elys-network/elys/v4/x/amm/types"
)

func TestDenomLiquidityQuerySingle(t *testing.T) {
	keeper, ctx, _, _ := keepertest.AmmKeeper(t)
	msgs := createNDenomLiquidity(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDenomLiquidityRequest
		response *types.QueryGetDenomLiquidityResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDenomLiquidityRequest{
				Denom: msgs[0].Denom,
			},
			response: &types.QueryGetDenomLiquidityResponse{DenomLiquidity: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetDenomLiquidityRequest{
				Denom: msgs[1].Denom,
			},
			response: &types.QueryGetDenomLiquidityResponse{DenomLiquidity: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDenomLiquidityRequest{
				Denom: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.DenomLiquidity(ctx, tc.request)
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

func TestDenomLiquidityQueryPaginated(t *testing.T) {
	keeper, ctx, _, _ := keepertest.AmmKeeper(t)
	msgs := createNDenomLiquidity(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDenomLiquidityRequest {
		return &types.QueryAllDenomLiquidityRequest{
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
			resp, err := keeper.DenomLiquidityAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DenomLiquidity), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.DenomLiquidity),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DenomLiquidityAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DenomLiquidity), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.DenomLiquidity),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.DenomLiquidityAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.DenomLiquidity),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.DenomLiquidityAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
