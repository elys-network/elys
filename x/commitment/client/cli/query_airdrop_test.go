package cli_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/x/commitment/client/cli"
	"github.com/elys-network/elys/x/commitment/types"
)

const validAddress = "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"

func TestCmdAirdrop(t *testing.T) {
	net := NetworkWithAirdropObjects(t)
	ctx := net.Validators[0].ClientCtx
	commonArgs := []string{
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	}

	t.Run("Show airdrop for valid address", func(t *testing.T) {
		args := []string{
			validAddress,
		}
		args = append(args, commonArgs...)

		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdAirdrop(), args)
		require.NoError(t, err)

		var resp types.QueryAirDropResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NotNil(t, resp)
		// Add more specific assertions about the response content if needed
		require.Equal(t, resp, types.QueryAirDropResponse{
			AtomStaking: math.NewInt(1000000000000000000),
			Cadet:       math.NewInt(1000000000000000000),
			NftHolder:   math.NewInt(1000000000000000000),
			Governor:    math.NewInt(1000000000000000000),
		})
	})

	t.Run("Show airdrop for invalid address", func(t *testing.T) {
		invalidAddress := "invalid"
		args := []string{
			invalidAddress,
		}
		args = append(args, commonArgs...)

		_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdAirdrop(), args)
		require.Error(t, err)
	})
}

// NetworkWithAirdropObjects creates a network with pre-configured airdrop data
func NetworkWithAirdropObjects(t *testing.T) *network.Network {
	t.Helper()
	cfg := network.DefaultConfig(t.TempDir())

	state := types.GenesisState{}

	// Modify state to include airdrop data
	state.AtomStakers = []*types.AtomStaker{
		{
			Address: validAddress,
			Amount:  math.NewInt(1000000000000000000),
		},
	}
	state.Cadets = []*types.Cadet{
		{
			Address: validAddress,
			Amount:  math.NewInt(1000000000000000000),
		},
	}
	state.NftHolders = []*types.NftHolder{
		{
			Address: validAddress,
			Amount:  math.NewInt(1000000000000000000),
		},
	}
	state.Governors = []*types.Governor{
		{
			Address: validAddress,
			Amount:  math.NewInt(1000000000000000000),
		},
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg)
}
