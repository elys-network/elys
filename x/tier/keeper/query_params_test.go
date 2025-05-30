package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/x/tier/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestParamsQueryRequest(t *testing.T) {
	keeper, ctx := testkeeper.MembershiptierKeeper(t)

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	_, err := keeper.Params(ctx, nil)
	want := status.Error(codes.InvalidArgument, "invalid request")

	require.ErrorIs(t, err, want)
}

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.MembershiptierKeeper(t)

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
