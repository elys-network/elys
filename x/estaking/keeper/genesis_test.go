package keeper_test

import (
	"testing"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/estaking/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k := app.EstakingKeeper
	k.InitGenesis(ctx, genesisState)
	got := k.ExportGenesis(ctx)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
