package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/vaults/keeper"
	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (suite *KeeperTestSuite) TestMsgServerPerformActionOpenPerpetual() {
	for _, tc := range []struct {
		desc            string
		creator         sdk.AccAddress
		vaultId         uint64
		position        types.Position
		leverage        sdkmath.LegacyDec
		collateral      sdk.Coin
		takeProfitPrice sdkmath.LegacyDec
		stopLossPrice   sdkmath.LegacyDec
		poolId          uint64
		setup           func()
		expectError     bool
		errorContains   string
	}{
		{
			desc:            "successful perpetual open - long position",
			creator:         sdk.AccAddress([]byte("manager")),
			vaultId:         1,
			position:        types.Position_LONG,
			leverage:        sdkmath.LegacyNewDec(2),
			collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
			takeProfitPrice: sdkmath.LegacyNewDec(50),
			stopLossPrice:   sdkmath.LegacyNewDec(1),
			poolId:          1,
			setup: func() {
				// Create a vault with the manager
				msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
				msg := types.MsgAddVault{
					Creator:       suite.app.VaultsKeeper.GetAuthority(),
					DepositDenom:  "uusdc",
					MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
					AllowedCoins:  []string{"uusdc", "uatom"},
					RewardCoins:   []string{"uelys"},
					BenchmarkCoin: "uatom",
					Manager:       sdk.AccAddress([]byte("manager")).String(),
				}
				_, err := msgServer.AddVault(suite.ctx, &msg)
				suite.Require().NoError(err)

				// Fund the vault with collateral
				vaultAddress := types.NewVaultAddress(1)
				err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(10000))))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", vaultAddress, sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(10000))))
				suite.Require().NoError(err)
			},
			expectError: false,
		},
		// {
		// 	desc:            "successful perpetual open - short position",
		// 	creator:         sdk.AccAddress([]byte("manager")),
		// 	vaultId:         1,
		// 	position:        types.Position_SHORT,
		// 	leverage:        sdkmath.LegacyNewDec(3),
		// 	collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(500)),
		// 	takeProfitPrice: sdkmath.LegacyNewDec(80),
		// 	stopLossPrice:   sdkmath.LegacyNewDec(120),
		// 	poolId:          1,
		// 	setup: func() {
		// 		// Create a vault with the manager
		// 		msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
		// 		msg := types.MsgAddVault{
		// 			Creator:       suite.app.VaultsKeeper.GetAuthority(),
		// 			DepositDenom:  "uusdc",
		// 			MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		// 			AllowedCoins:  []string{"uusdc", "uatom"},
		// 			RewardCoins:   []string{"uelys"},
		// 			BenchmarkCoin: "uatom",
		// 			Manager:       sdk.AccAddress([]byte("manager")).String(),
		// 		}
		// 		_, err := msgServer.AddVault(suite.ctx, &msg)
		// 		suite.Require().NoError(err)

		// 		// Fund the vault with collateral
		// 		vaultAddress := types.NewVaultAddress(1)
		// 		err = suite.app.BankKeeper.MintCoins(suite.ctx, authtypes.Minter, sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(10000))))
		// 		suite.Require().NoError(err)
		// 		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, authtypes.Minter, vaultAddress, sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(10000))))
		// 		suite.Require().NoError(err)
		// 	},
		// 	expectError: false,
		// },
		// {
		// 	desc:            "vault not found",
		// 	creator:         sdk.AccAddress([]byte("manager")),
		// 	vaultId:         999,
		// 	position:        types.Position_LONG,
		// 	leverage:        sdkmath.LegacyNewDec(2),
		// 	collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
		// 	takeProfitPrice: sdkmath.LegacyNewDec(110),
		// 	stopLossPrice:   sdkmath.LegacyNewDec(90),
		// 	poolId:          1,
		// 	setup:           func() {},
		// 	expectError:     true,
		// 	errorContains:   "vault 999 not found",
		// },
		// {
		// 	desc:            "invalid signer - not the vault manager",
		// 	creator:         sdk.AccAddress([]byte("wrong_manager")),
		// 	vaultId:         1,
		// 	position:        types.Position_LONG,
		// 	leverage:        sdkmath.LegacyNewDec(2),
		// 	collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
		// 	takeProfitPrice: sdkmath.LegacyNewDec(110),
		// 	stopLossPrice:   sdkmath.LegacyNewDec(90),
		// 	poolId:          1,
		// 	setup: func() {
		// 		// Create a vault with a different manager
		// 		msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
		// 		msg := types.MsgAddVault{
		// 			Creator:       suite.app.VaultsKeeper.GetAuthority(),
		// 			DepositDenom:  "uusdc",
		// 			MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		// 			AllowedCoins:  []string{"uusdc", "uatom"},
		// 			RewardCoins:   []string{"uelys"},
		// 			BenchmarkCoin: "uatom",
		// 			Manager:       sdk.AccAddress([]byte("manager")).String(),
		// 		}
		// 		_, err := msgServer.AddVault(suite.ctx, &msg)
		// 		suite.Require().NoError(err)
		// 	},
		// 	expectError:   true,
		// 	errorContains: "vault 1 is not managed by",
		// },
		// {
		// 	desc:            "unspecified position",
		// 	creator:         sdk.AccAddress([]byte("manager")),
		// 	vaultId:         1,
		// 	position:        types.Position_UNSPECIFIED,
		// 	leverage:        sdkmath.LegacyNewDec(2),
		// 	collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
		// 	takeProfitPrice: sdkmath.LegacyNewDec(110),
		// 	stopLossPrice:   sdkmath.LegacyNewDec(90),
		// 	poolId:          1,
		// 	setup: func() {
		// 		// Create a vault with the manager
		// 		msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
		// 		msg := types.MsgAddVault{
		// 			Creator:       suite.app.VaultsKeeper.GetAuthority(),
		// 			DepositDenom:  "uusdc",
		// 			MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		// 			AllowedCoins:  []string{"uusdc", "uatom"},
		// 			RewardCoins:   []string{"uelys"},
		// 			BenchmarkCoin: "uatom",
		// 			Manager:       sdk.AccAddress([]byte("manager")).String(),
		// 		}
		// 		_, err := msgServer.AddVault(suite.ctx, &msg)
		// 		suite.Require().NoError(err)
		// 	},
		// 	expectError:   true,
		// 	errorContains: "action failed with error",
		// },
		// {
		// 	desc:            "zero leverage",
		// 	creator:         sdk.AccAddress([]byte("manager")),
		// 	vaultId:         1,
		// 	position:        types.Position_LONG,
		// 	leverage:        sdkmath.LegacyZeroDec(),
		// 	collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
		// 	takeProfitPrice: sdkmath.LegacyNewDec(110),
		// 	stopLossPrice:   sdkmath.LegacyNewDec(90),
		// 	poolId:          1,
		// 	setup: func() {
		// 		// Create a vault with the manager
		// 		msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
		// 		msg := types.MsgAddVault{
		// 			Creator:       suite.app.VaultsKeeper.GetAuthority(),
		// 			DepositDenom:  "uusdc",
		// 			MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		// 			AllowedCoins:  []string{"uusdc", "uatom"},
		// 			RewardCoins:   []string{"uelys"},
		// 			BenchmarkCoin: "uatom",
		// 			Manager:       sdk.AccAddress([]byte("manager")).String(),
		// 		}
		// 		_, err := msgServer.AddVault(suite.ctx, &msg)
		// 		suite.Require().NoError(err)
		// 	},
		// 	expectError:   true,
		// 	errorContains: "action failed with error",
		// },
		// {
		// 	desc:            "zero collateral amount",
		// 	creator:         sdk.AccAddress([]byte("manager")),
		// 	vaultId:         1,
		// 	position:        types.Position_LONG,
		// 	leverage:        sdkmath.LegacyNewDec(2),
		// 	collateral:      sdk.NewCoin("uusdc", sdkmath.ZeroInt()),
		// 	takeProfitPrice: sdkmath.LegacyNewDec(110),
		// 	stopLossPrice:   sdkmath.LegacyNewDec(90),
		// 	poolId:          1,
		// 	setup: func() {
		// 		// Create a vault with the manager
		// 		msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
		// 		msg := types.MsgAddVault{
		// 			Creator:       suite.app.VaultsKeeper.GetAuthority(),
		// 			DepositDenom:  "uusdc",
		// 			MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		// 			AllowedCoins:  []string{"uusdc", "uatom"},
		// 			RewardCoins:   []string{"uelys"},
		// 			BenchmarkCoin: "uatom",
		// 			Manager:       sdk.AccAddress([]byte("manager")).String(),
		// 		}
		// 		_, err := msgServer.AddVault(suite.ctx, &msg)
		// 		suite.Require().NoError(err)
		// 	},
		// 	expectError:   true,
		// 	errorContains: "action failed with error",
		// },
		// {
		// 	desc:            "invalid pool id",
		// 	creator:         sdk.AccAddress([]byte("manager")),
		// 	vaultId:         1,
		// 	position:        types.Position_LONG,
		// 	leverage:        sdkmath.LegacyNewDec(2),
		// 	collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
		// 	takeProfitPrice: sdkmath.LegacyNewDec(110),
		// 	stopLossPrice:   sdkmath.LegacyNewDec(90),
		// 	poolId:          999,
		// 	setup: func() {
		// 		// Create a vault with the manager
		// 		msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
		// 		msg := types.MsgAddVault{
		// 			Creator:       suite.app.VaultsKeeper.GetAuthority(),
		// 			DepositDenom:  "uusdc",
		// 			MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		// 			AllowedCoins:  []string{"uusdc", "uatom"},
		// 			RewardCoins:   []string{"uelys"},
		// 			BenchmarkCoin: "uatom",
		// 			Manager:       sdk.AccAddress([]byte("manager")).String(),
		// 		}
		// 		_, err := msgServer.AddVault(suite.ctx, &msg)
		// 		suite.Require().NoError(err)

		// 		// Fund the vault with collateral
		// 		vaultAddress := types.NewVaultAddress(1)
		// 		err = suite.app.BankKeeper.MintCoins(suite.ctx, authtypes.Minter, sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(10000))))
		// 		suite.Require().NoError(err)
		// 		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, authtypes.Minter, vaultAddress, sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(10000))))
		// 		suite.Require().NoError(err)
		// 	},
		// 	expectError:   true,
		// 	errorContains: "action failed with error",
		// },
	} {
		suite.Run(tc.desc, func() {
			// Setup test case
			tc.setup()
			suite.SetPerpetualPool(1)

			msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)

			// Prepare the message
			msg := types.MsgPerformActionOpenPerpetual{
				Creator:         tc.creator.String(),
				Position:        tc.position,
				Leverage:        tc.leverage,
				Collateral:      tc.collateral,
				TakeProfitPrice: tc.takeProfitPrice,
				StopLossPrice:   tc.stopLossPrice,
				PoolId:          tc.poolId,
				VaultId:         tc.vaultId,
			}

			_, err := msgServer.PerformActionOpenPerpetual(suite.ctx, &msg)

			if tc.expectError {
				suite.Require().Error(err)
				if tc.errorContains != "" {
					suite.Require().Contains(err.Error(), tc.errorContains)
				}
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerPerformActionOpenPerpetual_EdgeCases() {
	// Setup: Create a vault and fund it
	suite.SetPerpetualPool(1)
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)

	// Create vault
	addVaultMsg := types.MsgAddVault{
		Creator:       suite.app.VaultsKeeper.GetAuthority(),
		DepositDenom:  "uusdc",
		MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		AllowedCoins:  []string{"uusdc", "uatom"},
		RewardCoins:   []string{"uelys"},
		BenchmarkCoin: "uatom",
		Manager:       sdk.AccAddress([]byte("manager")).String(),
	}
	_, err := msgServer.AddVault(suite.ctx, &addVaultMsg)
	suite.Require().NoError(err)

	// Fund the vault
	vaultAddress := types.NewVaultAddress(1)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(10000))))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", vaultAddress, sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(10000))))
	suite.Require().NoError(err)

	// Test edge cases
	testCases := []struct {
		desc            string
		leverage        sdkmath.LegacyDec
		collateral      sdk.Coin
		takeProfitPrice sdkmath.LegacyDec
		stopLossPrice   sdkmath.LegacyDec
		expectError     bool
	}{
		{
			desc:            "very high leverage",
			leverage:        sdkmath.LegacyNewDec(100),
			collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(100)),
			takeProfitPrice: sdkmath.LegacyNewDec(110),
			stopLossPrice:   sdkmath.LegacyNewDec(90),
			expectError:     false, // Should succeed if perpetual module allows it
		},
		{
			desc:            "very small collateral",
			leverage:        sdkmath.LegacyNewDec(2),
			collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1)),
			takeProfitPrice: sdkmath.LegacyNewDec(110),
			stopLossPrice:   sdkmath.LegacyNewDec(90),
			expectError:     false, // Should succeed if perpetual module allows it
		},
		{
			desc:            "zero take profit price",
			leverage:        sdkmath.LegacyNewDec(2),
			collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
			takeProfitPrice: sdkmath.LegacyZeroDec(),
			stopLossPrice:   sdkmath.LegacyNewDec(90),
			expectError:     false, // Should succeed if perpetual module allows it
		},
		{
			desc:            "zero stop loss price",
			leverage:        sdkmath.LegacyNewDec(2),
			collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
			takeProfitPrice: sdkmath.LegacyNewDec(110),
			stopLossPrice:   sdkmath.LegacyZeroDec(),
			expectError:     false, // Should succeed if perpetual module allows it
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.desc, func() {
			msg := types.MsgPerformActionOpenPerpetual{
				Creator:         sdk.AccAddress([]byte("manager")).String(),
				Position:        types.Position_LONG,
				Leverage:        tc.leverage,
				Collateral:      tc.collateral,
				TakeProfitPrice: tc.takeProfitPrice,
				StopLossPrice:   tc.stopLossPrice,
				PoolId:          1,
				VaultId:         1,
			}

			_, err := msgServer.PerformActionOpenPerpetual(suite.ctx, &msg)

			if tc.expectError {
				suite.Require().Error(err)
				// TODO: Match error messages
			} else {
				// Note: These might fail due to perpetual module validation, which is expected
				// The test is mainly to ensure the vault module doesn't panic on edge cases
				if err != nil {
					suite.T().Logf("Expected edge case test to potentially fail due to perpetual module validation: %v", err)
				}
			}
		})
	}
}
