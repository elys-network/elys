package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestCheckSameAssetPosition() {
	addr := suite.AddAccounts(1)
	msg := &types.MsgOpen{
		Creator:       addr[0].String(),
		Position:      types.Position_LONG,
		Leverage:      sdk.NewDec(1),
		TradingAsset:  ptypes.ATOM,
		Collateral:    sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100)),
		StopLossPrice: sdk.NewDec(100),
	}

	mtp := types.NewMTP(addr[0].String(), ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency, ptypes.ATOM, types.Position_LONG, sdk.NewDec(5), sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), 1)

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
			"mtp found",
			true,
			func() {
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
