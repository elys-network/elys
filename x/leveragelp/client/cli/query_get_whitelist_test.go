package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/x/leveragelp/client/cli"
	"github.com/elys-network/elys/x/leveragelp/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithWhitelistedObjects(t *testing.T, n int) (*network.Network, []string) {
	t.Helper()
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	state := types.GenesisState{}

	whitelistedAddrs := make([]string, 0)
	// Generate n random accounts with 1000000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, n, sdk.NewInt(1000000))

	cfg := network.DefaultConfig()
	for i := 0; i < n; i++ {
		whitelistedAddrs = append(whitelistedAddrs, addr[i].String())
		state.AddressWhitelist = append(state.AddressWhitelist, addr[i].String())
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), whitelistedAddrs
}

func TestGetWhitelistedAddresses(t *testing.T) {
	net, objs := networkWithWhitelistedObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdGetWhitelist(), args)
			require.NoError(t, err)
			var resp types.WhitelistResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Whitelist), step)
			require.Subset(t,
				objs,
				resp.Whitelist,
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdGetWhitelist(), args)
			require.NoError(t, err)
			var resp types.WhitelistResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Whitelist), step)
			require.Subset(t,
				objs,
				resp.Whitelist,
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdGetWhitelist(), args)
		require.NoError(t, err)
		var resp types.WhitelistResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			objs,
			resp.Whitelist,
		)
	})
}
