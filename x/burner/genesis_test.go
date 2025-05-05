package burner_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/burner"
	"github.com/elys-network/elys/x/burner/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		HistoryList: []types.History{
			{
				Block:       1,
				BurnedCoins: sdk.Coins{sdk.NewInt64Coin("uusdc", 1)},
			},
			{
				Block:       2,
				BurnedCoins: sdk.Coins{sdk.NewInt64Coin("uatom", 1)},
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _ := keepertest.BurnerKeeper(t)
	burner.InitGenesis(ctx, *k, genesisState)
	got := burner.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.HistoryList, got.HistoryList)
	// this line is used by starport scaffolding # genesis/test/assert
}
