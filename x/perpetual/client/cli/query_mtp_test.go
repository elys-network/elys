package cli_test

import (
	"cosmossdk.io/math"
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/testutil/network"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/client/cli"
	"github.com/elys-network/elys/x/perpetual/types"
)

func networkWithMTPObjects(t *testing.T, n int) (*network.Network, []*types.MTP) {
	t.Helper()
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)
	state := types.GenesisState{}

	mtps := make([]*types.MTP, 0)
	// Generate n random accounts with 1000000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, n, math.NewInt(1000000))

	cfg := network.DefaultConfig()
	for i := 0; i < n; i++ {
		mtp := types.MTP{
			Address:                        addr[i].String(),
			CollateralAsset:                ptypes.BaseCurrency,
			TradingAsset:                   "ATOM",
			LiabilitiesAsset:               ptypes.BaseCurrency,
			CustodyAsset:                   "ATOM",
			Collateral:                     math.NewInt(0),
			Liabilities:                    math.NewInt(0),
			BorrowInterestPaidCollateral:   math.NewInt(0),
			BorrowInterestPaidCustody:      math.NewInt(0),
			BorrowInterestUnpaidCollateral: math.NewInt(0),
			Custody:                        math.NewInt(0),
			TakeProfitLiabilities:          math.NewInt(0),
			TakeProfitCustody:              math.NewInt(0),
			Leverage:                       math.LegacyNewDec(0),
			MtpHealth:                      math.LegacyNewDec(0),
			Position:                       types.Position_LONG,
			Id:                             (uint64)(i + 1),
			AmmPoolId:                      (uint64)(i + 1),
			ConsolidateLeverage:            math.LegacyZeroDec(),
			SumCollateral:                  math.ZeroInt(),
			TakeProfitPrice:                math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
			TakeProfitBorrowRate:           math.LegacyOneDec(),
			FundingFeePaidCollateral:       math.NewInt(0),
			FundingFeePaidCustody:          math.NewInt(0),
			FundingFeeReceivedCollateral:   math.NewInt(0),
			FundingFeeReceivedCustody:      math.NewInt(0),
			OpenPrice:                      math.LegacyNewDec(0),
			StopLossPrice:                  math.LegacyNewDec(0),
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
