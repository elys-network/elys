package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/vaults/keeper"
	"github.com/elys-network/elys/v5/x/vaults/types"
)

func (suite *KeeperTestSuite) TestMsgServerWithdraw() {
	// Create test accounts
	depositor := sdk.AccAddress([]byte("depositor"))
	manager := sdk.AccAddress([]byte("manager"))

	// Create the vault first with the correct authority
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := types.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "uusdc",
		MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		AllowedCoins:  []string{"uusdc", "uatom"},
		RewardCoins:   []string{"uelys"},
		BenchmarkCoin: "uatom",
		Manager:       manager.String(),
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	// Setup initial deposit
	depositAmount := sdk.NewCoin("uusdc", sdkmath.NewInt(100000))
	coinsToSend := sdk.NewCoins(depositAmount)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coinsToSend)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", depositor, coinsToSend)
	suite.Require().NoError(err)

	// Make initial deposit
	depositMsg := types.MsgDeposit{
		VaultId:   1,
		Depositor: depositor.String(),
		Amount:    depositAmount,
	}
	depositResp, err := msgServer.Deposit(suite.ctx, &depositMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(depositResp)

	testCases := []struct {
		desc        string
		withdrawer  sdk.AccAddress
		vaultId     uint64
		shares      sdkmath.Int
		setup       func()
		expectError bool
		errMsg      string
	}{
		{
			desc:        "successful withdraw",
			withdrawer:  depositor,
			vaultId:     1,
			shares:      depositResp.Shares,
			setup:       func() {},
			expectError: false,
		},
		// {
		// 	desc:        "vault not found",
		// 	withdrawer:  depositor,
		// 	vaultId:     999,
		// 	shares:      sdkmath.NewInt(1000),
		// 	setup:       func() {},
		// 	expectError: true,
		// 	errMsg:      "vault not found",
		// },
		// {
		// 	desc:        "insufficient shares",
		// 	withdrawer:  depositor,
		// 	vaultId:     1,
		// 	shares:      sdkmath.NewInt(1000000), // More than deposited
		// 	setup:       func() {},
		// 	expectError: true,
		// 	errMsg:      "insufficient shares",
		// },
		// {
		// 	desc:        "invalid withdrawer",
		// 	withdrawer:  sdk.AccAddress([]byte("invalid")),
		// 	vaultId:     1,
		// 	shares:      sdkmath.NewInt(1000),
		// 	setup:       func() {},
		// 	expectError: true,
		// 	errMsg:      "invalid withdrawer",
		// },
	}

	for _, tc := range testCases {
		suite.Run(tc.desc, func() {
			// Setup test case if needed
			if tc.setup != nil {
				tc.setup()
			}

			// Prepare the withdraw message
			msg := types.MsgWithdraw{
				Withdrawer: tc.withdrawer.String(),
				VaultId:    tc.vaultId,
				Shares:     tc.shares,
			}

			// Execute withdraw
			_, err := msgServer.Withdraw(suite.ctx, &msg)

			if tc.expectError {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errMsg)
			} else {
				suite.Require().NoError(err)

				// Verify the withdraw was successful
				// 1. Check that shares were burned
				shareDenom := types.GetShareDenomForVault(tc.vaultId)
				shareBalance := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, tc.withdrawer)
				committedAmount := shareBalance.GetCommittedAmountForDenom(shareDenom)
				suite.Require().True(committedAmount.IsZero(), "shares should be burned")

				// 2. Check that tokens were returned to withdrawer
				withdrawerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.withdrawer, "uusdc")
				suite.Require().True(withdrawerBalance.Amount.GT(sdkmath.ZeroInt()), "withdrawer should receive tokens")

				// 3. Check that vault balance decreased
				vaultAddress := types.NewVaultAddress(tc.vaultId)
				vaultBalance := suite.app.BankKeeper.GetBalance(suite.ctx, vaultAddress, "uusdc")
				suite.Require().True(vaultBalance.Amount.LT(depositAmount.Amount), "vault balance should decrease")
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerWithdraw_WithRewards() {
	// Create test accounts
	depositor := sdk.AccAddress([]byte("depositor"))
	manager := sdk.AccAddress([]byte("manager"))

	// Create the vault first with the correct authority
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := types.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "uusdc",
		MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		AllowedCoins:  []string{"uusdc", "uatom"},
		RewardCoins:   []string{"uelys"},
		BenchmarkCoin: "uatom",
		Manager:       manager.String(),
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	// Setup initial deposit
	depositAmount := sdk.NewCoin("uusdc", sdkmath.NewInt(100000))
	coinsToSend := sdk.NewCoins(depositAmount)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coinsToSend)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", depositor, coinsToSend)
	suite.Require().NoError(err)

	// Make initial deposit
	depositMsg := types.MsgDeposit{
		VaultId:   1,
		Depositor: depositor.String(),
		Amount:    depositAmount,
	}
	depositResp, err := msgServer.Deposit(suite.ctx, &depositMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(depositResp)

	// Add rewards to vault
	vaultAddress := types.NewVaultAddress(1)
	rewardAmount := sdk.NewCoin("uelys", sdkmath.NewInt(1000))
	rewardCoins := sdk.NewCoins(rewardAmount)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", rewardCoins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", vaultAddress, rewardCoins)
	suite.Require().NoError(err)

	// Test withdraw with rewards
	msg := types.MsgWithdraw{
		Withdrawer: depositor.String(),
		VaultId:    1,
		Shares:     depositResp.Shares,
	}

	_, err = msgServer.Withdraw(suite.ctx, &msg)
	suite.Require().NoError(err)

	// Verify withdrawer received both deposit and rewards
	withdrawerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, depositor, "uusdc")
	suite.Require().True(withdrawerBalance.Amount.GT(sdkmath.ZeroInt()), "withdrawer should receive deposit tokens")

	rewardBalance := suite.app.BankKeeper.GetBalance(suite.ctx, depositor, "uelys")
	suite.Require().True(rewardBalance.Amount.GT(sdkmath.ZeroInt()), "withdrawer should receive reward tokens")
}
