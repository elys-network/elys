package masterchef_test

import (
	"testing"

	simapp "github.com/elys-network/elys/v5/app"

	"github.com/elys-network/elys/v5/testutil/nullify"
	"github.com/elys-network/elys/v5/x/masterchef"
	"github.com/elys-network/elys/v5/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

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
