package cli_test

import (
	"encoding/json"
	"os"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/elys-network/elys/v4/x/tradeshield/client/cli"
	"github.com/stretchr/testify/require"
)

func TestCancelSpotOrders(t *testing.T) {
	net := setupNetwork(t)
	ctx := net.Validators[0].ClientCtx
	val := net.Validators[0]

	tmpFile, err := os.CreateTemp("", "ids.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	validIds := []uint64{1}
	validJson, err := json.Marshal(validIds)
	require.NoError(t, err)
	_, err = tmpFile.Write(validJson)
	require.NoError(t, err)
	tmpFile.Close()

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...

	args := []string{
		tmpFile.Name(),
		"--from=" + val.Address.String(),
		"-y",
	}

	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdCancelSpotOrders(), args)
	require.NoError(t, err)
}
