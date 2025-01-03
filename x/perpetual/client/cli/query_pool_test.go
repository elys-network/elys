package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"cosmossdk.io/math"
	assetprofilemoduletypes "github.com/elys-network/elys/x/assetprofile/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/testutil/nullify"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/client/cli"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func networkWithPoolObjects(t *testing.T, n int) (*network.Network, []types.Pool) {
	t.Helper()
	cfg := network.DefaultConfig(t.TempDir())
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		ammPool := ammtypes.Pool{
			PoolId: uint64(i),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token: sdk.Coin{
						Denom:  "testAsset",
						Amount: math.NewInt(100),
					},
				},
			},
		}
		pool := types.NewPool(ammPool, math.LegacyMustNewDecFromStr("5.5"))
		nullify.Fill(&pool)
		state.PoolList = append(state.PoolList, pool)
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
	return network.New(t, cfg), state.PoolList
}

func (s *CLITestSuite) TestShowPool() {
	cmd := cli.CmdShowPool()
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
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryGetPoolResponse{
					Pool: types.PoolResponse{},
				})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"1",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryGetPoolResponse{},
			false,
		},
		{
			"invalid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryGetPoolResponse{
					Pool: types.PoolResponse{},
				})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"-1",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryGetPoolResponse{},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

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

//func TestShowPool(t *testing.T) {
//	net, objspool := networkWithPoolObjects(t, 2)
//
//	objs := make([]types.PoolResponse, len(objspool))
//
//	for k, v := range objspool {
//		objs[k].AmmPoolId = v.AmmPoolId
//		objs[k].BorrowInterestRate = v.BorrowInterestRate
//		objs[k].Health = v.Health
//		objs[k].FundingRate = v.FundingRate
//		objs[k].LastHeightBorrowInterestRateComputed = v.LastHeightBorrowInterestRateComputed
//		objs[k].PoolAssetsLong = v.PoolAssetsLong
//		objs[k].PoolAssetsShort = v.PoolAssetsShort
//		objs[k].NetOpenInterest = math.ZeroInt()
//	}
//
//	ctx := net.Validators[0].ClientCtx
//	common := []string{
//		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
//	}
//	tests := []struct {
//		desc    string
//		idIndex uint64
//
//		args []string
//		err  error
//		obj  types.PoolResponse
//	}{
//		{
//			desc:    "found",
//			idIndex: objs[0].AmmPoolId,
//
//			args: common,
//			obj:  objs[0],
//		},
//		{
//			desc:    "not found",
//			idIndex: (uint64)(100000),
//
//			args: common,
//			err:  status.Error(codes.NotFound, "not found"),
//		},
//	}
//	for _, tc := range tests {
//		t.Run(tc.desc, func(t *testing.T) {
//			args := []string{
//				fmt.Sprintf("%d", tc.idIndex),
//			}
//			args = append(args, tc.args...)
//			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowPool(), args)
//			if tc.err != nil {
//				stat, ok := status.FromError(tc.err)
//				require.True(t, ok)
//				require.ErrorIs(t, stat.Err(), tc.err)
//			} else {
//				require.NoError(t, err)
//				var resp types.QueryGetPoolResponse
//				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
//				require.NotNil(t, resp.Pool)
//				require.Equal(t,
//					nullify.Fill(&tc.obj),
//					nullify.Fill(&resp.Pool),
//				)
//			}
//		})
//	}
//}

/*
func TestListPool(t *testing.T) {
	net, objs := networkWithPoolObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPool(), args)
			require.NoError(t, err)
			var resp types.QueryAllPoolResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Pool), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Pool),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPool(), args)
			require.NoError(t, err)
			var resp types.QueryAllPoolResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Pool), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Pool),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPool(), args)
		require.NoError(t, err)
		var resp types.QueryAllPoolResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.Pool),
		)
	})
}
*/
