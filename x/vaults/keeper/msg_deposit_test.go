package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/vaults/keeper"
	"github.com/elys-network/elys/x/vaults/types"
)

func (suite *KeeperTestSuite) TestMsgServerDeposit() {
	for _, tc := range []struct {
		desc        string
		vaultId     uint64
		depositer   sdk.AccAddress
		amount      sdk.Coin
		setup       func()
		expectError bool
	}{
		{
			desc:        "successful deposit",
			vaultId:     1,
			depositer:   sdk.AccAddress([]byte("depositer1")),
			amount:      sdk.NewCoin("ustake", sdkmath.NewInt(1000)),
			setup:       func() {},
			expectError: false,
		},
		// {
		// 	desc:        "vault not found",
		// 	vaultId:     999,
		// 	depositer:   sdk.AccAddress([]byte("depositer2")),
		// 	amount:      sdk.NewCoin("stake", sdkmath.NewInt(1000)),
		// 	setup:       func() {}, // No vault setup
		// 	expectError: true,
		// },
		// {
		// 	desc:      "invalid coin denom",
		// 	vaultId:   1,
		// 	depositer: sdk.AccAddress([]byte("depositer3")),
		// 	amount:    sdk.NewCoin("invalid", sdkmath.NewInt(1000)),
		// 	setup: func() {
		// 		// Create a vault
		// 		vault := types.Vault{
		// 			Id:           1,
		// 			AllowedCoins: []string{"stake"},
		// 		}
		// 		suite.app.VaultsKeeper.SetVault(suite.ctx, vault)
		// 	},
		// 	expectError: true,
		// },
		// {
		// 	desc:      "insufficient funds",
		// 	vaultId:   1,
		// 	depositer: sdk.AccAddress([]byte("depositer4")),
		// 	amount:    sdk.NewCoin("stake", sdkmath.NewInt(1000)),
		// 	setup: func() {
		// 		// Create a vault
		// 		vault := types.Vault{
		// 			Id:           1,
		// 			AllowedCoins: []string{"stake"},
		// 		}
		// 		suite.app.VaultsKeeper.SetVault(suite.ctx, vault)

		// 		// bootstrap balances
		// 		tokens := sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(500)))
		// 		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tokens)
		// 		suite.Require().NoError(err)
		// 		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sdk.AccAddress([]byte("depositer4")), tokens)
		// 		suite.Require().NoError(err)
		// 	},
		// 	expectError: true,
		// },
	} {
		suite.Run(tc.desc, func() {
			// Setup test case
			tc.setup()

			msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)

			addVault := types.MsgAddVault{
				Creator:        suite.app.VaultsKeeper.GetAuthority(),
				DepositDenom:   tc.amount.Denom,
				MaxAmountUsd:   sdkmath.LegacyNewDec(1000000),
				AllowedCoins:   []string{tc.amount.Denom},
				AllowedActions: []uint64{},
				RewardCoins:    []string{},
			}
			_, err := msgServer.AddVault(suite.ctx, &addVault)

			// Prepare the message
			msg := types.MsgDeposit{
				VaultId:   tc.vaultId,
				Depositor: tc.depositer.String(),
				Amount:    tc.amount,
			}
			_, err = msgServer.Deposit(suite.ctx, &msg)

			if tc.expectError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// Verify the deposit was successful
				vaultAddress := types.NewVaultAddress(tc.vaultId)
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, vaultAddress)
				suite.Require().NotNil(commitments)
				found := false
				for _, token := range commitments.CommittedTokens {
					if token.Denom == tc.amount.Denom && token.Amount.Equal(tc.amount.Amount) {
						found = true
						break
					}
				}
				suite.Require().True(found, "deposit not found in commitments")
			}
		})
	}
}
