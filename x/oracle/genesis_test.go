package oracle_test

// import (
// 	"testing"

// 	keepertest "github.com/elys-network/elys/testutil/keeper"
// 	"github.com/elys-network/elys/testutil/nullify"
// 	"github.com/elys-network/elys/x/oracle"
// 	"github.com/elys-network/elys/x/oracle/types"
// 	"github.com/stretchr/testify/require"
// )

// func TestGenesis(t *testing.T) {
// 	genesisState := types.GenesisState{
// 		Params: types.DefaultParams(),
// 		PortId: types.PortID,
// 		AssetInfos: []types.AssetInfo{
// 			{
// 				Denom: "satoshi",
// 			},
// 			{
// 				Denom: "wei",
// 			},
// 		},
// 		Prices: []types.Price{
// 		{
// 			Index: "0",
// },
// 		{
// 			Index: "1",
// },
// 	},
// PriceFeeders: []types.PriceFeeder{
// 		{
// 			Index: "0",
// },
// 		{
// 			Index: "1",
// },
// 	},
// this line is used by starport scaffolding # genesis/test/state
// 	}

// 	k, ctx := keepertest.OracleKeeper(t)
// 	oracle.InitGenesis(ctx, *k, genesisState)
// 	got := oracle.ExportGenesis(ctx, *k)
// 	require.NotNil(t, got)

// 	nullify.Fill(&genesisState)
// 	nullify.Fill(got)

// 	require.Equal(t, genesisState.PortId, got.PortId)

// 	require.ElementsMatch(t, genesisState.AssetInfos, got.AssetInfos)
// 	require.ElementsMatch(t, genesisState.Prices, got.Prices)
// require.ElementsMatch(t, genesisState.PriceFeeders, got.PriceFeeders)
// this line is used by starport scaffolding # genesis/test/assert
// }
