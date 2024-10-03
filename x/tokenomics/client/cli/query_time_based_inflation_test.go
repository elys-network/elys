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
	"github.com/elys-network/elys/x/tokenomics/client/cli"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func networkWithTimeBasedInflationObjects(t *testing.T, n int) (*network.Network, []types.TimeBasedInflation) {
	t.Helper()
	cfg := network.DefaultConfig(t.TempDir())
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		timeBasedInflation := types.TimeBasedInflation{
			StartBlockHeight: uint64(i),
			EndBlockHeight:   uint64(i),
			Inflation: &types.InflationEntry{
				LmRewards:         9999999,
				IcsStakingRewards: 9999999,
				CommunityFund:     9999999,
				StrategicReserve:  9999999,
				TeamTokensVested:  9999999,
			},
		}
		nullify.Fill(&timeBasedInflation)
		state.TimeBasedInflationList = append(state.TimeBasedInflationList, timeBasedInflation)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.TimeBasedInflationList
}

func TestShowTimeBasedInflation(t *testing.T) {
	net, objs := networkWithTimeBasedInflationObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc               string
		idStartBlockHeight uint64
		idEndBlockHeight   uint64

		args []string
		err  error
		obj  types.TimeBasedInflation
	}{
		{
			desc:               "found",
			idStartBlockHeight: objs[0].StartBlockHeight,
			idEndBlockHeight:   objs[0].EndBlockHeight,

			args: common,
			obj:  objs[0],
		},
		{
			desc:               "not found",
			idStartBlockHeight: 100000,
			idEndBlockHeight:   100000,

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				strconv.Itoa(int(tc.idStartBlockHeight)),
				strconv.Itoa(int(tc.idEndBlockHeight)),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowTimeBasedInflation(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetTimeBasedInflationResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.TimeBasedInflation)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.TimeBasedInflation),
				)
			}
		})
	}
}

func TestListTimeBasedInflation(t *testing.T) {
	net, objs := networkWithTimeBasedInflationObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTimeBasedInflation(), args)
			require.NoError(t, err)
			var resp types.QueryAllTimeBasedInflationResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.TimeBasedInflation), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.TimeBasedInflation),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTimeBasedInflation(), args)
			require.NoError(t, err)
			var resp types.QueryAllTimeBasedInflationResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.TimeBasedInflation), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.TimeBasedInflation),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTimeBasedInflation(), args)
		require.NoError(t, err)
		var resp types.QueryAllTimeBasedInflationResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.TimeBasedInflation),
		)
	})
}
