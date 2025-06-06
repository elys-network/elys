package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/elys-network/elys/v6/x/vaults/keeper"
	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (suite *KeeperTestSuite) TestVaultQuery() {
	// Create test vault
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := types.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "uusdc",
		MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		AllowedCoins:  []string{"uusdc", "uatom"},
		RewardCoins:   []string{"uelys"},
		BenchmarkCoin: "uatom",
		Manager:       sdk.AccAddress([]byte("manager")).String(),
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	testCases := []struct {
		name      string
		req       *types.QueryVaultRequest
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "invalid request",
			req:       nil,
			expErr:    true,
			expErrMsg: "invalid request",
		},
		{
			name:      "vault not found",
			req:       &types.QueryVaultRequest{VaultId: 999},
			expErr:    true,
			expErrMsg: "vault not found",
		},
		{
			name:   "successful query",
			req:    &types.QueryVaultRequest{VaultId: 1},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.app.VaultsKeeper.Vault(suite.ctx, tc.req)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				suite.Require().Equal(uint64(1), resp.Vault.Id)
				suite.Require().Equal("uusdc", resp.Vault.DepositDenom)
				suite.Require().Equal(sdkmath.LegacyNewDec(1000000), resp.Vault.MaxAmountUsd)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestVaultsQuery() {
	// Create test vaults
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)

	// Create first vault
	addVault1 := types.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "uusdc",
		MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		AllowedCoins:  []string{"uusdc", "uatom"},
		RewardCoins:   []string{"uelys"},
		BenchmarkCoin: "uatom",
		Manager:       sdk.AccAddress([]byte("manager1")).String(),
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault1)
	suite.Require().NoError(err)

	// Create second vault
	addVault2 := types.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "uatom",
		MaxAmountUsd:  sdkmath.LegacyNewDec(2000000),
		AllowedCoins:  []string{"uatom", "uusdc"},
		RewardCoins:   []string{"uelys"},
		BenchmarkCoin: "uusdc",
		Manager:       sdk.AccAddress([]byte("manager2")).String(),
	}
	_, err = msgServer.AddVault(suite.ctx, &addVault2)
	suite.Require().NoError(err)

	testCases := []struct {
		name      string
		req       *types.QueryVaultsRequest
		expErr    bool
		expErrMsg string
		expCount  int
	}{
		{
			name:      "invalid request",
			req:       nil,
			expErr:    true,
			expErrMsg: "invalid request",
		},
		{
			name:     "successful query",
			req:      &types.QueryVaultsRequest{},
			expErr:   false,
			expCount: 2,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.app.VaultsKeeper.Vaults(suite.ctx, tc.req)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				suite.Require().Len(resp.Vaults, tc.expCount)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestVaultValueQuery() {
	// Create test vault
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := types.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "uusdc",
		MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		AllowedCoins:  []string{"uusdc", "uatom"},
		RewardCoins:   []string{"uelys"},
		BenchmarkCoin: "uatom",
		Manager:       sdk.AccAddress([]byte("manager")).String(),
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	testCases := []struct {
		name      string
		req       *types.QueryVaultValue
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "invalid request",
			req:       nil,
			expErr:    true,
			expErrMsg: "invalid request",
		},
		{
			name:      "vault not found",
			req:       &types.QueryVaultValue{VaultId: 999},
			expErr:    true,
			expErrMsg: "vault not found",
		},
		{
			name:   "successful query",
			req:    &types.QueryVaultValue{VaultId: 1},
			expErr: false,
		},
		// TODO: add case with deposits and rewards
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.app.VaultsKeeper.VaultValue(suite.ctx, tc.req)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				// Note: USD value might be 0 if no deposits or if price feed is not set up
				suite.Require().NotNil(resp.UsdValue)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestVaultPositionsQuery() {
	// Create test vault
	msgServer := keeper.NewMsgServerImpl(suite.app.VaultsKeeper)
	addVault := types.MsgAddVault{
		Creator:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		DepositDenom:  "uusdc",
		MaxAmountUsd:  sdkmath.LegacyNewDec(1000000),
		AllowedCoins:  []string{"uusdc", "uatom"},
		RewardCoins:   []string{"uelys"},
		BenchmarkCoin: "uatom",
		Manager:       sdk.AccAddress([]byte("manager")).String(),
	}
	_, err := msgServer.AddVault(suite.ctx, &addVault)
	suite.Require().NoError(err)

	testCases := []struct {
		name      string
		req       *types.QueryVaultPositionsRequest
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "invalid request",
			req:       nil,
			expErr:    true,
			expErrMsg: "invalid request",
		},
		{
			name:   "successful query",
			req:    &types.QueryVaultPositionsRequest{VaultId: 1},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.app.VaultsKeeper.VaultPositions(suite.ctx, tc.req)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				suite.Require().NotNil(resp.Positions)
			}
		})
	}
}
