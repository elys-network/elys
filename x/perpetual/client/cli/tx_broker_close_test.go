package cli_test

import (
	"testing"

	"cosmossdk.io/math"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/perpetual/client/cli"
)

func TestBrokerClosePosition(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	app := simapp.InitElysTestApp(true, t)
	basectx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, basectx)
	simapp.SetPerpetualParams(app, basectx)
	simapp.SetupAssetProfile(app, basectx)

	// Generate n random accounts with 1000000stake balanced
	addr := simapp.AddTestAddrs(app, basectx, 1, math.NewInt(1000000))
	args := []string{
		"1",
		"10000000",
		addr[0].String(),
		"--from=" + val.Address.String(),
		"-y",
	}
	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdBrokerClose(), args)
	require.NoError(t, err)
}
