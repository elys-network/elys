package tradeshield_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/tradeshield"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
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
				OrderId:      0,
				OwnerAddress: sdk.AccAddress([]byte("owner_address")).String(),
				PoolId:       1,
			},
			{
				OrderId:      1,
				OwnerAddress: sdk.AccAddress([]byte("owner_address")).String(),
				PoolId:       1,
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
