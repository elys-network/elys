package tradeshield_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tradeshield"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PendingSpotOrderList: []types.PendingSpotOrder{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		PendingSpotOrderCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TradeshieldKeeper(t)
	tradeshield.InitGenesis(ctx, *k, genesisState)
	got := tradeshield.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.PendingSpotOrderList, got.PendingSpotOrderList)
	require.Equal(t, genesisState.PendingSpotOrderCount, got.PendingSpotOrderCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
