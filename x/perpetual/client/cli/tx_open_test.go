package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/v7/testutil/network"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/client/cli"
)

func setupNetwork(t *testing.T) *network.Network {
	t.Helper()

	cfg := network.DefaultConfig(t.TempDir())
	return network.New(t, cfg)
}

func TestOpenPosition(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...
	args := []string{
		"long",
		"1.5",
		"1",
		"1000" + ptypes.BaseCurrency,
		"--from=" + val.Address.String(),
		"-y",
	}
	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdOpen(), args)
	require.NoError(t, err)
}
