package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/vaults/keeper"
	vaulttypes "github.com/elys-network/elys/v7/x/vaults/types"
)

func (suite *KeeperTestSuite) TestVaultFlow() {
	// Create test accounts
	manager := sdk.AccAddress([]byte("manager"))
	depositor := sdk.AccAddress([]byte("depositor"))

	// Setup initial balances for depositor
	coinsToSend := sdk.Coins{sdk.NewCoin("uusdc", sdkmath.NewInt(1000000))}
	err := suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coinsToSend)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", depositor, coinsToSend)
	suite.Require().NoError(err)

	err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coinsToSend)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, manager, coinsToSend)
	suite.Require().NoError(err)

	coinsToSend = sdk.Coins{sdk.NewCoin("uatom", sdkmath.NewInt(1000000))}
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coinsToSend)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", manager, coinsToSend)
	suite.Require().NoError(err)

	// Step 1: Add a new vault
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := vaulttypes.MsgAddVault{
		Creator:        authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:   "uusdc",
		MaxAmountUsd:   sdkmath.LegacyNewDec(1000000),
		AllowedCoins:   []string{"uusdc", "uatom", "amm/pool/1"},
		RewardCoins:    []string{"uelys"},
		BenchmarkCoin:  "uatom",
		Manager:        manager.String(),
		AllowedActions: []string{vaulttypes.ActionJoinPool, vaulttypes.ActionExitPool, vaulttypes.ActionSwap},
	}
	_, err = msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	// Verify vault was created correctly
	vault, found := suite.app.VaultsKeeper.GetVault(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal("uusdc", vault.DepositDenom)
	suite.Require().Equal(manager.String(), vault.Manager)

	// Step 2: Deposit to vault
	depositMsg := vaulttypes.MsgDeposit{
		VaultId:   1,
		Depositor: depositor.String(),
		Amount:    sdk.NewCoin("uusdc", sdkmath.NewInt(100000)),
	}
	_, err = msgServer.Deposit(suite.ctx, &depositMsg)
	suite.Require().NoError(err)

	// Verify deposit was successful
	vaultAddress := vaulttypes.NewVaultAddress(1)
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, vaultAddress, "uusdc")
	suite.Require().Equal(sdkmath.NewInt(100000), balance.Amount)

	// Verify depositor received share tokens
	shareDenom := vaulttypes.GetShareDenomForVault(1)
	commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, depositor)
	committedAmount := commitments.GetCommittedAmountForDenom(shareDenom)
	suite.Require().True(committedAmount.GT(sdkmath.ZeroInt()), "depositor should have received share tokens")

	// Step 3: Create a pool for the vault to join
	suite.CreateNewAmmPool(
		manager,
		false,                                    // useOracle
		sdkmath.LegacyMustNewDecFromStr("0.003"), // swapFee
		sdkmath.LegacyMustNewDecFromStr("0.003"), // exitFee
		"uatom",                                  // asset2
		sdkmath.NewInt(1000),                     // baseTokenAmount
		sdkmath.NewInt(1000),                     // assetAmount
	)

	// Step 4: Perform join pool action
	joinPoolMsg := vaulttypes.MsgPerformActionJoinPool{
		Creator:        manager.String(),
		VaultId:        1,
		PoolId:         1,
		ShareAmountOut: sdkmath.NewInt(100),
		MaxAmountsIn:   []sdk.Coin{{Denom: "uusdc", Amount: sdkmath.NewInt(100)}},
	}

	_, err = msgServer.PerformActionJoinPool(suite.ctx, &joinPoolMsg)
	suite.Require().NoError(err)

	// Verify vault's balance decreased after joining pool
	balance = suite.app.BankKeeper.GetBalance(suite.ctx, vaultAddress, "uusdc")
	suite.Require().True(balance.Amount.LT(sdkmath.NewInt(100000)), "vault balance should have decreased after joining pool")

	beforeWithdraw := suite.app.BankKeeper.GetAllBalances(suite.ctx, depositor)

	// Withdraw half shares
	exitVaultMsg := vaulttypes.MsgWithdraw{
		Withdrawer: depositor.String(),
		VaultId:    1,
		Shares:     sdkmath.NewInt(50000),
	}

	_, err = msgServer.Withdraw(suite.ctx, &exitVaultMsg)
	suite.Require().NoError(err)

	// Check balances
	afterWithdraw := suite.app.BankKeeper.GetAllBalances(suite.ctx, depositor)
	addedCoins := afterWithdraw.Sub(beforeWithdraw...)
	// TODO: verify numbers as per pool shares
	suite.Require().Equal(addedCoins, sdk.Coins{sdk.NewCoin("uatom", sdkmath.NewInt(23)), sdk.NewCoin("uusdc", sdkmath.NewInt(49975))})
}
