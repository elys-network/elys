package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx, _, _ := testkeeper.AmmKeeper(t)
	params := types.DefaultParams()
	params.BaseAssets = nil
	params.AllowedUpfrontSwapMakers = nil
	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
