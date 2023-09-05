package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/accountedpool/client/cli"
	"github.com/elys-network/elys/x/accountedpool/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithAccountedPoolObjects(t *testing.T, n int) (*network.Network, []types.AccountedPool) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		accountedPool := types.AccountedPool{
			PoolId: (uint64)(i),
		}
		nullify.Fill(&accountedPool)
		state.AccountedPoolList = append(state.AccountedPoolList, accountedPool)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.AccountedPoolList
}

func TestShowAccountedPool(t *testing.T) {
	net, objs := networkWithAccountedPoolObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc    string
		idIndex uint64

		args []string
		err  error
		obj  types.AccountedPool
	}{
		{
			desc:    "found",
			idIndex: objs[0].PoolId,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			idIndex: (uint64)(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				fmt.Sprintf("%d", tc.idIndex),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowAccountedPool(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetAccountedPoolResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.AccountedPool)
				// total weight is not set in genesis state
				tc.obj.TotalWeight = resp.AccountedPool.TotalWeight
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.AccountedPool),
				)
			}
		})
	}
}

func TestListAccountedPool(t *testing.T) {
	net, objs := networkWithAccountedPoolObjects(t, 5)
	ctx := net.Validators[0].ClientCtx
	const stepSize = 2

	type RequestArgs struct {
		Next   []byte
		Offset uint64
		Limit  uint64
		Total  bool
	}

	request := func(args RequestArgs) []string {
		var requestArgs []string
		requestArgs = append(requestArgs, fmt.Sprintf("--%s=json", tmcli.OutputFlag))

		if args.Next == nil {
			requestArgs = append(requestArgs, fmt.Sprintf("--%s=%d", flags.FlagOffset, args.Offset))
		} else {
			requestArgs = append(requestArgs, fmt.Sprintf("--%s=%s", flags.FlagPageKey, args.Next))
		}

		requestArgs = append(requestArgs, fmt.Sprintf("--%s=%d", flags.FlagLimit, args.Limit))

		if args.Total {
			requestArgs = append(requestArgs, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return requestArgs
	}

	executeCmdAndCheck := func(t *testing.T, args RequestArgs) (types.QueryAllAccountedPoolResponse, error) {
		cmdArgs := request(args)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListAccountedPool(), cmdArgs)
		if err != nil {
			return types.QueryAllAccountedPoolResponse{}, err
		}

		var resp types.QueryAllAccountedPoolResponse
		err = net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp)
		return resp, err
	}

	t.Run("ByOffset", func(t *testing.T) {
		for i := 0; i < len(objs); i += stepSize {
			resp, err := executeCmdAndCheck(t, RequestArgs{Next: nil, Offset: uint64(i), Limit: uint64(stepSize), Total: false})
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AccountedPool), stepSize)

			for j, accountedPool := range resp.AccountedPool {
				objs[i+j].TotalWeight = accountedPool.TotalWeight
			}

			require.Subset(t, nullify.Fill(objs), nullify.Fill(resp.AccountedPool))
		}
	})

	t.Run("ByKey", func(t *testing.T) {
		var next []byte
		for i := 0; i < len(objs); i += stepSize {
			resp, err := executeCmdAndCheck(t, RequestArgs{Next: next, Offset: 0, Limit: uint64(stepSize), Total: false})
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AccountedPool), stepSize)

			for j, accountedPool := range resp.AccountedPool {
				objs[i+j].TotalWeight = accountedPool.TotalWeight
			}

			require.Subset(t, nullify.Fill(objs), nullify.Fill(resp.AccountedPool))
			next = resp.Pagination.NextKey
		}
	})

	t.Run("Total", func(t *testing.T) {
		resp, err := executeCmdAndCheck(t, RequestArgs{Next: nil, Offset: 0, Limit: uint64(len(objs)), Total: true})
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))

		for i, accountedPool := range resp.AccountedPool {
			objs[i].TotalWeight = accountedPool.TotalWeight
		}

		require.ElementsMatch(t, nullify.Fill(objs), nullify.Fill(resp.AccountedPool))
	})
}
