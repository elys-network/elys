package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/x/margin/client/cli"
	"github.com/elys-network/elys/x/margin/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestGetStatus(t *testing.T) {
	net, _ := networkWithMTPObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdGetStatus(), common)
	require.NoError(t, err)
	var resp types.StatusResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))

	// Gensis MTP count should be 5
	require.Equal(t, resp.OpenMtpCount, (uint64)(5))
	require.Equal(t, resp.LifetimeMtpCount, (uint64)(5))
}
