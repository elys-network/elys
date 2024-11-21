package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"io"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/gogoproto/proto"
	"github.com/elys-network/elys/x/masterchef/client/cli"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (s *CLITestSuite) TestApr() {
	cmd := cli.CmdApr()
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
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryAprResponse{
					//Apr: sdkmath.LegacyZeroDec(),
				})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"2",
				"ueden",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.QueryAprResponse{
				//Apr: sdkmath.LegacyZeroDec(),
			},
			false,
		},
		// {
		// 	"invalid query",
		// 	func() client.Context {
		// 		bz, _ := s.encCfg.Codec.Marshal(&types.PositionsByPoolResponse{
		// 			Mtps: make([]*types.MtpAndPrice, 0),
		// 		})
		// 		c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
		// 			Value: bz,
		// 		})
		// 		return s.baseCtx.WithClient(c)
		// 	},
		// 	[]string{
		// 		"-1",
		// 		fmt.Sprintf("--%s=json", flags.FlagOutput),
		// 	},
		// 	&types.PositionsByPoolResponse{},
		// 	true,
		// },
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
