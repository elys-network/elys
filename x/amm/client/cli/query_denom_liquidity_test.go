package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	assetprofilemoduletypes "github.com/elys-network/elys/x/assetprofile/types"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/amm/client/cli"
	"github.com/elys-network/elys/x/amm/types"
)

func networkWithDenomLiquidityObjects(t *testing.T, n int) (*network.Network, []types.DenomLiquidity) {
	t.Helper()
	cfg := network.DefaultConfig(t.TempDir())
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		denomLiquidity := types.DenomLiquidity{
			Denom:     strconv.Itoa(i),
			Liquidity: math.ZeroInt(),
		}
		nullify.Fill(&denomLiquidity)
		state.DenomLiquidityList = append(state.DenomLiquidityList, denomLiquidity)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	assetProfileGenesisState := assetprofilemoduletypes.DefaultGenesis()
	usdcEntry := assetprofilemoduletypes.Entry{
		BaseDenom:                "uusdc",
		Decimals:                 6,
		Denom:                    "uusdc",
		Path:                     "",
		IbcChannelId:             "",
		IbcCounterpartyChannelId: "",
		DisplayName:              "",
		DisplaySymbol:            "",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "",
		TransferLimit:            "",
		Permissions:              nil,
		UnitDenom:                "",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	}
	assetProfileGenesisState.EntryList = []assetprofilemoduletypes.Entry{usdcEntry}
	buf, err = cfg.Codec.MarshalJSON(assetProfileGenesisState)
	require.NoError(t, err)
	cfg.GenesisState[assetprofilemoduletypes.ModuleName] = buf

	return network.New(t, cfg), state.DenomLiquidityList
}

func TestShowDenomLiquidity(t *testing.T) {
	net, objs := networkWithDenomLiquidityObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc    string
		idDenom string

		args []string
		err  error
		obj  types.DenomLiquidity
	}{
		{
			desc:    "found",
			idDenom: objs[0].Denom,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			idDenom: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idDenom,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowDenomLiquidity(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetDenomLiquidityResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.DenomLiquidity)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.DenomLiquidity),
				)
			}
		})
	}
}

func TestListDenomLiquidity(t *testing.T) {
	net, objs := networkWithDenomLiquidityObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDenomLiquidity(), args)
			require.NoError(t, err)
			var resp types.QueryAllDenomLiquidityResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.DenomLiquidity), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.DenomLiquidity),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDenomLiquidity(), args)
			require.NoError(t, err)
			var resp types.QueryAllDenomLiquidityResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.DenomLiquidity), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.DenomLiquidity),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDenomLiquidity(), args)
		require.NoError(t, err)
		var resp types.QueryAllDenomLiquidityResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.DenomLiquidity),
		)
	})
}
