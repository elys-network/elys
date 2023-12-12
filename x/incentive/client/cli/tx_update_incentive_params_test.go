package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/x/incentive/client/cli"
)

func setupNetwork(t *testing.T) *network.Network {
	t.Helper()

	cfg := network.DefaultConfig()
	return network.New(t, cfg)
}

func TestGovUpdateIncentiveParams(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...
	// [community-tax] [withdraw-addr-enabled] [reward-portion-for-lps] [elys-stake-tracking-rate] [max-eden-reward-apr-stakers] [max-eden-reward-par-lps] [distribution-epoch-for-stakers] [distribution-epoch-for-lps]
	args := []string{
		"0.00",
		"true",
		"0.60",
		"0.30",
		"10",
		"0.30",
		"0.30",
		"10",
		"10",
		"--title=test",
		"--summary=test",
		"--metadata=test",
		"--deposit=1000000uelys",
		"--from=" + val.Address.String(),
		"-y",
	}

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdUpdateIncentiveParams(), args)
	require.NoError(t, err)
}
