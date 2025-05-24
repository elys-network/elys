package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestCheckUserAuthorization() {
	msg := &types.MsgOpen{Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5"}

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{
			"user not whitelisted",
			types.ErrUnauthorised.Error(),
			func() {
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.WhitelistingEnabled = true
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
		},
		{
			"success",
			"",
			func() {
				suite.app.PerpetualKeeper.WhitelistAddress(suite.ctx, sdk.MustAccAddressFromBech32(msg.Creator))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			err := suite.app.PerpetualKeeper.CheckUserAuthorization(suite.ctx, msg)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
