package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func TestAirdropQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNAirdrop(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetAirdropRequest
		response *types.QueryGetAirdropResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetAirdropRequest{
				Intent: msgs[0].Intent,
			},
			response: &types.QueryGetAirdropResponse{Airdrop: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetAirdropRequest{
				Intent: msgs[1].Intent,
			},
			response: &types.QueryGetAirdropResponse{Airdrop: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetAirdropRequest{
				Intent: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Airdrop(wctx, tc.request)
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

func TestAirdropQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNAirdrop(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAirdropRequest {
		return &types.QueryAllAirdropRequest{
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
			resp, err := keeper.AirdropAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Airdrop), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Airdrop),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.AirdropAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Airdrop), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Airdrop),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AirdropAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Airdrop),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AirdropAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
