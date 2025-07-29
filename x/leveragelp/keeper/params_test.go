package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.LeveragelpKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, &params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
