package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/vaults/keeper"
	"github.com/elys-network/elys/v5/x/vaults/types"
)

func (suite *KeeperTestSuite) TestMsgServerAddVault() {
	for _, tc := range []struct {
		desc          string
		creator       sdk.AccAddress
		depositDenom  string
		maxAmountUsd  sdkmath.LegacyDec
		allowedCoins  []string
		rewardCoins   []string
		benchmarkCoin string
		Manager       sdk.AccAddress
		setup         func()
		expectError   bool
	}{
		{
			desc:          "successful vault creation",
			creator:       authtypes.NewModuleAddress(govtypes.ModuleName),
			depositDenom:  "uusdc",
			maxAmountUsd:  sdkmath.LegacyNewDec(1000000),
			allowedCoins:  []string{"uusdc", "uatom"},
			rewardCoins:   []string{"uelys"},
			benchmarkCoin: "uatom",
			Manager:       sdk.AccAddress([]byte("manager")),
			setup:         func() {},
			expectError:   false,
		},
		{
			desc:          "invalid authority",
			creator:       sdk.AccAddress([]byte("invalid_creator")),
			depositDenom:  "uusdc",
			maxAmountUsd:  sdkmath.LegacyNewDec(1000000),
			allowedCoins:  []string{"uusdc"},
			rewardCoins:   []string{"uelys"},
			benchmarkCoin: "uatom",
			setup:         func() {},
			expectError:   true,
		},
		{
			desc:          "empty deposit denom",
			creator:       sdk.AccAddress([]byte("creator2")),
			depositDenom:  "",
			maxAmountUsd:  sdkmath.LegacyNewDec(1000000),
			allowedCoins:  []string{"uusdc"},
			rewardCoins:   []string{"uelys"},
			benchmarkCoin: "uatom",
			setup:         func() {},
			expectError:   true,
		},
	} {
		suite.Run(tc.desc, func() {
			// Setup test case
			tc.setup()

			msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)

			// Prepare the message
			msg := types.MsgAddVault{
				Creator:       tc.creator.String(),
				DepositDenom:  tc.depositDenom,
				MaxAmountUsd:  tc.maxAmountUsd,
				AllowedCoins:  tc.allowedCoins,
				RewardCoins:   tc.rewardCoins,
				BenchmarkCoin: tc.benchmarkCoin,
				Manager:       tc.Manager.String(),
			}

			_, err := msgServer.AddVault(suite.ctx, &msg)

			if tc.expectError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// Verify vault was created correctly
				vault, found := suite.app.VaultsKeeper.GetVault(suite.ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(tc.depositDenom, vault.DepositDenom)
				suite.Require().Equal(tc.maxAmountUsd, vault.MaxAmountUsd)
				suite.Require().Equal(tc.allowedCoins, vault.AllowedCoins)
				suite.Require().Equal(tc.rewardCoins, vault.RewardCoins)
				suite.Require().Equal(tc.benchmarkCoin, vault.BenchmarkCoin)
				suite.Require().Equal(tc.Manager.String(), vault.Manager)

				// Verify module account was created
				vaultAddress := types.NewVaultAddress(1)
				account := suite.app.AccountKeeper.GetAccount(suite.ctx, vaultAddress)
				suite.Require().NotNil(account)
				suite.Require().Equal(vaultAddress, account.GetAddress())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerAddVault_SequentialIds() {
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)

	// Create multiple vaults and verify IDs are sequential
	for i := 1; i <= 5; i++ {
		msg := types.MsgAddVault{
			Creator:       suite.app.VaultsKeeper.GetAuthority(),
			DepositDenom:  "uusdc",
			MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
			AllowedCoins:  []string{"uusdc"},
			RewardCoins:   []string{"uelys"},
			BenchmarkCoin: "uatom",
		}

		_, err := msgServer.AddVault(suite.ctx, &msg)
		suite.Require().NoError(err)

		// Verify vault was created with correct ID
		vault, found := suite.app.VaultsKeeper.GetVault(suite.ctx, uint64(i))
		suite.Require().True(found)
		suite.Require().Equal(uint64(i), vault.Id)
	}
}
