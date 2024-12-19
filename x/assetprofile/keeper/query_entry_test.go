package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/testutil/keeper"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/assetprofile/types"
)

func TestEntryQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.AssetprofileKeeper(t)
	msgs := createNEntry(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryEntryRequest
		response *types.QueryEntryResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryEntryRequest{
				BaseDenom: msgs[0].BaseDenom,
			},
			response: &types.QueryEntryResponse{Entry: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryEntryRequest{
				BaseDenom: msgs[1].BaseDenom,
			},
			response: &types.QueryEntryResponse{Entry: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryEntryRequest{
				BaseDenom: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Entry(ctx, tc.request)
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

func TestEntryQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.AssetprofileKeeper(t)
	msgs := createNEntry(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllEntryRequest {
		return &types.QueryAllEntryRequest{
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
			resp, err := keeper.EntryAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Entry), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Entry),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.EntryAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Entry), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Entry),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.EntryAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Entry),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.EntryAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestEntryByDenom(t *testing.T) {
	keeper, ctx := keeper.AssetprofileKeeper(t)
	msgs := createNEntry(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryEntryByDenomRequest
		response *types.QueryEntryByDenomResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryEntryByDenomRequest{
				Denom: msgs[0].Denom,
			},
			response: &types.QueryEntryByDenomResponse{Entry: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryEntryByDenomRequest{
				Denom: msgs[1].Denom,
			},
			response: &types.QueryEntryByDenomResponse{Entry: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryEntryByDenomRequest{
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
			response, err := keeper.EntryByDenom(ctx, tc.request)
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
