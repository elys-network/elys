package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/x/commitment/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.CommitmentKeeper(t)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	params.TotalCommitted = nil
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}

// TestParamsQueryNilRequest tests the case where the request is nil
func TestParamsQueryNilRequest(t *testing.T) {
	keeper, ctx := testkeeper.CommitmentKeeper(t)

	_, err := keeper.Params(ctx, nil)
	require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
}
