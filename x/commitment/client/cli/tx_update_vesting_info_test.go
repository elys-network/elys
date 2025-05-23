package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/v5/testutil/network"
	"github.com/elys-network/elys/v5/x/commitment/client/cli"
)

func setupNetwork(t *testing.T) *network.Network {
	t.Helper()

	cfg := network.DefaultConfig(t.TempDir())
	return network.New(t, cfg)
}

func TestGovUpdateVestingInfo(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...

	args := []string{
		"--title=test",
		"--summary=test",
		"--metadata=test",
		"--deposit=1000000uelys",
		"--base-denom=ueden",
		"--vesting-denom=uelys",
		"--num-epochs=100",
		"--vest-now-factor=1",
		"--num-max-vestings=10",
		"--expedited=true",
		"--from=" + val.Address.String(),
		"-y",
	}

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdUpdateVestingInfo(), args)
	require.NoError(t, err)
}
