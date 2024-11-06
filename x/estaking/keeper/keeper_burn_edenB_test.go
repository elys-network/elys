package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func (suite *EstakingKeeperTestSuite) TestBurnEdenBFromElysUnstaked(t *testing.T) {
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

func TestBurnEdenBFromEdenUncommitted(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	ek, commitmentKeeper := app.EstakingKeeper, app.CommitmentKeeper

	var committed sdk.Coins
	var unclaimed sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000))
	uedenBToken := sdk.NewCoin(ptypes.EdenB, sdk.NewInt(20000))
	unclaimed = unclaimed.Add(uedenToken, uedenBToken)

	// Mint coins
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, unclaimed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, unclaimed)
	require.NoError(t, err)

	// Prepare committed tokens
	uedenToken = sdk.NewCoin(ptypes.Eden, sdk.NewInt(10000))
	uedenBToken = sdk.NewCoin(ptypes.EdenB, sdk.NewInt(5000))
	committed = committed.Add(uedenToken, uedenBToken)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true, WithdrawEnabled: true})

	commitment := app.CommitmentKeeper.GetCommitments(ctx, genAccount)
	commitment.Claimed = commitment.Claimed.Add(committed...)
	app.CommitmentKeeper.SetCommitments(ctx, commitment)

	msgServer := commkeeper.NewMsgServerImpl(commitmentKeeper)
	_, err = msgServer.CommitClaimedRewards(ctx, &ctypes.MsgCommitClaimedRewards{
		Creator: genAccount.String(),
		Amount:  sdk.NewInt(1000),
		Denom:   ptypes.Eden,
	})
	require.NoError(t, err)

	// Track Elys staked amount
	ek.EndBlocker(ctx)

	// Uncommit tokens
	_, err = msgServer.UncommitTokens(sdk.WrapSDKContext(ctx), &ctypes.MsgUncommitTokens{
		Creator: genAccount.String(),
		Amount:  sdk.NewInt(1000),
		Denom:   ptypes.Eden,
	})
	require.NoError(t, err)
}
