package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/elys-network/elys/v5/testutil/network"
	"github.com/elys-network/elys/v5/x/assetprofile/client/cli"
	"github.com/stretchr/testify/require"
)

func setupNetwork(t *testing.T) *network.Network {
	t.Helper()

	cfg := network.DefaultConfig(t.TempDir())
	return network.New(t, cfg)
}

func TestCmdAddEntry(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	args := []string{
		"axl",
		"6",
		"axl",
		"transfer/channel-11",
		"channel-11",
		"channel-198",
		"AXL",
		"AXL",
		"",
		"",
		"uaxl",
		"",
		"",
		"uaxl",
		"",
		"",
		"true",
		"true",
		"--from=" + val.Address.String(),
		"-y",
		"--gas=auto",
		"--gas-adjustment=2",
		"--fees=3stake",
	}

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdAddEntry(), args)

	require.NoError(t, err)
}
