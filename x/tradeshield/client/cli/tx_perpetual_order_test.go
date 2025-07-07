package cli_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/v6/testutil/network"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/tradeshield/client/cli"
)

func TestCreatePerpetualOpenOrder(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	fields := []string{}
	tests := []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10))).String()),
				"long",                          // position
				"10",                            // leverage
				"1",                             // pool id
				"1000000" + ptypes.BaseCurrency, // collateral
				"0.5",                           // trigger price
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			//require.NoError(t, net.WaitForNextBlock()) // Need to figure out how to have next block for testing without provider

			args := []string{}
			args = append(args, fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreatePerpetualOpenOrder(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			var resp sdk.TxResponse
			require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
		})
	}
}

func TestCancelPerpertualOrders(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	tmpFile, err := os.CreateTemp("", "ids.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Create the correct format with PerpetualOrderPoolKey objects
	validOrders := []map[string]uint64{
		{
			"pool_id":  1,
			"order_id": 1,
		},
	}
	validJson, err := json.Marshal(validOrders)
	require.NoError(t, err)
	_, err = tmpFile.Write(validJson)
	require.NoError(t, err)
	tmpFile.Close()

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...

	args := []string{
		tmpFile.Name(),
		"--from=" + val.Address.String(),
		"-y",
	}

	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdCancelPerpetualOrders(), args)
	require.NoError(t, err)
}
