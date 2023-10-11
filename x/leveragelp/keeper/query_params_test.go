package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.LeveragelpKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, &params)

	response, err := keeper.Params(wctx, &types.ParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.ParamsResponse{Params: params}, response)
}
