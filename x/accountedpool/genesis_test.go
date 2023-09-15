package tvl_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	tvl "github.com/elys-network/elys/x/accountedpool"
	"github.com/elys-network/elys/x/accountedpool/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		AccountedPoolList: []types.AccountedPool{
			{
				PoolId: (uint64)(0),
			},
			{
				PoolId: (uint64)(1),
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AccountedPoolKeeper(t)
	tvl.InitGenesis(ctx, *k, genesisState)
	got := tvl.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.AccountedPoolList, got.AccountedPoolList)
	// this line is used by starport scaffolding # genesis/test/assert
}
