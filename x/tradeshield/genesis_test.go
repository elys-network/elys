package tradeshield_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/v5/testutil/keeper"
	"github.com/elys-network/elys/v5/testutil/nullify"
	"github.com/elys-network/elys/v5/x/tradeshield"
	"github.com/elys-network/elys/v5/x/tradeshield/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PendingSpotOrderList: []types.SpotOrder{
			{
				OrderId: 0,
			},
			{
				OrderId: 1,
			},
		},
		PendingSpotOrderCount: 2,
		PendingPerpetualOrderList: []types.PerpetualOrder{
			{
				OrderId: 0,
			},
			{
				OrderId: 1,
			},
		},
		PendingPerpetualOrderCount: 2,
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
	require.ElementsMatch(t, genesisState.PendingPerpetualOrderList, got.PendingPerpetualOrderList)
	require.Equal(t, genesisState.PendingPerpetualOrderCount, got.PendingPerpetualOrderCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
