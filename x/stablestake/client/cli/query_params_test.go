package cli_test

import (
	"bytes"
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/elys-network/elys/x/stablestake/client/cli"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/gogo/protobuf/proto"
)

func (s *CLITestSuite) TestCmdQueryParams() {
	cmd := cli.CmdQueryParams()

	testCases := []struct {
		name         string
		args         []string
		ctxGen       func() client.Context
		expectErr    bool
		expectResult proto.Message
	}{
		{
			name: "successful query",
			args: []string{
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			ctxGen: func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryParamsResponse{
					Params: types.DefaultParams(),
				})
				c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			expectErr: false,
			expectResult: &types.QueryParamsResponse{
				Params: types.DefaultParams(),
			},
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
