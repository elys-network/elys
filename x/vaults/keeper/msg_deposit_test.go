package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/vaults/keeper"
	"github.com/elys-network/elys/x/vaults/types"
)

func (suite *KeeperTestSuite) TestMsgServerDeposit() {
	// Create the vault first with the correct authority
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := types.MsgAddVault{
		Creator:        authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:   "ustake",
		MaxAmountUsd:   sdkmath.LegacyNewDec(1000000),
		AllowedCoins:   []string{"ustake"},
		AllowedActions: []uint64{},
		RewardCoins:    []string{},
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

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
			desc:        "invalid coin denom",
			vaultId:     1,
			depositer:   sdk.AccAddress([]byte("depositer3")),
			amount:      sdk.NewCoin("invalid", sdkmath.NewInt(1000)),
			setup:       func() {},
			expectError: true,
		},
		{
			desc:      "insufficient funds",
			vaultId:   1,
			depositer: sdk.AccAddress([]byte("depositer4")),
			amount:    sdk.NewCoin("stake", sdkmath.NewInt(100000)),
			setup: func() {
			},
			expectError: true,
		},
	} {
		suite.Run(tc.desc, func() {
			// Setup test case
			tc.setup = func() {
				// Mint coins for the depositer
				coins := sdk.NewCoins(tc.amount)
				err := suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", tc.depositer, coins)
				suite.Require().NoError(err)
			}
			tc.setup()

			msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)

			// Prepare the message
			msg := types.MsgDeposit{
				VaultId:   tc.vaultId,
				Depositor: tc.depositer.String(),
				Amount:    tc.amount,
			}
			_, err := msgServer.Deposit(suite.ctx, &msg)

			if tc.expectError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// Verify the deposit was successful by checking the vault's balance
				vaultAddress := types.NewVaultAddress(tc.vaultId)
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, vaultAddress, tc.amount.Denom)
				suite.Require().Equal(tc.amount.Amount, balance.Amount)

				// Verify the depositer received the share tokens
				shareDenom := types.GetShareDenomForVault(tc.vaultId)
				shareBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.depositer, shareDenom)
				suite.Require().True(shareBalance.Amount.GT(sdkmath.ZeroInt()), "depositer should have received share tokens")
			}
		})
	}
}
