package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.ParameterKeeper(t)

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
