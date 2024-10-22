package cli_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/testutil/network"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/client/cli"
	"github.com/elys-network/elys/x/perpetual/types"
)

func networkWithMTPObjects(t *testing.T, n int) (*network.Network, []*types.MtpAndPrice) {
	t.Helper()
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	state := types.GenesisState{}
	mtps := make([]*types.MtpAndPrice, 0)

	simapp.SetStakingParam(app, ctx)
	simapp.SetPerpetualParams(app, ctx)

	// Generate n random accounts with 1000000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, n, math.NewInt(1000000))

	cfg := network.DefaultConfig(t.TempDir())
	for i := 0; i < n; i++ {
		mtp := types.MtpAndPrice{
			Mtp: &types.MTP{
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
				MtpHealth:                      math.LegacyNewDec(0),
				Position:                       types.Position_LONG,
				Id:                             (uint64)(i + 1),
				AmmPoolId:                      (uint64)(i + 1),
				TakeProfitPrice:                math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
				TakeProfitBorrowRate:           math.LegacyOneDec(),
				FundingFeePaidCollateral:       math.NewInt(0),
				FundingFeePaidCustody:          math.NewInt(0),
				FundingFeeReceivedCollateral:   math.NewInt(0),
				FundingFeeReceivedCustody:      math.NewInt(0),
				OpenPrice:                      math.LegacyNewDec(0),
				StopLossPrice:                  math.LegacyNewDec(0),
			},
			TradingAssetPrice: math.LegacyZeroDec(),
			Pnl:               math.LegacyZeroDec(),
		}

		mtps = append(mtps, &mtp)
		state.MtpList = append(state.MtpList, *mtp.Mtp)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	// Set oracle price and info
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	stateOracle := oracletypes.GenesisState{}
	stateOracle.Prices = append(stateOracle.Prices, oracletypes.Price{
		Asset:       "ATOM",
		Price:       math.LegacyNewDec(4),
		Source:      oracletypes.BAND,
		Provider:    provider.String(),
		Timestamp:   uint64(ctx.BlockTime().Unix()),
		BlockHeight: uint64(ctx.BlockHeight()),
	})
	stateOracle.Params = oracletypes.DefaultParams()
	stateOracle.PortId = "portid"
	stateOracle.AssetInfos = append(stateOracle.AssetInfos, oracletypes.AssetInfo{
		Denom:      "ATOM",
		Display:    "ATOM",
		Decimal:    6,
		BandTicker: "ATOM",
	})

	bufO, err := cfg.Codec.MarshalJSON(&stateOracle)
	require.NoError(t, err)
	cfg.GenesisState[oracletypes.ModuleName] = bufO
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
		obj  *types.MtpAndPrice
	}{
		{
			desc:    "found",
			addr:    objs[0].Mtp.Address,
			idIndex: objs[0].Mtp.Id,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			addr:    objs[0].Mtp.Address,
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
