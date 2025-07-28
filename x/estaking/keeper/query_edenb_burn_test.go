package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/estaking/types"
)

func (suite *EstakingKeeperTestSuite) TestEdenBBurnAmount() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func() (addr sdk.AccAddress)
		postValidateFunction func(addr sdk.AccAddress)
	}{
		{
			"burn EdenB from Elys unstaked",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()
				addr = suite.genAccount

				// Set asset profile entries
				suite.SetAssetProfile()

				// Prepare unclaimed and committed tokens
				//_ = suite.PrepareUnclaimedTokens()
				committed := suite.PrepareCommittedTokens()
				suite.AddTestCommitment(committed)

				// Set Elys staked
				suite.app.EstakingKeeper.SetElysStaked(suite.ctx, types.ElysStaked{
					Address: addr.String(),
					Amount:  math.NewInt(1000000),
				})

				return addr
			},
			func(addr sdk.AccAddress) {
				// Test burning EdenB from Elys unstaked
				res, err := suite.app.EstakingKeeper.EdenBBurnAmount(suite.ctx, &types.QueryEdenBBurnAmountRequest{
					Address:   addr.String(),
					TokenType: types.TokenType_TOKEN_TYPE_ELYS,
					Amount:    math.NewInt(100000),
				})
				suite.Require().NoError(err)
				// Expected burn amount = 100000 (unbonded amt) / (1000000 (elys staked) + 10000 (Eden committed)) * (5000 EdenB + 5000 EdenB committed)
				suite.Require().Equal(res.BurnEdenbAmount.String(), "990")
			},
		},
		{
			"burn EdenB from Eden uncommitted",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()
				addr = suite.genAccount

				// Set asset profile entries
				suite.SetAssetProfile()

				// Prepare unclaimed and committed tokens
				committed := suite.PrepareCommittedTokens()
				suite.AddTestCommitment(committed)

				// Set Elys staked
				suite.app.EstakingKeeper.SetElysStaked(suite.ctx, types.ElysStaked{
					Address: addr.String(),
					Amount:  math.NewInt(1000000),
				})

				return addr
			},
			func(addr sdk.AccAddress) {
				// Test burning EdenB from Eden uncommitted
				res, err := suite.app.EstakingKeeper.EdenBBurnAmount(suite.ctx, &types.QueryEdenBBurnAmountRequest{
					Address:   addr.String(),
					TokenType: types.TokenType_TOKEN_TYPE_EDEN,
					Amount:    math.NewInt(5000),
				})
				suite.Require().NoError(err)
				// Expected burn amount = 5000 (uncommitted amt) / (1000000 (elys staked) + 10000 (Eden committed) + 5000 (uncommitted amt)) * (5000 EdenB + 5000 EdenB committed)
				suite.Require().Equal(res.BurnEdenbAmount.String(), "49")
			},
		},
		{
			"invalid request - nil request",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()
				addr = suite.AddAccounts(1, nil)[0]
				return addr
			},
			func(addr sdk.AccAddress) {
				_, err := suite.app.EstakingKeeper.EdenBBurnAmount(suite.ctx, nil)
				suite.Require().Error(err)
			},
		},
		{
			"invalid request - invalid address",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()
				addr = suite.AddAccounts(1, nil)[0]
				return addr
			},
			func(addr sdk.AccAddress) {
				_, err := suite.app.EstakingKeeper.EdenBBurnAmount(suite.ctx, &types.QueryEdenBBurnAmountRequest{
					Address:   "invalid_address",
					TokenType: types.TokenType_TOKEN_TYPE_ELYS,
					Amount:    math.NewInt(100000),
				})
				suite.Require().Error(err)
			},
		},
		{
			"no staked Elys",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()
				addr = suite.AddAccounts(1, nil)[0]
				return addr
			},
			func(addr sdk.AccAddress) {
				_, err := suite.app.EstakingKeeper.EdenBBurnAmount(suite.ctx, &types.QueryEdenBBurnAmountRequest{
					Address:   addr.String(),
					TokenType: types.TokenType_TOKEN_TYPE_ELYS,
					Amount:    math.NewInt(100000),
				})
				suite.Require().Error(err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			addr := tc.prerequisiteFunction()
			tc.postValidateFunction(addr)
		})
	}
}
