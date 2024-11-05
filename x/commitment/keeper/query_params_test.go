package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.CommitmentKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}

// TestParamsQueryNilRequest tests the case where the request is nil
func TestParamsQueryNilRequest(t *testing.T) {
	keeper, ctx := testkeeper.CommitmentKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	_, err := keeper.Params(wctx, nil)
	require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
}
