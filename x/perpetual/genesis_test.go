package perpetual_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/perpetual"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PoolList: []types.Pool{
			{
				AmmPoolId: 0,
			},
			{
				AmmPoolId: 1,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PerpetualKeeper(t)
	perpetual.InitGenesis(ctx, *k, genesisState)
	got := perpetual.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.PoolList, got.PoolList)
	// this line is used by starport scaffolding # genesis/test/assert
}
