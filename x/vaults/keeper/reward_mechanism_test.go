package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mastercheftypes "github.com/elys-network/elys/v6/x/masterchef/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/vaults/keeper"
	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (suite *KeeperTestSuite) TestRewardMechanism() {
	// Create test accounts
	depositor := sdk.AccAddress([]byte("depositor"))
	manager := sdk.AccAddress([]byte("manager"))
	protocolAddress := sdk.AccAddress([]byte("protocol"))

	// Set protocol address in masterchef params
	suite.app.MasterchefKeeper.SetParams(suite.ctx, mastercheftypes.Params{
		ProtocolRevenueAddress: protocolAddress.String(),
	})

	// Create the vault
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
	_, err = msgServer.Deposit(suite.ctx, &depositMsg)
	suite.Require().NoError(err)

	// Create a pool for the vault to join
	suite.CreateNewAmmPool(
		manager,
		false,                                    // useOracle
		sdkmath.LegacyMustNewDecFromStr("0.003"), // swapFee
		sdkmath.LegacyMustNewDecFromStr("0.003"), // exitFee
		"uatom",                                  // asset2
		sdkmath.NewInt(1000),                     // baseTokenAmount
		sdkmath.NewInt(1000),                     // assetAmount
	)

	// Join pool
	joinPoolMsg := types.MsgPerformActionJoinPool{
		Creator:        manager.String(),
		VaultId:        1,
		PoolId:         1,
		ShareAmountOut: sdkmath.NewInt(100),
		MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: sdkmath.NewInt(100)}},
	}
	_, err = msgServer.PerformActionJoinPool(suite.ctx, &joinPoolMsg)
	suite.Require().NoError(err)

	// Simulate rewards from masterchef
	rewardAmount := sdkmath.NewInt(1000)
	rewardDenom := ptypes.Eden
	suite.app.VaultsKeeper.UpdateAccPerShare(suite.ctx, 1, rewardDenom, rewardAmount)

	// Verify pool reward info was updated
	poolRewardInfo, found := suite.app.VaultsKeeper.GetPoolRewardInfo(suite.ctx, 1, rewardDenom)
	suite.Require().True(found)
	suite.Require().True(poolRewardInfo.PoolAccRewardPerShare.IsPositive())

	// Update user reward pending
	suite.app.VaultsKeeper.UpdateUserRewardPending(suite.ctx, 1, rewardDenom, depositor, false, sdkmath.ZeroInt())

	// Verify user reward info was updated
	userRewardInfo, found := suite.app.VaultsKeeper.GetUserRewardInfo(suite.ctx, depositor, 1, rewardDenom)
	suite.Require().True(found)
	suite.Require().True(userRewardInfo.RewardPending.IsPositive())

	// Update user reward debt
	suite.app.VaultsKeeper.UpdateUserRewardDebt(suite.ctx, 1, rewardDenom, depositor)

	// Verify user reward debt was updated
	userRewardInfo, found = suite.app.VaultsKeeper.GetUserRewardInfo(suite.ctx, depositor, 1, rewardDenom)
	suite.Require().True(found)
	suite.Require().True(userRewardInfo.RewardDebt.IsPositive())

	// Test claiming rewards
	claimMsg := types.MsgClaimRewards{
		Sender:   depositor.String(),
		VaultIds: []uint64{1},
	}
	_, err = msgServer.ClaimRewards(suite.ctx, &claimMsg)
	suite.Require().NoError(err)

	// Verify rewards were claimed
	userRewardInfo, found = suite.app.VaultsKeeper.GetUserRewardInfo(suite.ctx, depositor, 1, rewardDenom)
	suite.Require().True(found)
	suite.Require().True(userRewardInfo.RewardPending.IsZero())

	// Verify depositor received the rewards
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, depositor, rewardDenom)
	suite.Require().True(balance.Amount.IsPositive())
}

