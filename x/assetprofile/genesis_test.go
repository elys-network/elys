package assetprofile_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/testutil/nullify"
	"github.com/elys-network/elys/v4/x/assetprofile"
	"github.com/elys-network/elys/v4/x/assetprofile/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		EntryList: []types.Entry{
			{
				BaseDenom: "0",
			},
			{
				BaseDenom: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AssetprofileKeeper(t)
	assetprofile.InitGenesis(ctx, *k, genesisState)
	got := assetprofile.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.EntryList, got.EntryList)
	// this line is used by starport scaffolding # genesis/test/assert
}
