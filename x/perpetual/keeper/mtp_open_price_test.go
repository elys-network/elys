package keeper_test

import (
	"cosmossdk.io/math"
	"errors"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestGetAndSetOpenPrice() {
	addr := suite.AddAccounts(1, nil)
	amount := math.NewInt(1000)

	testCases := []struct {
		name              string
		mtp               *types.MTP
		expectedOpenPrice math.LegacyDec
		expectedErr       error
		setup             func(mtp *types.MTP) // Helper to set up custody/collateral easily
	}{
		{
			name: "Error no asset price set",
			mtp: &types.MTP{
				TradingAsset:     ptypes.ATOM,
				LiabilitiesAsset: ptypes.BaseCurrency,
				CollateralAsset:  ptypes.BaseCurrency,
				Position:         types.Position_LONG,
				Liabilities:      math.NewInt(15_000_000),
				Custody:          math.NewInt(10_000_000),
				Collateral:       math.NewInt(2_000_000),
			},
			expectedErr: errors.New("denom price uusdc is zero"),
		},
		{
			name: "LONG position, collateral is trading asset, valid calculation",
			mtp: &types.MTP{
				TradingAsset:     ptypes.ATOM,
				CollateralAsset:  ptypes.ATOM,
				LiabilitiesAsset: ptypes.BaseCurrency,
				Position:         types.Position_LONG,
				Liabilities:      math.NewInt(15_000_000),
				Custody:          math.NewInt(10_000_000),
				Collateral:       math.NewInt(2_000_000),
			},
			expectedOpenPrice: math.LegacyMustNewDecFromStr("1.875"), // liabilities / (custody - collateral)
			expectedErr:       nil,
			setup: func(mtp *types.MTP) {
				_, _ = suite.ResetAndSetSuite(addr, true, amount.MulRaw(1000), math.NewInt(2))
			},
		},
		{
			name: "ERROR LONG position, collateral is trading asset, denominator (uusdcPerBaseAsset) is zero",
			mtp: &types.MTP{
				TradingAsset:     ptypes.ATOM,
				CollateralAsset:  ptypes.ATOM,
				LiabilitiesAsset: ptypes.BaseCurrency,
				Position:         types.Position_LONG,
				Liabilities:      math.NewInt(1000_000_000),
				Custody:          math.NewInt(10_000_000),
				Collateral:       math.NewInt(10_000_000),
			},
			expectedOpenPrice: math.LegacyZeroDec(), // Expect 0 as denominator becomes 0
			expectedErr:       errors.New("(custody - collateral) is zero while calculating open price"),
		},
		{
			name: "LONG position, collateral is different asset, valid calculation",
			mtp: &types.MTP{
				TradingAsset:     ptypes.ATOM,
				CollateralAsset:  ptypes.BaseCurrency,
				LiabilitiesAsset: ptypes.BaseCurrency,
				Position:         types.Position_LONG,
				Collateral:       math.NewInt(50_000_000),
				Liabilities:      math.NewInt(10_000_000),
				Custody:          math.NewInt(20_000_000),
			},
			expectedOpenPrice: math.LegacyMustNewDecFromStr("3.0"), // (collateral + liabilities) / custody
			expectedErr:       nil,
		},
		{
			name: "ERROR LONG position, collateral is different asset, denominator (custody) is zero",
			mtp: &types.MTP{
				TradingAsset:     ptypes.ATOM,
				CollateralAsset:  ptypes.BaseCurrency,
				LiabilitiesAsset: ptypes.BaseCurrency,
				Position:         types.Position_LONG,
				Collateral:       math.NewInt(500_000_000),
				Liabilities:      math.NewInt(1000_000_000),
				Custody:          math.NewInt(0), // Custody is 0
			},
			expectedOpenPrice: math.LegacyZeroDec(), // Expect 0 if denominator is 0
			expectedErr:       errors.New("custody is zero while calculating open price"),
		},
		{
			name: "ERROR SHORT position, liabilities are zero",
			mtp: &types.MTP{
				TradingAsset:     ptypes.ATOM,
				LiabilitiesAsset: ptypes.ATOM,
				CollateralAsset:  ptypes.BaseCurrency,
				Position:         types.Position_SHORT,
				Liabilities:      math.NewInt(0), // Liabilities are zero
				Custody:          math.NewInt(1000_000_000),
				Collateral:       math.NewInt(100_000_000),
			},
			expectedOpenPrice: math.LegacyZeroDec(), // Expect 0
			expectedErr:       nil,
		},
		{
			name: "SHORT position, valid calculation",
			mtp: &types.MTP{
				TradingAsset:     ptypes.ATOM,
				LiabilitiesAsset: ptypes.ATOM,
				CollateralAsset:  ptypes.BaseCurrency,
				Position:         types.Position_SHORT,
				Custody:          math.NewInt(40_000_000),
				Collateral:       math.NewInt(20_000_000),
				Liabilities:      math.NewInt(10_000_000), // (custody - collateral) / liabilities
			},
			expectedOpenPrice: math.LegacyMustNewDecFromStr("2.0"),
			expectedErr:       nil,
		},
		{
			name: "LONG ETH position, collateral is trading asset, valid calculation",
			mtp: &types.MTP{
				TradingAsset:     "wei",
				CollateralAsset:  "wei",
				LiabilitiesAsset: ptypes.BaseCurrency,
				Position:         types.Position_LONG,
				Liabilities:      math.NewInt(6000_000_000),                                         // 6000 USDC
				Custody:          math.LegacyMustNewDecFromStr("6000000000000000000").TruncateInt(), // 6 eth
				Collateral:       math.LegacyMustNewDecFromStr("2000000000000000000").TruncateInt(), // 2 ETH
			},
			expectedOpenPrice: math.LegacyMustNewDecFromStr("1500"), // liabilities / (custody - collateral)
			expectedErr:       nil,
		},
		{
			name: "LONG ETH position, collateral is different asset, valid calculation",
			mtp: &types.MTP{
				TradingAsset:     "wei",
				CollateralAsset:  ptypes.BaseCurrency,
				LiabilitiesAsset: ptypes.BaseCurrency,
				Position:         types.Position_LONG,
				Liabilities:      math.NewInt(4000_000_000),                                         // 4000 usdc
				Custody:          math.LegacyMustNewDecFromStr("4000000000000000000").TruncateInt(), // 4 ETH
				Collateral:       math.NewInt(2000_000_000),                                         // 2000 usdc
			},
			expectedOpenPrice: math.LegacyMustNewDecFromStr("1500"), // (collateral + liabilities) / custody
			expectedErr:       nil,
		},
		{
			name: "SHORT ETH position, valid calculation",
			mtp: &types.MTP{
				TradingAsset:     "wei",
				LiabilitiesAsset: "wei",
				CollateralAsset:  ptypes.BaseCurrency,
				Position:         types.Position_SHORT,
				Custody:          math.NewInt(5000_000_000), // 6000
				Collateral:       math.NewInt(2000_000_000),
				Liabilities:      math.LegacyMustNewDecFromStr("2000000000000000000").TruncateInt(), // 2 ETH
			},
			expectedOpenPrice: math.LegacyMustNewDecFromStr("1500"), // (custody - collateral) / liabilities
			expectedErr:       nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			if tc.setup != nil {
				tc.setup(tc.mtp)
			}

			err := suite.app.PerpetualKeeper.GetAndSetOpenPrice(suite.ctx, tc.mtp)

			if tc.expectedErr != nil {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErr.Error())
			} else {
				suite.Require().NoError(err)
				suite.Require().True(tc.expectedOpenPrice.Equal(tc.mtp.OpenPrice), "Expected OpenPrice %s, got %s", tc.expectedOpenPrice.String(), tc.mtp.OpenPrice.String())
			}
		})
	}
}
