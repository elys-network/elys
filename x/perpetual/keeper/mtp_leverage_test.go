package keeper_test

import (
	"cosmossdk.io/math"
	assetprofiletypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestGetEffectiveLeverage() {
	var mtp types.MTP
	mtp.TradingAsset = ptypes.ATOM

	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.BaseCurrency,
		Denom:       ptypes.BaseCurrency,
		Decimals:    6,
		DisplayName: "USDC",
	})
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
	})

	testCases := []struct {
		name                 string
		expectErrMsg         string
		result               math.LegacyDec
		prerequisiteFunction func()
	}{
		{
			"price not set",
			"asset info uatom not found",
			math.LegacyDec{},
			func() {
			},
		},
		{
			"LONG",
			"",
			math.LegacyMustNewDecFromStr("1.176470588235294117"),
			func() {
				suite.SetupCoinPrices()
				mtp.Position = types.Position_LONG
				mtp.Custody = math.NewInt(100)
				mtp.Liabilities = math.NewInt(75)
			},
		},
		{
			"SHORT",
			"",
			math.LegacyMustNewDecFromStr("1.666666666666666666"),
			func() {
				suite.SetupCoinPrices()
				mtp.Position = types.Position_SHORT
				mtp.Custody = math.NewInt(600)
				mtp.Liabilities = math.NewInt(75)
			},
		},
		{
			"SHORT, bot should liquidate before leverage goes negative",
			"",
			math.LegacyMustNewDecFromStr("-1.363636363636363636"),
			func() {
				suite.SetupCoinPrices()
				mtp.Position = types.Position_SHORT
				mtp.Custody = math.NewInt(100)
				mtp.Liabilities = math.NewInt(75)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			effectiveLeverage, err := suite.app.PerpetualKeeper.GetEffectiveLeverage(suite.ctx, mtp)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().Equal(tc.result, effectiveLeverage)
				suite.Require().NoError(err)
			}
		})
	}
}
