package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/x/perpetual/client/cli"
)

func TestClosePosition(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...
	args := []string{
		"1",
		"10000000",
		"--from=" + val.Address.String(),
		"-y",
	}
	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdClose(), args)
	require.NoError(t, err)
}
