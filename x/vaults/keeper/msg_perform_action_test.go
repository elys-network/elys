package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	"github.com/elys-network/elys/x/vaults/keeper"
	vaulttypes "github.com/elys-network/elys/x/vaults/types"
)

func (suite *KeeperTestSuite) TestMsgServerPerformAction() {
	// Create the vault first with the correct authority
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := vaulttypes.MsgAddVault{
		Creator:        authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:   "ustake",
		MaxAmountUsd:   math.LegacyMustNewDecFromStr("1000000"),
		AllowedCoins:   []string{"ustake", "uusdc", "uelys"},
		AllowedActions: []uint64{},
		RewardCoins:    []string{},
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
		msg         *vaulttypes.MsgPerformAction
		setup       func()
		expectError bool
		errMsg      string
	}{
		{
			desc: "successful join pool",
			msg: &vaulttypes.MsgPerformAction{
				Creator: manager.String(),
				VaultId: 1,
				Action: &vaulttypes.MsgPerformAction_JoinPool{
					JoinPool: &types.MsgJoinPool{
						Sender:         manager.String(),
						PoolId:         1,
						ShareAmountOut: math.NewInt(100),
						MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
					},
				},
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
			msg: &vaulttypes.MsgPerformAction{
				Creator: manager.String(),
				VaultId: 999, // Non-existent vault
				Action: &vaulttypes.MsgPerformAction_JoinPool{
					JoinPool: &types.MsgJoinPool{
						Sender:         manager.String(),
						PoolId:         1,
						ShareAmountOut: math.NewInt(100),
						MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
					},
				},
			},
			expectError: true,
			errMsg:      "vault 999 not found",
		},
		{
			desc: "invalid signer",
			msg: &vaulttypes.MsgPerformAction{
				Creator: invalidManager.String(),
				VaultId: 1,
				Action: &vaulttypes.MsgPerformAction_JoinPool{
					JoinPool: &types.MsgJoinPool{
						Sender:         invalidManager.String(),
						PoolId:         1,
						ShareAmountOut: math.NewInt(100),
						MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
					},
				},
			},
			expectError: true,
			errMsg:      "vault 1 is not managed by",
		},
		{
			desc: "invalid action - nil join pool",
			msg: &vaulttypes.MsgPerformAction{
				Creator: manager.String(),
				VaultId: 1,
				Action: &vaulttypes.MsgPerformAction_JoinPool{
					JoinPool: nil,
				},
			},
			expectError: true,
			errMsg:      "vault 1 does not allow this action",
		},
		{
			desc: "invalid action - zero pool id",
			msg: &vaulttypes.MsgPerformAction{
				Creator: manager.String(),
				VaultId: 1,
				Action: &vaulttypes.MsgPerformAction_JoinPool{
					JoinPool: &types.MsgJoinPool{
						Sender:         manager.String(),
						PoolId:         0,
						ShareAmountOut: math.NewInt(100),
						MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
					},
				},
			},
			expectError: true,
			errMsg:      "vault 1 does not allow this action",
		},
		{
			desc: "invalid action - empty max amounts in",
			msg: &vaulttypes.MsgPerformAction{
				Creator: manager.String(),
				VaultId: 1,
				Action: &vaulttypes.MsgPerformAction_JoinPool{
					JoinPool: &types.MsgJoinPool{
						Sender:         manager.String(),
						PoolId:         1,
						ShareAmountOut: math.NewInt(100),
						MaxAmountsIn:   []sdk.Coin{},
					},
				},
			},
			expectError: true,
			errMsg:      "vault 1 does not allow this action",
		},
		{
			desc: "invalid action - zero share amount out",
			msg: &vaulttypes.MsgPerformAction{
				Creator: manager.String(),
				VaultId: 1,
				Action: &vaulttypes.MsgPerformAction_JoinPool{
					JoinPool: &types.MsgJoinPool{
						Sender:         manager.String(),
						PoolId:         1,
						ShareAmountOut: math.ZeroInt(),
						MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: math.NewInt(100)}},
					},
				},
			},
			expectError: true,
			errMsg:      "vault 1 does not allow this action",
		},
		{
			desc: "invalid action - invalid swap by denom",
			msg: &vaulttypes.MsgPerformAction{
				Creator: manager.String(),
				VaultId: 1,
				Action: &vaulttypes.MsgPerformAction_SwapByDenom{
					SwapByDenom: &types.MsgSwapByDenom{
						Sender:    manager.String(),
						Amount:    sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)},
						MinAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(90)},
						MaxAmount: sdk.Coin{Denom: "uelys", Amount: math.NewInt(110)},
						DenomIn:   "uusdc",
						DenomOut:  "uelys",
						Recipient: "invalid", // Should be vault address
					},
				},
			},
			expectError: true,
			errMsg:      "vault 1 does not allow this action",
		},
	} {
		suite.Run(tc.desc, func() {
			// Setup test case if needed
			if tc.setup != nil {
				tc.setup()
			}

			_, err := msgServer.PerformAction(suite.ctx, tc.msg)
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
