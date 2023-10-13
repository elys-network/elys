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

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/testutil/network"
	"github.com/elys-network/elys/x/margin/client/cli"
	"github.com/elys-network/elys/x/margin/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithMTPObjects(t *testing.T, n int) (*network.Network, []*types.MTP) {
	t.Helper()
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	state := types.GenesisState{}

	mtps := make([]*types.MTP, 0)
	// Generate n random accounts with 1000000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, n, sdk.NewInt(1000000))

	cfg := network.DefaultConfig()
	for i := 0; i < n; i++ {
		mtp := types.MTP{
			Address:                   addr[i].String(),
			Collaterals:               []sdk.Coin{sdk.NewCoin(paramtypes.BaseCurrency, sdk.NewInt(0))},
			Liabilities:               sdk.NewInt(0),
			InterestPaidCollaterals:   []sdk.Coin{sdk.NewCoin(paramtypes.BaseCurrency, sdk.NewInt(0))},
			InterestPaidCustodies:     []sdk.Coin{sdk.NewCoin("ATOM", sdk.NewInt(0))},
			InterestUnpaidCollaterals: []sdk.Coin{sdk.NewCoin(paramtypes.BaseCurrency, sdk.NewInt(0))},
			Custodies:                 []sdk.Coin{sdk.NewCoin("ATOM", sdk.NewInt(0))},
			Leverages:                 []sdk.Dec{sdk.NewDec(0)},
			MtpHealth:                 sdk.NewDec(0),
			Position:                  types.Position_LONG,
			Id:                        (uint64)(i + 1),
			AmmPoolId:                 (uint64)(i + 1),
			ConsolidateLeverage:       sdk.ZeroDec(),
			SumCollateral:             sdk.ZeroInt(),
		}

		mtps = append(mtps, &mtp)
		state.MtpList = append(state.MtpList, mtp)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), mtps
}

func TestShowMTP(t *testing.T) {
	net, objs := networkWithMTPObjects(t, 2)

	ctx := net.Validators[0].ClientCtx

	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc    string
		addr    string
		idIndex uint64

		args []string
		err  error
		obj  *types.MTP
	}{
		{
			desc:    "found",
			addr:    objs[0].Address,
			idIndex: objs[0].Id,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			addr:    objs[0].Address,
			idIndex: (uint64)(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.addr,
				fmt.Sprintf("%d", tc.idIndex),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdMtp(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.MTPResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Mtp)
				require.Equal(t,
					tc.obj,
					resp.Mtp,
				)
			}
		})
	}
}
