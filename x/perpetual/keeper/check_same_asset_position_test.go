package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestCheckSameAssetPosition() {
	addr := suite.AddAccounts(1, nil)
	msg := &types.MsgOpen{
		Creator:       addr[0].String(),
		Position:      types.Position_LONG,
		Leverage:      math.LegacyNewDec(1),
		TradingAsset:  ptypes.ATOM,
		Collateral:    sdk.NewCoin(ptypes.ATOM, math.NewInt(100)),
		StopLossPrice: math.LegacyZeroDec(),
	}

	mtp := types.NewMTP(suite.ctx, addr[0].String(), ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency, ptypes.ATOM, types.Position_LONG, types.TakeProfitPriceDefault, 1)
	mtp.StopLossPrice = math.LegacyZeroDec()
	testCases := []struct {
		name                 string
		exists               bool
		prerequisiteFunction func()
	}{
		{
			"mtp not found",
			false,
			func() {
			},
		},
		{
			"mtp not found because position is different",
			false,
			func() {
				msg.Position = types.Position_SHORT
				err := suite.app.PerpetualKeeper.SetMTP(suite.ctx, mtp)
				suite.Require().NoError(err)
			},
		},
		{
			"mtp not found because collateral is different",
			false,
			func() {
				msg.Position = types.Position_LONG
				err := suite.app.PerpetualKeeper.SetMTP(suite.ctx, mtp)
				suite.Require().NoError(err)
			},
		},
		{
			"mtp found",
			true,
			func() {
				msg.Collateral = sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			gotMtp := suite.app.PerpetualKeeper.CheckSameAssetPosition(suite.ctx, msg)
			if tc.exists {
				suite.Require().NotNil(gotMtp)
				suite.Require().Equal(mtp, gotMtp)
			} else {
				suite.Require().Nil(gotMtp)
			}
		})
	}
}
