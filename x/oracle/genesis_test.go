package oracle_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	simapp "github.com/elys-network/elys/v7/app"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/oracle"
	"github.com/elys-network/elys/v7/x/oracle/types"
	"github.com/stretchr/testify/require"
)

const (
	initChain = true
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		AssetInfos: []types.AssetInfo{
			{
				Denom:   "satoshi",
				Display: "BTC",
			},
			{
				Denom:   "wei",
				Display: "ETH",
			},
		},
		Prices: []types.Price{
			{
				Asset: "BTC",
				Price: sdkmath.LegacyNewDec(30000),
			},
			{
				Asset: "ETH",
				Price: sdkmath.LegacyNewDec(2000),
			},
		},
		PriceFeeders: []types.PriceFeeder{
			{
				Feeder:   "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				IsActive: true,
			},
			{
				Feeder:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				IsActive: false,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	app := simapp.InitElysTestApp(initChain, t)
	ctx := app.BaseApp.NewContext(initChain)
	oracle.InitGenesis(ctx, app.OracleKeeper, genesisState)
	got := oracle.ExportGenesis(ctx, app.OracleKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	require.ElementsMatch(t, genesisState.AssetInfos, got.AssetInfos)
	require.ElementsMatch(t, genesisState.Prices, got.Prices)
	require.ElementsMatch(t, genesisState.PriceFeeders, got.PriceFeeders)
	// this line is used by starport scaffolding # genesis/test/assert
}
