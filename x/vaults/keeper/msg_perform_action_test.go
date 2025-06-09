package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"

	"github.com/elys-network/elys/v6/x/vaults/keeper"
	vaulttypes "github.com/elys-network/elys/v6/x/vaults/types"
)

func (suite *KeeperTestSuite) TestMsgServerPerformActionJoinPool() {
	// Create the vault first with the correct authority
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := vaulttypes.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "ustake",
		MaxAmountUsd:  math.LegacyMustNewDecFromStr("1000000"),
		AllowedCoins:  []string{"ustake", "uusdc", "uelys"},
		BenchmarkCoin: "uatom",
		RewardCoins:   []string{},
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	// Create test accounts
	manager := sdk.AccAddress("manager")
	invalidManager := sdk.AccAddress("invalid")

	// Update vault with manager
	vault, found := suite.app.VaultsKeeper.GetVault(suite.ctx, 1)
	suite.Require().True(found)
	vault.Manager = manager.String()
	suite.app.VaultsKeeper.SetVault(suite.ctx, vault)

	for _, tc := range []struct {
		desc        string
		msg         *vaulttypes.MsgPerformActionJoinPool
		setup       func()
		expectError bool
		errMsg      string
	}{
		{
			desc: "successful join pool",
			msg: &vaulttypes.MsgPerformActionJoinPool{
				Creator:        manager.String(),
				VaultId:        1,
				PoolId:         1,
				ShareAmountOut: math.NewInt(100),
				MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
			},
			setup: func() {
				coinsToSend := sdk.Coins{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000))}
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coinsToSend)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, manager, coinsToSend)
				suite.Require().NoError(err)

				coinsToSend = sdk.Coins{sdk.NewCoin(ptypes.Elys, math.NewInt(100000000))}
				err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coinsToSend)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, manager, coinsToSend)
				suite.Require().NoError(err)

				// Create a pool using CreateNewAmmPool
				suite.CreateNewAmmPool(
					manager,
					false,                                 // useOracle
					math.LegacyMustNewDecFromStr("0.003"), // swapFee
					math.LegacyMustNewDecFromStr("0.003"), // exitFee
					"uelys",                               // asset2
					math.NewInt(1000),                     // baseTokenAmount
					math.NewInt(1000),                     // assetAmount
				)

				// Fund the vault with tokens
				vaultAddress := vaulttypes.NewVaultAddress(1)
				coins := sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000)))
				err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", vaultAddress, coins)
				suite.Require().NoError(err)
			},
			expectError: false,
		},
		{
			desc: "invalid vault id",
			msg: &vaulttypes.MsgPerformActionJoinPool{
				Creator:        manager.String(),
				VaultId:        999, // Non-existent vault
				PoolId:         1,
				ShareAmountOut: math.NewInt(100),
				MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
			},
			expectError: true,
			errMsg:      "vault 999 not found",
		},
		{
			desc: "invalid signer",
			msg: &vaulttypes.MsgPerformActionJoinPool{
				Creator:        invalidManager.String(),
				VaultId:        1,
				PoolId:         1,
				ShareAmountOut: math.NewInt(100),
				MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
			},
			expectError: true,
			errMsg:      "vault 1 is not managed by",
		},
		{
			desc: "invalid action - zero pool id",
			msg: &vaulttypes.MsgPerformActionJoinPool{
				Creator:        manager.String(),
				VaultId:        1,
				PoolId:         0,
				ShareAmountOut: math.NewInt(100),
				MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
		{
			desc: "invalid action - empty max amounts in",
			msg: &vaulttypes.MsgPerformActionJoinPool{
				Creator:        manager.String(),
				VaultId:        1,
				PoolId:         1,
				ShareAmountOut: math.NewInt(100),
				MaxAmountsIn:   []sdk.Coin{},
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
		{
			desc: "invalid action - zero share amount out",
			msg: &vaulttypes.MsgPerformActionJoinPool{
				Creator:        manager.String(),
				VaultId:        1,
				PoolId:         1,
				ShareAmountOut: math.ZeroInt(),
				MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
	} {
		suite.Run(tc.desc, func() {
			// Setup test case if needed
			if tc.setup != nil {
				tc.setup()
			}

			_, err := msgServer.PerformActionJoinPool(suite.ctx, tc.msg)
			if tc.expectError {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errMsg)
			} else {
				suite.Require().NoError(err)

				// For successful join pool, verify the vault's balance decreased
				if tc.desc == "successful join pool" {
					vaultAddress := vaulttypes.NewVaultAddress(1)
					balance := suite.app.BankKeeper.GetBalance(suite.ctx, vaultAddress, "uusdc")
					suite.Require().True(balance.Amount.LT(math.NewInt(1000)), "vault balance should have decreased")
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerPerformActionExitPool() {
	// Create the vault first with the correct authority
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := vaulttypes.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "ustake",
		MaxAmountUsd:  math.LegacyMustNewDecFromStr("1000000"),
		AllowedCoins:  []string{"ustake", "uusdc", "uelys"},
		BenchmarkCoin: "uatom",
		RewardCoins:   []string{},
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	// Create test accounts
	manager := sdk.AccAddress("manager")
	invalidManager := sdk.AccAddress("invalid")

	// Update vault with manager
	vault, found := suite.app.VaultsKeeper.GetVault(suite.ctx, 1)
	suite.Require().True(found)
	vault.Manager = manager.String()
	suite.app.VaultsKeeper.SetVault(suite.ctx, vault)

	for _, tc := range []struct {
		desc        string
		msg         *vaulttypes.MsgPerformActionExitPool
		setup       func()
		expectError bool
		errMsg      string
	}{
		{
			desc: "successful exit pool",
			msg: &vaulttypes.MsgPerformActionExitPool{
				Creator:       manager.String(),
				VaultId:       1,
				PoolId:        1,
				ShareAmountIn: math.NewInt(100),
				MinAmountsOut: []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(90)}},
				TokenOutDenom: "uusdc",
			},
			setup: func() {
				// Create a pool using CreateNewAmmPool
				suite.CreateNewAmmPool(
					manager,
					false,                                 // useOracle
					math.LegacyMustNewDecFromStr("0.003"), // swapFee
					math.LegacyMustNewDecFromStr("0.003"), // exitFee
					"uelys",                               // asset2
					math.NewInt(1000),                     // baseTokenAmount
					math.NewInt(1000),                     // assetAmount
				)

				// Fund the vault with pool shares
				vaultAddress := vaulttypes.NewVaultAddress(1)
				coins := sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000)))
				err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", vaultAddress, coins)
				suite.Require().NoError(err)

				// Join the pool first
				_, _, err := suite.app.AmmKeeper.JoinPoolNoSwap(suite.ctx, vaultAddress, 1, math.NewInt(100), []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}})
				suite.Require().NoError(err)
			},
			expectError: false,
		},
		{
			desc: "invalid vault id",
			msg: &vaulttypes.MsgPerformActionExitPool{
				Creator:       manager.String(),
				VaultId:       999, // Non-existent vault
				PoolId:        1,
				ShareAmountIn: math.NewInt(100),
				MinAmountsOut: []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(90)}},
				TokenOutDenom: "uusdc",
			},
			expectError: true,
			errMsg:      "vault 999 not found",
		},
		{
			desc: "invalid signer",
			msg: &vaulttypes.MsgPerformActionExitPool{
				Creator:       invalidManager.String(),
				VaultId:       1,
				PoolId:        1,
				ShareAmountIn: math.NewInt(100),
				MinAmountsOut: []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(90)}},
				TokenOutDenom: "uusdc",
			},
			expectError: true,
			errMsg:      "vault 1 is not managed by",
		},
		{
			desc: "invalid action - zero pool id",
			msg: &vaulttypes.MsgPerformActionExitPool{
				Creator:       manager.String(),
				VaultId:       1,
				PoolId:        0,
				ShareAmountIn: math.NewInt(100),
				MinAmountsOut: []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(90)}},
				TokenOutDenom: "uusdc",
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
		{
			desc: "invalid action - empty min amounts out",
			msg: &vaulttypes.MsgPerformActionExitPool{
				Creator:       manager.String(),
				VaultId:       1,
				PoolId:        1,
				ShareAmountIn: math.NewInt(100),
				MinAmountsOut: []sdk.Coin{},
				TokenOutDenom: "uusdc",
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
		{
			desc: "invalid action - zero share amount in",
			msg: &vaulttypes.MsgPerformActionExitPool{
				Creator:       manager.String(),
				VaultId:       1,
				PoolId:        1,
				ShareAmountIn: math.ZeroInt(),
				MinAmountsOut: []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(90)}},
				TokenOutDenom: "uusdc",
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
	} {
		suite.Run(tc.desc, func() {
			// Setup test case if needed
			if tc.setup != nil {
				tc.setup()
			}

			_, err := msgServer.PerformActionExitPool(suite.ctx, tc.msg)
			if tc.expectError {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errMsg)
			} else {
				suite.Require().NoError(err)

				// For successful exit pool, verify the vault's balance increased
				if tc.desc == "successful exit pool" {
					vaultAddress := vaulttypes.NewVaultAddress(1)
					balance := suite.app.BankKeeper.GetBalance(suite.ctx, vaultAddress, "uusdc")
					suite.Require().True(balance.Amount.GT(math.NewInt(1000)), "vault balance should have increased")
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerPerformActionSwapByDenom() {
	// Create the vault first with the correct authority
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := vaulttypes.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "ustake",
		MaxAmountUsd:  math.LegacyMustNewDecFromStr("1000000"),
		AllowedCoins:  []string{"ustake", "uusdc", "uelys"},
		BenchmarkCoin: "uatom",
		RewardCoins:   []string{},
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	// Create test accounts
	manager := sdk.AccAddress("manager")
	invalidManager := sdk.AccAddress("invalid")

	// Update vault with manager
	vault, found := suite.app.VaultsKeeper.GetVault(suite.ctx, 1)
	suite.Require().True(found)
	vault.Manager = manager.String()
	suite.app.VaultsKeeper.SetVault(suite.ctx, vault)

	for _, tc := range []struct {
		desc        string
		msg         *vaulttypes.MsgPerformActionSwapByDenom
		setup       func()
		expectError bool
		errMsg      string
	}{
		{
			desc: "successful swap by denom",
			msg: &vaulttypes.MsgPerformActionSwapByDenom{
				Creator:   manager.String(),
				VaultId:   1,
				Amount:    sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)},
				MinAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(90)},
				MaxAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(110)},
				DenomIn:   "uusdc",
				DenomOut:  "uelys",
			},
			setup: func() {
				// Create a pool using CreateNewAmmPool
				suite.CreateNewAmmPool(
					manager,
					false,                                 // useOracle
					math.LegacyMustNewDecFromStr("0.003"), // swapFee
					math.LegacyMustNewDecFromStr("0.003"), // exitFee
					"uelys",                               // asset2
					math.NewInt(1000),                     // baseTokenAmount
					math.NewInt(1000),                     // assetAmount
				)

				// Fund the vault with tokens
				vaultAddress := vaulttypes.NewVaultAddress(1)
				coins := sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000)))
				err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", vaultAddress, coins)
				suite.Require().NoError(err)
			},
			expectError: false,
		},
		{
			desc: "invalid vault id",
			msg: &vaulttypes.MsgPerformActionSwapByDenom{
				Creator:   manager.String(),
				VaultId:   999, // Non-existent vault
				Amount:    sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)},
				MinAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(90)},
				MaxAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(110)},
				DenomIn:   "uusdc",
				DenomOut:  "uelys",
			},
			expectError: true,
			errMsg:      "vault 999 not found",
		},
		{
			desc: "invalid signer",
			msg: &vaulttypes.MsgPerformActionSwapByDenom{
				Creator:   invalidManager.String(),
				VaultId:   1,
				Amount:    sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)},
				MinAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(90)},
				MaxAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(110)},
				DenomIn:   "uusdc",
				DenomOut:  "uelys",
			},
			expectError: true,
			errMsg:      "vault 1 is not managed by",
		},
		{
			desc: "invalid action - zero amount",
			msg: &vaulttypes.MsgPerformActionSwapByDenom{
				Creator:   manager.String(),
				VaultId:   1,
				Amount:    sdk.Coin{Denom: "uusdc", Amount: math.ZeroInt()},
				MinAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(90)},
				MaxAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(110)},
				DenomIn:   "uusdc",
				DenomOut:  "uelys",
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
		{
			desc: "invalid action - zero min amount",
			msg: &vaulttypes.MsgPerformActionSwapByDenom{
				Creator:   manager.String(),
				VaultId:   1,
				Amount:    sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)},
				MinAmount: sdk.Coin{Denom: "uelys", Amount: math.ZeroInt()},
				MaxAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(110)},
				DenomIn:   "uusdc",
				DenomOut:  "uelys",
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
		{
			desc: "invalid action - zero max amount",
			msg: &vaulttypes.MsgPerformActionSwapByDenom{
				Creator:   manager.String(),
				VaultId:   1,
				Amount:    sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)},
				MinAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(90)},
				MaxAmount: sdk.Coin{Denom: "uelys", Amount: math.ZeroInt()},
				DenomIn:   "uusdc",
				DenomOut:  "uelys",
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
		{
			desc: "invalid action - empty denoms",
			msg: &vaulttypes.MsgPerformActionSwapByDenom{
				Creator:   manager.String(),
				VaultId:   1,
				Amount:    sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)},
				MinAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(90)},
				MaxAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(110)},
				DenomIn:   "",
				DenomOut:  "",
			},
			expectError: true,
			errMsg:      "action failed with error",
		},
	} {
		suite.Run(tc.desc, func() {
			// Setup test case if needed
			if tc.setup != nil {
				tc.setup()
			}

			_, err := msgServer.PerformActionSwapByDenom(suite.ctx, tc.msg)
			if tc.expectError {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errMsg)
			} else {
				suite.Require().NoError(err)

				// For successful swap, verify the vault's balance changed
				if tc.desc == "successful swap by denom" {
					vaultAddress := vaulttypes.NewVaultAddress(1)
					balance := suite.app.BankKeeper.GetBalance(suite.ctx, vaultAddress, "uusdc")
					suite.Require().True(balance.Amount.LT(math.NewInt(1000)), "vault balance should have decreased")
					balance = suite.app.BankKeeper.GetBalance(suite.ctx, vaultAddress, "uelys")
					suite.Require().True(balance.Amount.GT(math.ZeroInt()), "vault should have received uelys tokens")
				}
			}
		})
	}
}
