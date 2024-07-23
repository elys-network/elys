package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (suite *KeeperTestSuite) TestOpen() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	testCases := []struct {
		name                 string
		input                *types.MsgOpen
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"failed user authorization",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  "stake",
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("10.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			errorsmod.Wrap(types.ErrUnauthorised, "unauthorised").Error(),
			func() {
				suite.EnableWhiteListing()
			},
		},
		{"no positions allowed",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  "stake",
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("10.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			errorsmod.Wrap(types.ErrMaxOpenPositions, "cannot open new positions").Error(),
			func() {
				suite.DisableWhiteListing()
				suite.SetMaxOpenPositions(0)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
