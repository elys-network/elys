package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/elys-network/elys/testutil/network"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/client/cli"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func networkWithMTPObjects(t *testing.T, n int) (*network.Network, []*types.MtpAndPrice) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}

	mtps := make([]*types.MtpAndPrice, 0)

	for i := 0; i < n; i++ {
		mtp := types.MtpAndPrice{
			Mtp: &types.MTP{
				Address:                       authtypes.NewModuleAddress("test").String(),
				CollateralAsset:               ptypes.BaseCurrency,
				TradingAsset:                  "ATOM",
				LiabilitiesAsset:              ptypes.BaseCurrency,
				CustodyAsset:                  "ATOM",
				Collateral:                    sdk.NewInt(0),
				Liabilities:                   sdk.NewInt(0),
				BorrowInterestPaidCustody:     sdk.NewInt(0),
				BorrowInterestUnpaidLiability: sdk.NewInt(0),
				Custody:                       sdk.NewInt(0),
				TakeProfitLiabilities:         sdk.NewInt(0),
				TakeProfitCustody:             sdk.NewInt(0),
				MtpHealth:                     sdk.NewDec(0),
				Position:                      types.Position_LONG,
				Id:                            (uint64)(i + 1),
				AmmPoolId:                     (uint64)(i + 1),
				TakeProfitPrice:               types.TakeProfitPriceDefault,
				TakeProfitBorrowFactor:        sdk.OneDec(),
				FundingFeePaidCustody:         sdk.NewInt(0),
				FundingFeeReceivedCustody:     sdk.NewInt(0),
				OpenPrice:                     sdk.NewDec(0),
				StopLossPrice:                 sdk.NewDec(0),
			},
			TradingAssetPrice: sdk.ZeroDec(),
			Pnl:               sdk.NewCoin("USDC", sdk.NewInt(0)),
			LiquidationPrice:  sdk.ZeroDec(),
			Fees: &types.Fees{
				TotalFeesBaseCurrency:            sdk.NewInt(0),
				BorrowInterestFeesLiabilityAsset: sdk.NewInt(0),
				BorrowInterestFeesBaseCurrency:   sdk.NewInt(0),
				FundingFeesLiquidityAsset:        sdk.NewInt(0),
				FundingFeesBaseCurrency:          sdk.NewInt(0),
			},
		}

		mtps = append(mtps, &mtp)
		state.MtpList = append(state.MtpList, *mtp.Mtp)
	}

	state.Params = types.NewParams()
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	// Set oracle price and info
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	stateOracle := oracletypes.GenesisState{}
	stateOracle.Prices = []oracletypes.Price{
		{
			Asset:       "ATOM",
			Price:       sdk.NewDec(4),
			Source:      oracletypes.BAND,
			Provider:    provider.String(),
			Timestamp:   uint64(1729430615),
			BlockHeight: uint64(1),
		},
		{
			Asset:       "ATOM",
			Price:       sdk.NewDec(4),
			Source:      oracletypes.ELYS,
			Provider:    provider.String(),
			Timestamp:   uint64(1729430615),
			BlockHeight: uint64(1),
		},
		{
			Asset:       "USDC",
			Price:       sdk.NewDec(1),
			Source:      oracletypes.BAND,
			Provider:    provider.String(),
			Timestamp:   uint64(1729430615),
			BlockHeight: uint64(1),
		},
		{
			Asset:       "USDC",
			Price:       sdk.NewDec(1),
			Source:      oracletypes.ELYS,
			Provider:    provider.String(),
			Timestamp:   uint64(1729430615),
			BlockHeight: uint64(1),
		},
		{
			Asset:       "BTC",
			Price:       sdk.NewDec(60000),
			Source:      oracletypes.ELYS,
			Provider:    provider.String(),
			Timestamp:   uint64(1729430615),
			BlockHeight: uint64(1),
		},
	}
	stateOracle.Params = oracletypes.DefaultParams()
	stateOracle.PortId = "portid"
	stateOracle.AssetInfos = []oracletypes.AssetInfo{
		{
			Denom:      "ATOM",
			Display:    "ATOM",
			Decimal:    6,
			BandTicker: "ATOM",
		},
		{
			Denom:      "USDC",
			Display:    "USDC",
			Decimal:    6,
			BandTicker: "USDC",
		},
		{
			Denom:      "BTC",
			Display:    "BTC",
			Decimal:    6,
			BandTicker: "BTC",
		},
	}

	bufO, err := cfg.Codec.MarshalJSON(&stateOracle)
	require.NoError(t, err)
	cfg.GenesisState[oracletypes.ModuleName] = bufO
	return network.New(t, cfg), mtps
}

func (s *CLITestSuite) TestShowMTP() {
	cmd := cli.CmdMtp()
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
				bz, _ := s.encCfg.Codec.Marshal(&types.MTPResponse{
					Mtp: &types.MtpAndPrice{},
				})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				authtypes.NewModuleAddress("test").String(),
				"1",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.MTPResponse{},
			false,
		},
		{
			"invalid query",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.MTPResponse{
					Mtp: &types.MtpAndPrice{},
				})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				authtypes.NewModuleAddress("test").String(),
				"-1",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			&types.MTPResponse{},
			true,
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

//func TestShowMTP(t *testing.T) {
//	net, objs := networkWithMTPObjects(t, 2)
//
//	ctx := net.Validators[0].ClientCtx
//
//	common := []string{
//		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
//	}
//	tests := []struct {
//		desc    string
//		addr    string
//		idIndex uint64
//
//		args []string
//		err  error
//		obj  *types.MtpAndPrice
//	}{
//		{
//			desc:    "found",
//			addr:    objs[0].Mtp.Address,
//			idIndex: objs[0].Mtp.Id,
//
//			args: common,
//			obj:  objs[0],
//		},
//		{
//			desc:    "not found",
//			addr:    objs[0].Mtp.Address,
//			idIndex: (uint64)(100000),
//
//			args: common,
//			err:  status.Error(codes.NotFound, "not found"),
//		},
//	}
//	for _, tc := range tests {
//		t.Run(tc.desc, func(t *testing.T) {
//			args := []string{
//				tc.addr,
//				fmt.Sprintf("%d", tc.idIndex),
//			}
//			args = append(args, tc.args...)
//			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdMtp(), args)
//			if tc.err != nil {
//				stat, ok := status.FromError(tc.err)
//				require.True(t, ok)
//				require.ErrorIs(t, stat.Err(), tc.err)
//			} else {
//				require.NoError(t, err)
//				var resp types.MTPResponse
//				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
//				require.NotNil(t, resp.Mtp)
//				require.Equal(t,
//					tc.obj,
//					resp.Mtp,
//				)
//			}
//		})
//	}
//}
