package tvl_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/v5/testutil/keeper"
	"github.com/elys-network/elys/v5/testutil/nullify"
	tvl "github.com/elys-network/elys/v5/x/accountedpool"
	"github.com/elys-network/elys/v5/x/accountedpool/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
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
