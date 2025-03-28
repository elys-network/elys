package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

    keepertest "github.com/elys-network/elys/testutil/keeper"
    "github.com/elys-network/elys/x/vaults/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.VaultsKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