func (suite *KeeperTestSuite) TestRewardMechanismWithMultipleUsers() {
	// Create test accounts
	depositor1 := sdk.AccAddress([]byte("depositor1"))
	depositor2 := sdk.AccAddress([]byte("depositor2"))
	manager := sdk.AccAddress([]byte("manager"))

	// Create the vault
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

	// Setup initial deposits for both users
	depositAmount := sdk.NewCoin("uusdc", sdkmath.NewInt(100000))
	coinsToSend := sdk.NewCoins(depositAmount)

	// Fund and deposit for user 1
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coinsToSend)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", depositor1, coinsToSend)
	suite.Require().NoError(err)
	depositMsg1 := types.MsgDeposit{
		VaultId:   1,
		Depositor: depositor1.String(),
		Amount:    depositAmount,
	}
	_, err = msgServer.Deposit(suite.ctx, &depositMsg1)
	suite.Require().NoError(err)

	// Fund and deposit for user 2
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coinsToSend)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", depositor2, coinsToSend)
	suite.Require().NoError(err)
	depositMsg2 := types.MsgDeposit{
		VaultId:   1,
		Depositor: depositor2.String(),
		Amount:    depositAmount,
	}
	_, err = msgServer.Deposit(suite.ctx, &depositMsg2)
	suite.Require().NoError(err)

	// Create a pool and join it
	suite.CreateNewAmmPool(
		manager,
		false,
		sdkmath.LegacyMustNewDecFromStr("0.003"),
		sdkmath.LegacyMustNewDecFromStr("0.003"),
		"uatom",
		sdkmath.NewInt(1000),
		sdkmath.NewInt(1000),
	)

	joinPoolMsg := types.MsgPerformActionJoinPool{
		Creator:        manager.String(),
		VaultId:        1,
		PoolId:         1,
		ShareAmountOut: sdkmath.NewInt(100),
		MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: sdkmath.NewInt(100)}},
	}
	_, err = msgServer.PerformActionJoinPool(suite.ctx, &joinPoolMsg)
	suite.Require().NoError(err)

	// Simulate rewards from masterchef
	rewardAmount := sdkmath.NewInt(1000)
	rewardDenom := ptypes.Eden
	suite.app.VaultsKeeper.UpdateAccPerShare(suite.ctx, 1, rewardDenom, rewardAmount)

	// Update rewards for both users
	suite.app.VaultsKeeper.UpdateUserRewardPending(suite.ctx, 1, rewardDenom, depositor1, false, sdkmath.ZeroInt())
	suite.app.VaultsKeeper.UpdateUserRewardPending(suite.ctx, 1, rewardDenom, depositor2, false, sdkmath.ZeroInt())

	// Verify both users have pending rewards
	user1RewardInfo, found := suite.app.VaultsKeeper.GetUserRewardInfo(suite.ctx, depositor1, 1, rewardDenom)
	suite.Require().True(found)
	suite.Require().True(user1RewardInfo.RewardPending.IsPositive())

	user2RewardInfo, found := suite.app.VaultsKeeper.GetUserRewardInfo(suite.ctx, depositor2, 1, rewardDenom)
	suite.Require().True(found)
	suite.Require().True(user2RewardInfo.RewardPending.IsPositive())

	// Test claiming rewards for both users
	claimMsg1 := types.MsgClaimRewards{
		Sender:   depositor1.String(),
		VaultIds: []uint64{1},
	}
	_, err = msgServer.ClaimRewards(suite.ctx, &claimMsg1)
	suite.Require().NoError(err)

	claimMsg2 := types.MsgClaimRewards{
		Sender:   depositor2.String(),
		VaultIds: []uint64{1},
	}
	_, err = msgServer.ClaimRewards(suite.ctx, &claimMsg2)
	suite.Require().NoError(err)

	// Verify both users received their rewards
	balance1 := suite.app.BankKeeper.GetBalance(suite.ctx, depositor1, rewardDenom)
	balance2 := suite.app.BankKeeper.GetBalance(suite.ctx, depositor2, rewardDenom)
	suite.Require().True(balance1.Amount.IsPositive())
	suite.Require().True(balance2.Amount.IsPositive())
}
