package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
			desc:      "successful deposit",
			vaultId:   1,
			depositer: sdk.AccAddress([]byte("depositer1")),
			amount:    sdk.NewCoin("stake", sdkmath.NewInt(1000)),
			setup: func() {
				// Create a vault
				vault := types.Vault{
					Id:           1,
					AllowedCoins: []string{"stake"},
				}
				suite.app.VaultsKeeper.SetVault(suite.ctx, vault)
			},
			expectError: false,
		},
		{
			desc:        "vault not found",
			vaultId:     999,
			depositer:   sdk.AccAddress([]byte("depositer2")),
			amount:      sdk.NewCoin("stake", sdkmath.NewInt(1000)),
			setup:       func() {}, // No vault setup
			expectError: true,
		},
		{
			desc:      "invalid coin denom",
			vaultId:   1,
			depositer: sdk.AccAddress([]byte("depositer3")),
			amount:    sdk.NewCoin("invalid", sdkmath.NewInt(1000)),
			setup: func() {
				// Create a vault
				vault := types.Vault{
					Id:           1,
					AllowedCoins: []string{"stake"},
				}
				suite.app.VaultsKeeper.SetVault(suite.ctx, vault)
			},
			expectError: true,
		},
		{
			desc:      "insufficient funds",
			vaultId:   1,
			depositer: sdk.AccAddress([]byte("depositer4")),
			amount:    sdk.NewCoin("stake", sdkmath.NewInt(1000)),
			setup: func() {
				// Create a vault
				vault := types.Vault{
					Id:           1,
					AllowedCoins: []string{"stake"},
				}
				suite.app.VaultsKeeper.SetVault(suite.ctx, vault)

				// Set depositer balance to less than the deposit amount
				suite.app.BankKeeper.SetBalances(suite.ctx, sdk.AccAddress([]byte("depositer4")), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(500))))
			},
			expectError: true,
		},
	} {
		suite.Run(tc.desc, func() {
			// Setup test case
			tc.setup()

			// Prepare the message
			msg := types.MsgDeposit{
				VaultId:   tc.vaultId,
				Depositor: tc.depositer.String(),
				Amount:    tc.amount,
			}

			// Call the handler
			_, err := suite.app.VaultsKeeper.Deposit(suite.ctx, &msg)

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
