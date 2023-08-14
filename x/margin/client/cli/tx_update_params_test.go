package cli_test

import (
	"strconv"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/margin/client/cli"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestGovUpdateParams(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	app := simapp.InitElysTestApp(true)
	baseCtx := app.BaseApp.NewContext(true, tmproto.Header{})

	// Generate n random accounts with 1000000stake balanced
	addr := simapp.AddTestAddrs(app, baseCtx, 1, sdk.NewInt(1000000))

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...

	args := []string{
		"--title=test",
		"--summary=test",
		"--metadata=test",
		"--deposit=1000000uelys",
		"--leverage-max=5.0",
		"--interest-rate-min=1.0",
		"--interest-rate-max=1.0",
		"--interest-rate-increase=1.0",
		"--interest-rate-decrease=1.0",
		"--health-gain-factor=1.0",
		"--epoch-length=5",
		"--max-open-positions=100",
		"--removal-queue-threshold=0.1",
		"--pool-open-threshold=1.0",
		"--force-close-fund-percentage=1.0",
		"--force-close-fund-address=" + addr[0].String(),
		"--incremental-interest-payment-enabled=true",
		"--incremental-interest-payment-fund-percentage=1.0",
		"--incremental-interest-payment-fund-address=" + addr[0].String(),
		"--sq-modifier=0.1",
		"--safety-factor=0.1",
		"--whitelisting-enabled=true",
		"--from=" + val.Address.String(),
		"-y",
	}

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdUpdateParams(), args)
	require.NoError(t, err)
}
