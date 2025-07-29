package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.PerpetualKeeper(t)
	params := types.DefaultParams()

	err := k.SetParams(ctx, &params)
	require.NoError(t, err)

	require.EqualValues(t, params, k.GetParams(ctx))
}
