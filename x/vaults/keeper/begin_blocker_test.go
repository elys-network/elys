package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mastercheftypes "github.com/elys-network/elys/v7/x/masterchef/types"
	parametertypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/vaults/keeper"
	"github.com/elys-network/elys/v7/x/vaults/types"
)

func (suite *KeeperTestSuite) TestBeginBlocker_ManagementFee() {
	// Create test accounts
	depositor := sdk.AccAddress([]byte("depositor"))
	manager := sdk.AccAddress([]byte("manager"))
	protocolAddress := sdk.AccAddress([]byte("protocol"))

	// Set protocol address in masterchef params
	suite.app.MasterchefKeeper.SetParams(suite.ctx, mastercheftypes.Params{
		ProtocolRevenueAddress: protocolAddress.String(),
	})

	// Create the vault with management fee
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := types.MsgAddVault{
		Creator:          authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:     "uusdc",
		MaxAmountUsd:     sdkmath.LegacyNewDec(1000000),
		AllowedCoins:     []string{"uusdc", "uatom"},
		RewardCoins:      []string{"uelys"},
		BenchmarkCoin:    "uatom",
		Manager:          manager.String(),
		ManagementFee:    sdkmath.LegacyNewDecWithPrec(2, 2), // 2% management fee
		ProtocolFeeShare: sdkmath.LegacyNewDecWithPrec(5, 1), // 50% protocol fee share
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

	// Set total blocks per year
	suite.app.ParameterKeeper.SetParams(suite.ctx, parametertypes.Params{
		TotalBlocksPerYear: 1000,
	})

	// Run begin blocker
	suite.app.VaultsKeeper.BeginBlocker(suite.ctx)

	// Calculate expected fees
	// Management fee = deposit * (fee_rate / blocks_per_year)
	expectedFee := depositAmount.Amount.ToLegacyDec().Mul(addVault.ManagementFee).Quo(sdkmath.LegacyNewDec(1000))
	expectedManagerFee := expectedFee.Mul(sdkmath.LegacyNewDec(1).Sub(addVault.ProtocolFeeShare))
	expectedProtocolFee := expectedFee.Mul(addVault.ProtocolFeeShare)

	// Verify manager received their share of the fee
	managerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, manager, "uusdc")
	suite.Require().True(managerBalance.Amount.Equal(expectedManagerFee.TruncateInt()),
		"manager should receive correct fee amount")

	// Verify protocol address received their share of the fee
	protocolBalance := suite.app.BankKeeper.GetBalance(suite.ctx, protocolAddress, "uusdc")
	suite.Require().True(protocolBalance.Amount.Equal(expectedProtocolFee.TruncateInt()),
		"protocol address should receive correct fee amount")
}

func (suite *KeeperTestSuite) TestBeginBlocker_PerformanceFee() {
	// Create test accounts
	depositor := sdk.AccAddress([]byte("depositor"))
	manager := sdk.AccAddress([]byte("manager"))
	protocolAddress := sdk.AccAddress([]byte("protocol"))

	// Set protocol address in masterchef params
	suite.app.MasterchefKeeper.SetParams(suite.ctx, mastercheftypes.Params{
		ProtocolRevenueAddress: protocolAddress.String(),
	})

	// Create the vault with performance fee
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := types.MsgAddVault{
		Creator:          authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:     "uusdc",
		MaxAmountUsd:     sdkmath.LegacyNewDec(1000000),
		AllowedCoins:     []string{"uusdc", "uatom"},
		RewardCoins:      []string{"uelys"},
		BenchmarkCoin:    "uatom",
		Manager:          manager.String(),
		PerformanceFee:   sdkmath.LegacyNewDecWithPrec(2, 1), // 20% performance fee
		ProtocolFeeShare: sdkmath.LegacyNewDecWithPrec(5, 1), // 50% protocol fee share
		ManagementFee:    sdkmath.LegacyZeroDec(),            // 0% management fee
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

	// Set total blocks per year
	suite.app.ParameterKeeper.SetParams(suite.ctx, parametertypes.Params{
		TotalBlocksPerYear: 1000,
	})

	// Set performance fee epoch length
	suite.app.VaultsKeeper.SetParams(suite.ctx, types.Params{
		PerformanceFeeEpochLength: 1,
	})

	// Add some profit to the vault
	vaultAddress := types.NewVaultAddress(1)
	profitAmount := sdk.NewCoin("uusdc", sdkmath.NewInt(20000)) // 20% profit
	profitCoins := sdk.NewCoins(profitAmount)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", profitCoins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", vaultAddress, profitCoins)
	suite.Require().NoError(err)

	// query vault usd value
	usdValue, err := suite.app.VaultsKeeper.VaultUsdValue(suite.ctx, 1)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.LegacyMustNewDecFromStr("0.12"), usdValue.Dec())

	// Run begin blocker
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.VaultsKeeper.BeginBlocker(suite.ctx)

	// Verify manager received their share of the performance fee
	managerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, manager, "uusdc")
	suite.Require().Equal(managerBalance.Amount, sdkmath.NewInt(2),
		"manager should receive correct performance fee amount")

	// Verify protocol address received their share of the performance fee
	protocolBalance := suite.app.BankKeeper.GetBalance(suite.ctx, protocolAddress, "uusdc")
	suite.Require().Equal(protocolBalance.Amount, sdkmath.NewInt(1),
		"protocol address should receive correct performance fee amount")
}
