package keeper_test

import (
	"testing"
)

func (suite *KeeperTestSuite) TestBeginBlocker_ManagementFee(t *testing.T) {
	// Create test accounts
	// depositor := sdk.AccAddress([]byte("depositor"))
	// manager := sdk.AccAddress([]byte("manager"))
	// protocolAddress := sdk.AccAddress([]byte("protocol"))

	// Set protocol address in masterchef params
	// suite.app.MasterchefKeeper.SetParams(suite.ctx, types.Params{
	// 	ProtocolRevenueAddress: protocolAddress.String(),
	// })

	// // Create the vault with management fee
	// msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	// addVault := types.MsgAddVault{
	// 	Creator:          authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	// 	DepositDenom:     "uusdc",
	// 	MaxAmountUsd:     sdkmath.LegacyNewDec(1000000),
	// 	AllowedCoins:     []string{"uusdc", "uatom"},
	// 	RewardCoins:      []string{"uelys"},
	// 	BenchmarkCoin:    "uatom",
	// 	Manager:          manager.String(),
	// 	ManagementFee:    sdkmath.LegacyNewDecWithPrec(2, 2), // 2% management fee
	// 	ProtocolFeeShare: sdkmath.LegacyNewDecWithPrec(5, 1), // 50% protocol fee share
	// }
	// _, err := msgServer.AddVault(suite.ctx, &addVault)
	// suite.Require().NoError(err)

	// // Setup initial deposit
	// depositAmount := sdk.NewCoin("uusdc", sdkmath.NewInt(100000))
	// coinsToSend := sdk.NewCoins(depositAmount)
	// err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coinsToSend)
	// suite.Require().NoError(err)
	// err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", depositor, coinsToSend)
	// suite.Require().NoError(err)

	// // Make initial deposit
	// depositMsg := types.MsgDeposit{
	// 	VaultId:   1,
	// 	Depositor: depositor.String(),
	// 	Amount:    depositAmount,
	// }
	// _, err = msgServer.Deposit(suite.ctx, &depositMsg)
	// suite.Require().NoError(err)

	// // Set total blocks per year
	// suite.app.ParamsKeeper.SetParams(suite.ctx, types.Params{
	// 	TotalBlocksPerYear: 1000,
	// })

	// // Run begin blocker
	// suite.app.VaultsKeeper.BeginBlocker(suite.ctx)

	// // Calculate expected fees
	// // Management fee = deposit * (fee_rate / blocks_per_year)
	// expectedFee := depositAmount.Amount.ToLegacyDec().Mul(addVault.ManagementFee).Quo(sdkmath.LegacyNewDec(1000))
	// expectedManagerFee := expectedFee.Mul(sdkmath.LegacyNewDec(1).Sub(addVault.ProtocolFeeShare))
	// expectedProtocolFee := expectedFee.Mul(addVault.ProtocolFeeShare)

	// // Verify manager received their share of the fee
	// managerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, manager, "uusdc")
	// require.True(t, managerBalance.Amount.Equal(expectedManagerFee.TruncateInt()),
	// 	"manager should receive correct fee amount")

	// // Verify protocol address received their share of the fee
	// protocolBalance := suite.app.BankKeeper.GetBalance(suite.ctx, protocolAddress, "uusdc")
	// require.True(t, protocolBalance.Amount.Equal(expectedProtocolFee.TruncateInt()),
	// 	"protocol address should receive correct fee amount")
}

// func (suite *KeeperTestSuite) TestBeginBlocker_PerformanceFee(t *testing.T) {
// 	// Create test accounts
// 	depositor := sdk.AccAddress([]byte("depositor"))
// 	manager := sdk.AccAddress([]byte("manager"))
// 	protocolAddress := sdk.AccAddress([]byte("protocol"))

// 	// Set protocol address in masterchef params
// 	suite.app.MasterchefKeeper.SetParams(suite.ctx, types.Params{
// 		ProtocolRevenueAddress: protocolAddress.String(),
// 	})

// 	// Create the vault with performance fee
// 	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
// 	addVault := types.MsgAddVault{
// 		Creator:          authtypes.NewModuleAddress(govtypes.ModuleName).String(),
// 		DepositDenom:     "uusdc",
// 		MaxAmountUsd:     sdkmath.LegacyNewDec(1000000),
// 		AllowedCoins:     []string{"uusdc", "uatom"},
// 		RewardCoins:      []string{"uelys"},
// 		BenchmarkCoin:    "uatom",
// 		Manager:          manager.String(),
// 		PerformanceFee:   sdkmath.LegacyNewDecWithPrec(2, 1), // 20% performance fee
// 		ProtocolFeeShare: sdkmath.LegacyNewDecWithPrec(5, 1), // 50% protocol fee share
// 	}
// 	_, err := msgServer.AddVault(suite.ctx, &addVault)
// 	suite.Require().NoError(err)

// 	// Setup initial deposit
// 	depositAmount := sdk.NewCoin("uusdc", sdkmath.NewInt(100000))
// 	coinsToSend := sdk.NewCoins(depositAmount)
// 	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coinsToSend)
// 	suite.Require().NoError(err)
// 	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", depositor, coinsToSend)
// 	suite.Require().NoError(err)

// 	// Make initial deposit
// 	depositMsg := types.MsgDeposit{
// 		VaultId:   1,
// 		Depositor: depositor.String(),
// 		Amount:    depositAmount,
// 	}
// 	_, err = msgServer.Deposit(suite.ctx, &depositMsg)
// 	suite.Require().NoError(err)

// 	// Set total blocks per year
// 	suite.app.ParamsKeeper.SetParams(suite.ctx, types.Params{
// 		TotalBlocksPerYear: 1000,
// 	})

// 	// Set performance fee epoch length
// 	suite.app.VaultsKeeper.SetParams(suite.ctx, types.Params{
// 		PerformanceFeeEpochLength: 1,
// 	})

// 	// Add some profit to the vault
// 	vaultAddress := types.NewVaultAddress(1)
// 	profitAmount := sdk.NewCoin("uusdc", sdkmath.NewInt(20000)) // 20% profit
// 	profitCoins := sdk.NewCoins(profitAmount)
// 	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", profitCoins)
// 	suite.Require().NoError(err)
// 	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", vaultAddress, profitCoins)
// 	suite.Require().NoError(err)

// 	// Run begin blocker
// 	suite.app.VaultsKeeper.BeginBlocker(suite.ctx)

// 	// Calculate expected performance fee
// 	// Performance fee = profit * performance_fee_rate
// 	expectedFee := profitAmount.Amount.ToLegacyDec().Mul(addVault.PerformanceFee)
// 	expectedManagerFee := expectedFee.Mul(sdkmath.LegacyNewDec(1).Sub(addVault.ProtocolFeeShare))
// 	expectedProtocolFee := expectedFee.Mul(addVault.ProtocolFeeShare)

// 	// Verify manager received their share of the performance fee
// 	managerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, manager, "uusdc")
// 	require.True(t, managerBalance.Amount.Equal(expectedManagerFee.TruncateInt()),
// 		"manager should receive correct performance fee amount")

// 	// Verify protocol address received their share of the performance fee
// 	protocolBalance := suite.app.BankKeeper.GetBalance(suite.ctx, protocolAddress, "uusdc")
// 	require.True(t, protocolBalance.Amount.Equal(expectedProtocolFee.TruncateInt()),
// 		"protocol address should receive correct performance fee amount")
// }
