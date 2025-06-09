package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/v6/testutil/network"
	"github.com/elys-network/elys/v6/x/vaults/client/cli"
)

func setupNetwork(t *testing.T) *network.Network {
	t.Helper()
	cfg := network.DefaultConfig(t.TempDir())
	return network.New(t, cfg)
}

func TestCmdPerformAction(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid join pool action",
			args: []string{
				"1", // vault ID
				"--action", `{"join_pool":{"pool_id":1,"max_amounts_in":[{"denom":"uusdc","amount":"100000000"}],"share_amount_out":"100000000"}}`,
				"--from=" + val.Address.String(),
				"--gas=auto",
				"--gas-adjustment=1.5",
				"--gas-prices=0.1uelys",
				"-y",
			},
			expectErr: false,
		},
		// {
		// 	name: "invalid vault ID",
		// 	args: []string{
		// 		"invalid", // invalid vault ID
		// 		"--action", `{"join_pool":{"pool_id":1,"max_amounts_in":[{"denom":"uusdc","amount":"100000000"}],"share_amount_out":"100000000"}}`,
		// 		"--from=" + val.Address.String(),
		// 		"--gas=auto",
		// 		"--gas-adjustment=1.5",
		// 		"--gas-prices=0.1uelys",
		// 		"-y",
		// 	},
		// 	expectErr: true,
		// },
		// {
		// 	name: "missing action",
		// 	args: []string{
		// 		"1", // vault ID
		// 		"--from=" + val.Address.String(),
		// 		"--gas=auto",
		// 		"--gas-adjustment=1.5",
		// 		"--gas-prices=0.1uelys",
		// 		"-y",
		// 	},
		// 	expectErr: true,
		// },
		// {
		// 	name: "invalid action JSON",
		// 	args: []string{
		// 		"1", // vault ID
		// 		"--action", `invalid json`,
		// 		"--from=" + val.Address.String(),
		// 		"--gas=auto",
		// 		"--gas-adjustment=1.5",
		// 		"--gas-prices=0.1uelys",
		// 		"-y",
		// 	},
		// 	expectErr: true,
		// },
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPerformAction(), tc.args)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
