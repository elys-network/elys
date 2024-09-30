package cli_test

import (
	"cosmossdk.io/math"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/perpetual/client/cli"
)

func TestGovUpdateParams(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	app := simapp.InitElysTestApp(true)
	basectx := app.BaseApp.NewContext(true)

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
		"--leverage-max=5.0",
		"--borrow-interest-rate-min=1.0",
		"--borrow-interest-rate-max=1.0",
		"--borrow-interest-rate-increase=1.0",
		"--borrow-interest-rate-decrease=1.0",
		"--health-gain-factor=1.0",
		"--epoch-length=5",
		"--max-open-positions=100",
		"--pool-open-threshold=1.0",
		"--force-close-fund-percentage=1.0",
		"--force-close-fund-address=" + addr[0].String(),
		"--incremental-borrow-interest-payment-enabled=true",
		"--incremental-borrow-interest-payment-fund-percentage=1.0",
		"--incremental-borrow-interest-payment-fund-address=" + addr[0].String(),
		"--safety-factor=0.1",
		"--whitelisting-enabled=true",
		"--from=" + val.Address.String(),
		"-y",
	}

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdUpdateParams(), args)
	require.NoError(t, err)
}
