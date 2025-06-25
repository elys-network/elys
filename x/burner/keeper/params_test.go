package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/x/burner/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx, _ := testkeeper.BurnerKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, &params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
