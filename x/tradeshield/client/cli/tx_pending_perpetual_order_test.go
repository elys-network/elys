package cli_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/x/tradeshield/client/cli"
)

func TestCreatePendingPerpetualOrder(t *testing.T) {
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
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			require.NoError(t, net.WaitForNextBlock())

			args := []string{}
			args = append(args, fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreatePendingPerpetualOrder(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			var resp sdk.TxResponse
			require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, clitestutil.CheckTxCode(net, ctx, resp.TxHash, tc.code))
		})
	}
}

// TODO: Add this in message task
// func TestUpdatePendingPerpetualOrder(t *testing.T) {
// 	net := network.New(t)

// 	val := net.Validators[0]
// 	ctx := val.ClientCtx

// 	fields := []string{"xyz"}
// 	common := []string{
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10))).String()),
// 	}
// 	args := []string{}
// 	args = append(args, fields...)
// 	args = append(args, common...)
// 	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreatePendingPerpetualOrder(), args)
// 	require.NoError(t, err)

// 	tests := []struct {
// 		desc string
// 		id   string
// 		args []string
// 		code uint32
// 		err  error
// 	}{
// 		{
// 			desc: "valid",
// 			id:   "0",
// 			args: common,
// 		},
// 		{
// 			desc: "key not found",
// 			id:   "1",
// 			args: common,
// 			code: sdkerrors.ErrKeyNotFound.ABCICode(),
// 		},
// 		{
// 			desc: "invalid key",
// 			id:   "invalid",
// 			err:  strconv.ErrSyntax,
// 		},
// 	}
// 	for _, tc := range tests {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			require.NoError(t, net.WaitForNextBlock())

// 			args := []string{tc.id}
// 			args = append(args, fields...)
// 			args = append(args, tc.args...)
// 			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdUpdatePendingPerpetualOrder(), args)
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 				return
// 			}
// 			require.NoError(t, err)

// 			var resp sdk.TxResponse
// 			require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 			require.NoError(t, clitestutil.CheckTxCode(net, ctx, resp.TxHash, tc.code))
// 		})
// 	}
// }

// func TestDeletePendingPerpetualOrder(t *testing.T) {
// 	net := network.New(t)

// 	val := net.Validators[0]
// 	ctx := val.ClientCtx

// 	fields := []string{"xyz"}
// 	common := []string{
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10))).String()),
// 	}
// 	args := []string{}
// 	args = append(args, fields...)
// 	args = append(args, common...)
// 	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreatePendingPerpetualOrder(), args)
// 	require.NoError(t, err)

// 	tests := []struct {
// 		desc string
// 		id   string
// 		args []string
// 		code uint32
// 		err  error
// 	}{
// 		{
// 			desc: "valid",
// 			id:   "0",
// 			args: common,
// 		},
// 		{
// 			desc: "key not found",
// 			id:   "1",
// 			args: common,
// 			code: sdkerrors.ErrKeyNotFound.ABCICode(),
// 		},
// 		{
// 			desc: "invalid key",
// 			id:   "invalid",
// 			err:  strconv.ErrSyntax,
// 		},
// 	}
// 	for _, tc := range tests {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			require.NoError(t, net.WaitForNextBlock())

// 			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdDeletePendingPerpetualOrder(), append([]string{tc.id}, tc.args...))
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 				return
// 			}
// 			require.NoError(t, err)

// 			var resp sdk.TxResponse
// 			require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 			require.NoError(t, clitestutil.CheckTxCode(net, ctx, resp.TxHash, tc.code))
// 		})
// 	}
// }
