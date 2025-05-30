package cli_test

import (
	"testing"

	"cosmossdk.io/math"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	simapp "github.com/elys-network/elys/v6/app"
	"github.com/elys-network/elys/v6/x/perpetual/client/cli"
)

func TestGovWhitelist(t *testing.T) {
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

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...

	args := []string{
		"--title=test",
		"--summary=test",
		"--metadata=test",
		"--deposit=1000000uelys",
		"--from=" + val.Address.String(),
		"-y",
	}

	args = append(args, addr[0].String())

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdWhitelist(), args)
	require.NoError(t, err)
}
