package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *EstakingKeeperTestSuite) TestHooksStaking() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"AfterUnbondingInitiated",
			func() {
				suite.ResetSuite()

				hooks := suite.app.EstakingKeeper.Hooks()
				err := hooks.AfterUnbondingInitiated(suite.ctx, 0)
				suite.Require().Nil(err)
			},
			func() {},
		},
		{
			"BeforeValidatorModified",
			func() {
				suite.ResetSuite()

				hooks := suite.app.EstakingKeeper.Hooks()
				err := hooks.BeforeValidatorModified(suite.ctx, suite.valAddr)
				suite.Require().Nil(err)
			},
			func() {},
		},
		{
			"AfterValidatorRemoved",
			func() {
				suite.ResetSuite()

				hooks := suite.app.EstakingKeeper.Hooks()
				consAddr := sdk.ConsAddress(suite.ctx.BlockHeader().ProposerAddress)
				err := hooks.AfterValidatorRemoved(suite.ctx, consAddr, suite.valAddr)
				suite.Require().Nil(err)
			},
			func() {},
		},
		{
			"AfterValidatorBonded",
			func() {
				suite.ResetSuite()

				hooks := suite.app.EstakingKeeper.Hooks()
				consAddr := sdk.ConsAddress(suite.ctx.BlockHeader().ProposerAddress)
				err := hooks.AfterValidatorBonded(suite.ctx, consAddr, suite.valAddr)
				suite.Require().Nil(err)
			},
			func() {},
		},
		{
			"AfterValidatorBeginUnbonding",
			func() {
				suite.ResetSuite()

				hooks := suite.app.EstakingKeeper.Hooks()
				consAddr := sdk.ConsAddress(suite.ctx.BlockHeader().ProposerAddress)
				err := hooks.AfterValidatorBeginUnbonding(suite.ctx, consAddr, suite.valAddr)
				suite.Require().Nil(err)
			},
			func() {},
		},
		{
			"BeforeValidatorSlashed",
			func() {
				suite.ResetSuite()

				hooks := suite.app.EstakingKeeper.Hooks()
				err := hooks.BeforeValidatorSlashed(suite.ctx, suite.valAddr, math.LegacyOneDec())
				suite.Require().Nil(err)
			},
			func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
