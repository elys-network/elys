package cli_test

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"io"
	"strconv"
	"testing"

	"context"
	abci "github.com/cometbft/cometbft/abci/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/gogoproto/proto"
	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/assetprofile/client/cli"
	"github.com/elys-network/elys/x/assetprofile/types"
	"github.com/stretchr/testify/require"
)

var usdcEntry = types.Entry{
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

func networkWithEntryObjects(t *testing.T, n int) (*network.Network, []types.Entry) {
	t.Helper()
	cfg := network.DefaultConfig(t.TempDir())
	assetProfileGenesisState := types.DefaultGenesis()

	for i := 0; i < n; i++ {
		entry := types.Entry{
			BaseDenom: strconv.Itoa(i),
		}
		nullify.Fill(&entry)
		assetProfileGenesisState.EntryList = append(assetProfileGenesisState.EntryList, entry)
	}
	assetProfileGenesisState.EntryList = append(assetProfileGenesisState.EntryList, usdcEntry)
	buf, err := cfg.Codec.MarshalJSON(assetProfileGenesisState)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), assetProfileGenesisState.EntryList
}

func (s *CLITestSuite) TestShowEntry() {
	cmd := cli.CmdShowEntry()
	cmd.SetOutput(io.Discard)

	testCases := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"valid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryEntryResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"uusdc",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryEntryResponse{},
			false,
		},
		{
			"invalid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryEntryResponse{})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"-1",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryEntryResponse{},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var outBuf bytes.Buffer

			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd.SetContext(ctx)
			cmd.SetArgs(tc.args)

			s.Require().NoError(client.SetCmdClientContextHandler(clientCtx, cmd))

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(outBuf.Bytes(), tc.expectResult))
				s.Require().NoError(err)
			}
		})
	}
}

func TestListEntry(t *testing.T) {
	net, objs := networkWithEntryObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListEntry(), args)
			require.NoError(t, err)
			var resp types.QueryAllEntryResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Entry), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Entry),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListEntry(), args)
			require.NoError(t, err)
			var resp types.QueryAllEntryResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Entry), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Entry),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListEntry(), args)
		require.NoError(t, err)
		var resp types.QueryAllEntryResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.Entry),
		)
	})
}
