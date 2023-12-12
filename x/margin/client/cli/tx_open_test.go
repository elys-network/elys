package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/x/margin/client/cli"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func setupNetwork(t *testing.T) *network.Network {
	t.Helper()

	cfg := network.DefaultConfig()
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
		"open",
		"1.5",
		"uatom",
		"1000" + ptypes.BaseCurrency,
		"--from=" + val.Address.String(),
		"-y",
	}
	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdOpen(), args)
	require.NoError(t, err)
}
