package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/testutil/nullify"
	"github.com/elys-network/elys/v6/x/tokenomics/types"
)

func TestTimeBasedInflationQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)

	msgs := createNTimeBasedInflation(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetTimeBasedInflationRequest
		response *types.QueryGetTimeBasedInflationResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetTimeBasedInflationRequest{
				StartBlockHeight: msgs[0].StartBlockHeight,
				EndBlockHeight:   msgs[0].EndBlockHeight,
			},
			response: &types.QueryGetTimeBasedInflationResponse{TimeBasedInflation: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetTimeBasedInflationRequest{
				StartBlockHeight: msgs[1].StartBlockHeight,
				EndBlockHeight:   msgs[1].EndBlockHeight,
			},
			response: &types.QueryGetTimeBasedInflationResponse{TimeBasedInflation: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetTimeBasedInflationRequest{
				StartBlockHeight: 100000,
				EndBlockHeight:   100000,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.TimeBasedInflation(ctx, tc.request)
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

func TestTimeBasedInflationQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)

	msgs := createNTimeBasedInflation(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllTimeBasedInflationRequest {
		return &types.QueryAllTimeBasedInflationRequest{
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
			resp, err := keeper.TimeBasedInflationAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.TimeBasedInflation), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.TimeBasedInflation),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.TimeBasedInflationAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.TimeBasedInflation), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.TimeBasedInflation),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.TimeBasedInflationAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.TimeBasedInflation),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.TimeBasedInflationAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
