package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *EstakingKeeperTestSuite) TestBurnEdenBFromElysUnstaked() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"burn EdenB from Elys unstaked",
			func() {
				var committed sdk.Coins
				var unclaimed sdk.Coins

				// Prepare unclaimed tokens
				uedenToken := sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000))
				uedenBToken := sdk.NewCoin(ptypes.EdenB, sdk.NewInt(20000))
				unclaimed = unclaimed.Add(uedenToken, uedenBToken)

				// Mint coins
				err := suite.app.BankKeeper.MintCoins(suite.ctx, ctypes.ModuleName, unclaimed)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ctypes.ModuleName, suite.genAccount, unclaimed)
				suite.Require().NoError(err)

				// Prepare committed tokens
				uedenToken = sdk.NewCoin(ptypes.Eden, sdk.NewInt(10000))
				uedenBToken = sdk.NewCoin(ptypes.EdenB, sdk.NewInt(5000))
				committed = committed.Add(uedenToken, uedenBToken)

				// Mint coins
				err = suite.app.BankKeeper.MintCoins(suite.ctx, ctypes.ModuleName, committed)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ctypes.ModuleName, suite.genAccount, committed)
				suite.Require().NoError(err)

				// Add testing commitment
				suite.AddTestCommitment(committed)

				// Take elys staked snapshot
				suite.app.EstakingKeeper.TakeDelegationSnapshot(suite.ctx, suite.genAccount)
			},
			func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()

			// burn amount = 100000 (unbonded amt) / (1000000 (elys staked) + 10000 (Eden committed)) * (20000 EdenB + 5000 EdenB committed)
			unbondAmt, err := suite.app.StakingKeeper.Unbond(suite.ctx, suite.genAccount, suite.valAddr, sdk.NewDecWithPrec(10, 2))
			suite.Require().Equal(unbondAmt, sdk.NewInt(100000))
			suite.Require().NoError(err)

			// Process EdenB burn operation
			suite.app.EstakingKeeper.EndBlocker(suite.ctx)

			tc.postValidateFunction()
		})
	}
}

func (suite *EstakingKeeperTestSuite) TestBurnEdenBFromEdenUncommitted() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"burn EdenB from Eden uncommitted",
			func() {
				var committed sdk.Coins
				var unclaimed sdk.Coins

				// Prepare unclaimed tokens
				uedenToken := sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000))
				uedenBToken := sdk.NewCoin(ptypes.EdenB, sdk.NewInt(20000))
				unclaimed = unclaimed.Add(uedenToken, uedenBToken)

				// Mint coins
				err := suite.app.BankKeeper.MintCoins(suite.ctx, ctypes.ModuleName, unclaimed)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ctypes.ModuleName, suite.genAccount, unclaimed)
				suite.Require().NoError(err)

				// Prepare committed tokens
				uedenToken = sdk.NewCoin(ptypes.Eden, sdk.NewInt(10000))
				uedenBToken = sdk.NewCoin(ptypes.EdenB, sdk.NewInt(5000))
				committed = committed.Add(uedenToken, uedenBToken)

				// Set assetprofile entry for denom
				suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true, WithdrawEnabled: true})

				commitment := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, suite.genAccount)
				commitment.Claimed = commitment.Claimed.Add(committed...)
				suite.app.CommitmentKeeper.SetCommitments(suite.ctx, commitment)
			},
			func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			// Track Elys staked amount
			suite.app.EstakingKeeper.EndBlocker(suite.ctx)

			msgServer := commkeeper.NewMsgServerImpl(suite.app.CommitmentKeeper)

			_, err := msgServer.CommitClaimedRewards(suite.ctx, &ctypes.MsgCommitClaimedRewards{
				Creator: suite.genAccount.String(),
				Amount:  sdk.NewInt(1000),
				Denom:   ptypes.Eden,
			})
			suite.Require().NoError(err)

			// Uncommit tokens
			_, err = msgServer.UncommitTokens(sdk.WrapSDKContext(suite.ctx), &ctypes.MsgUncommitTokens{
				Creator: suite.genAccount.String(),
				Amount:  sdk.NewInt(1000),
				Denom:   ptypes.Eden,
			})
			suite.Require().NoError(err)

			tc.postValidateFunction()
		})
	}
}
