package cli_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/margin/client/cli"
)

func TestGovDeWhitelist(t *testing.T) {
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
		"--from=" + val.Address.String(),
		"-y",
	}

	args = append(args, addr[0].String())

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdDewhitelist(), args)
	require.NoError(t, err)
}
