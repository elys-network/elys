package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	simapp "github.com/elys-network/elys/v5/app"
	"github.com/elys-network/elys/v5/testutil/network"
	"github.com/elys-network/elys/v5/x/tradeshield/client/cli"
)

func setupNetwork(t *testing.T) *network.Network {
	t.Helper()

	cfg := network.DefaultConfig(t.TempDir())
	return network.New(t, cfg)
}

func TestGovUpdateParams(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	app := simapp.InitElysTestApp(true, t)
	basectx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, basectx)
	simapp.SetPerpetualParams(app, basectx)
	simapp.SetupAssetProfile(app, basectx)

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...

	args := []string{
		"--title=test",
		"--summary=test",
		"--metadata=test",
		"--deposit=1000000uelys",
		"--market-order-enabled=true",
		"--stake-enabled=true",
		"--process-orders-enabled=true",
		"--swap-enabled=true",
		"--perpetual-enabled=true",
		"--reward-enabled=true",
		"--leverage-enabled=true",
		"--limit-process-order=10000",
		"--reward-percentage=0.1",
		"--margin-error=0.1",
		"--minimum-deposit=1000000",
		"--from=" + val.Address.String(),
		"-y",
	}

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdUpdateParams(), args)
	require.NoError(t, err)
}
