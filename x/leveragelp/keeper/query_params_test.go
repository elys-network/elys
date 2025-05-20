package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.LeveragelpKeeper(t)

	params := types.DefaultParams()
	keeper.SetParams(ctx, &params)

	response, err := keeper.Params(ctx, &types.ParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.ParamsResponse{Params: params}, response)
}
