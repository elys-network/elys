package amm_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/amm"
	"github.com/elys-network/elys/v7/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PoolList: []types.Pool{
			{
				PoolId:  0,
				Address: types.NewPoolAddress(0).String(),
			},
			{
				PoolId:  1,
				Address: types.NewPoolAddress(1).String(),
			},
		},
		DenomLiquidityList: []types.DenomLiquidity{
			{
				Denom: "0",
			},
			{
				Denom: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _, _ := keepertest.AmmKeeper(t)
	amm.InitGenesis(ctx, *k, genesisState)
	got := amm.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.PoolList, got.PoolList)
	require.ElementsMatch(t, genesisState.DenomLiquidityList, got.DenomLiquidityList)
	// this line is used by starport scaffolding # genesis/test/assert
}
