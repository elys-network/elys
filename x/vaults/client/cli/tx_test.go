package cli_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/v6/testutil/network"
	"github.com/elys-network/elys/v6/x/vaults/client/cli"
)

func setupNetwork(t *testing.T) *network.Network {
	cfg := network.DefaultConfig(t.TempDir())
	cfg.NumValidators = 1
	return network.New(t, cfg)
}

func TestCmdPerformActionJoinPool(t *testing.T) {
	net := setupNetwork(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	common := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, math.NewInt(10))).String()),
	}

	testCases := []struct {
		desc    string
		vaultId string
		poolId  string
		share   string
		maxIn   string
		expErr  bool
	}{
		{
			desc:    "valid",
			vaultId: "1",
			poolId:  "1",
			share:   "1000000",
			maxIn:   "1000000uatom,1000000uusdc",
			expErr:  false,
		},
		{
			desc:    "invalid vault ID",
			vaultId: "invalid",
			poolId:  "1",
			share:   "1000000",
			maxIn:   "1000000uatom,1000000uusdc",
			expErr:  true,
		},
		{
			desc:    "invalid pool ID",
			vaultId: "1",
			poolId:  "invalid",
			share:   "1000000",
			maxIn:   "1000000uatom,1000000uusdc",
			expErr:  true,
		},
		{
			desc:    "invalid share amount",
			vaultId: "1",
			poolId:  "1",
			share:   "invalid",
			maxIn:   "1000000uatom,1000000uusdc",
			expErr:  true,
		},
		{
			desc:    "invalid max amounts in",
			vaultId: "1",
			poolId:  "1",
			share:   "1000000",
			maxIn:   "invalid",
			expErr:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.vaultId,
				tc.poolId,
				tc.share,
				tc.maxIn,
			}
			args = append(args, common...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPerformActionJoinPool(), args)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, uint32(0), resp.Code)
			}
		})
	}
}

func TestCmdPerformActionExitPool(t *testing.T) {
	net := setupNetwork(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	common := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, math.NewInt(10))).String()),
	}

	testCases := []struct {
		desc    string
		vaultId string
		poolId  string
		share   string
		minOut  string
		denom   string
		expErr  bool
	}{
		{
			desc:    "valid",
			vaultId: "1",
			poolId:  "1",
			share:   "1000000",
			minOut:  "1000000uatom,1000000uusdc",
			denom:   "uatom",
			expErr:  false,
		},
		{
			desc:    "invalid vault ID",
			vaultId: "invalid",
			poolId:  "1",
			share:   "1000000",
			minOut:  "1000000uatom,1000000uusdc",
			denom:   "uatom",
			expErr:  true,
		},
		{
			desc:    "invalid pool ID",
			vaultId: "1",
			poolId:  "invalid",
			share:   "1000000",
			minOut:  "1000000uatom,1000000uusdc",
			denom:   "uatom",
			expErr:  true,
		},
		{
			desc:    "invalid share amount",
			vaultId: "1",
			poolId:  "1",
			share:   "invalid",
			minOut:  "1000000uatom,1000000uusdc",
			denom:   "uatom",
			expErr:  true,
		},
		{
			desc:    "invalid min amounts out",
			vaultId: "1",
			poolId:  "1",
			share:   "1000000",
			minOut:  "invalid",
			denom:   "uatom",
			expErr:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.vaultId,
				tc.poolId,
				tc.share,
				tc.minOut,
				tc.denom,
			}
			args = append(args, common...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPerformActionExitPool(), args)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, uint32(0), resp.Code)
			}
		})
	}
}

func TestCmdPerformActionSwapByDenom(t *testing.T) {
	net := setupNetwork(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	common := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, math.NewInt(10))).String()),
	}

	testCases := []struct {
		desc     string
		vaultId  string
		amount   string
		minAmt   string
		maxAmt   string
		denomIn  string
		denomOut string
		expErr   bool
	}{
		{
			desc:     "valid",
			vaultId:  "1",
			amount:   "1000000uatom",
			minAmt:   "900000uusdc",
			maxAmt:   "1100000uusdc",
			denomIn:  "uatom",
			denomOut: "uusdc",
			expErr:   false,
		},
		{
			desc:     "invalid vault ID",
			vaultId:  "invalid",
			amount:   "1000000uatom",
			minAmt:   "900000uusdc",
			maxAmt:   "1100000uusdc",
			denomIn:  "uatom",
			denomOut: "uusdc",
			expErr:   true,
		},
		{
			desc:     "invalid amount",
			vaultId:  "1",
			amount:   "invalid",
			minAmt:   "900000uusdc",
			maxAmt:   "1100000uusdc",
			denomIn:  "uatom",
			denomOut: "uusdc",
			expErr:   true,
		},
		{
			desc:     "invalid min amount",
			vaultId:  "1",
			amount:   "1000000uatom",
			minAmt:   "invalid",
			maxAmt:   "1100000uusdc",
			denomIn:  "uatom",
			denomOut: "uusdc",
			expErr:   true,
		},
		{
			desc:     "invalid max amount",
			vaultId:  "1",
			amount:   "1000000uatom",
			minAmt:   "900000uusdc",
			maxAmt:   "invalid",
			denomIn:  "uatom",
			denomOut: "uusdc",
			expErr:   true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.vaultId,
				tc.amount,
				tc.minAmt,
				tc.maxAmt,
				tc.denomIn,
				tc.denomOut,
			}
			args = append(args, common...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPerformActionSwapByDenom(), args)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, uint32(0), resp.Code)
			}
		})
	}
}
