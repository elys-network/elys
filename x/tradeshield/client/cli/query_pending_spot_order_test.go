package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func networkWithPendingSpotOrderObjects(t *testing.T, n int) (*network.Network, []types.SpotOrder) {
	t.Helper()
	cfg := network.DefaultConfig(t.TempDir())
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		pendingSpotOrder := types.SpotOrder{
			OrderType:    types.SpotOrderType_MARKETBUY, // Assuming a BUY order type
			OrderId:      uint64(i + 1),
			OwnerAddress: fmt.Sprintf("address%d", i+1),
			Status:       types.Status_PENDING, // Assuming a PENDING status
		}
		nullify.Fill(&pendingSpotOrder)
		state.PendingSpotOrderList = append(state.PendingSpotOrderList, pendingSpotOrder)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.PendingSpotOrderList
}

// TODO: Add tests for the CLI queries in query task
// func TestShowPendingSpotOrder(t *testing.T) {
// 	net, objs := networkWithPendingSpotOrderObjects(t, 2)

// 	ctx := net.Validators[0].ClientCtx
// 	common := []string{
// 		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
// 	}
// 	tests := []struct {
// 		desc string
// 		id   string
// 		args []string
// 		err  error
// 		obj  types.SpotOrder
// 	}{
// 		{
// 			desc: "found",
// 			id:   fmt.Sprintf("%d", objs[0].OrderId),
// 			args: common,
// 			obj:  objs[0],
// 		},
// 		{
// 			desc: "not found",
// 			id:   "not_found",
// 			args: common,
// 			err:  status.Error(codes.NotFound, "not found"),
// 		},
// 	}
// 	for _, tc := range tests {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			args := []string{tc.id}
// 			args = append(args, tc.args...)
// 			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowPendingSpotOrder(), args)
// 			if tc.err != nil {
// 				stat, ok := status.FromError(tc.err)
// 				require.True(t, ok)
// 				require.ErrorIs(t, stat.Err(), tc.err)
// 			} else {
// 				require.NoError(t, err)
// 				var resp types.QueryGetPendingSpotOrderResponse
// 				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 				require.NotNil(t, resp.PendingSpotOrder)
// 				require.Equal(t,
// 					nullify.Fill(&tc.obj),
// 					nullify.Fill(&resp.PendingSpotOrder),
// 				)
// 			}
// 		})
// 	}
// }

// func TestListPendingSpotOrder(t *testing.T) {
// 	net, objs := networkWithPendingSpotOrderObjects(t, 5)

// 	ctx := net.Validators[0].ClientCtx
// 	request := func(next []byte, offset, limit uint64, total bool) []string {
// 		args := []string{
// 			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
// 		}
// 		if next == nil {
// 			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
// 		} else {
// 			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
// 		}
// 		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
// 		if total {
// 			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
// 		}
// 		return args
// 	}
// 	t.Run("ByOffset", func(t *testing.T) {
// 		step := 2
// 		for i := 0; i < len(objs); i += step {
// 			args := request(nil, uint64(i), uint64(step), false)
// 			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPendingSpotOrder(), args)
// 			require.NoError(t, err)
// 			var resp types.QueryAllPendingSpotOrderResponse
// 			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 			require.LessOrEqual(t, len(resp.PendingSpotOrder), step)
// 			require.Subset(t,
// 				nullify.Fill(objs),
// 				nullify.Fill(resp.PendingSpotOrder),
// 			)
// 		}
// 	})
// 	t.Run("ByKey", func(t *testing.T) {
// 		step := 2
// 		var next []byte
// 		for i := 0; i < len(objs); i += step {
// 			args := request(next, 0, uint64(step), false)
// 			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPendingSpotOrder(), args)
// 			require.NoError(t, err)
// 			var resp types.QueryAllPendingSpotOrderResponse
// 			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 			require.LessOrEqual(t, len(resp.PendingSpotOrder), step)
// 			require.Subset(t,
// 				nullify.Fill(objs),
// 				nullify.Fill(resp.PendingSpotOrder),
// 			)
// 			next = resp.Pagination.NextKey
// 		}
// 	})
// 	t.Run("Total", func(t *testing.T) {
// 		args := request(nil, 0, uint64(len(objs)), true)
// 		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPendingSpotOrder(), args)
// 		require.NoError(t, err)
// 		var resp types.QueryAllPendingSpotOrderResponse
// 		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 		require.NoError(t, err)
// 		require.Equal(t, len(objs), int(resp.Pagination.Total))
// 		require.ElementsMatch(t,
// 			nullify.Fill(objs),
// 			nullify.Fill(resp.PendingSpotOrder),
// 		)
// 	})
// }
