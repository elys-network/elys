package masterchef_test

import (
	simapp "github.com/elys-network/elys/app"
	"testing"

	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/masterchef"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	app := simapp.InitElysTestApp(true)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}
	ctx := app.BaseApp.NewContext(true)
	k := app.MasterchefKeeper
	masterchef.InitGenesis(ctx, k, genesisState)
	got := masterchef.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
