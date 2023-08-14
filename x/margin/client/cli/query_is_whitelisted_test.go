package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/x/margin/client/cli"
	"github.com/elys-network/elys/x/margin/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestIsWhitelisted(t *testing.T) {
	net, objs := networkWithWhitelistedObjects(t, 2)

	ctx := net.Validators[0].ClientCtx

	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc          string
		addr          string
		isWhitelisted bool
		args          []string
		err           error
		obj           string
	}{
		{
			desc:          "found",
			addr:          objs[0],
			isWhitelisted: true,
			args:          common,
			obj:           objs[0],
		},
		{
			desc:          "not found",
			addr:          "invalid address",
			isWhitelisted: false,
			args:          common,
			err:           status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.addr,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdIsWhitelisted(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.IsWhitelistedResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.obj, resp.Address)
				require.Equal(t, resp.IsWhitelisted, true)
			}
		})
	}
}
