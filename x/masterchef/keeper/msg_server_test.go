package keeper_test

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v6/x/masterchef/types"
)

func (suite *MasterchefKeeperTestSuite) TestMsgUpdateParams() {

	pk := ed25519.GenPrivKey().PubKey()
	nonAuthority := sdk.AccAddress(pk.Address()).String()
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	suite.AddPoolInfo()

	testCases := []struct {
		name         string
		msg          *types.MsgUpdateParams
		expectErr    bool
		expectErrMsg error
	}{
		{
			name: "Invalid signer",
			msg: &types.MsgUpdateParams{
				Authority: "",
				Params:    types.DefaultParams(),
			},
			expectErr:    true,
			expectErrMsg: errors.New("invalid authority"),
		},
		{
			name: "Invalid Signer Address",
			msg: &types.MsgUpdateParams{
				Authority: nonAuthority,
				Params:    types.DefaultParams(),
			},
			expectErr:    true,
			expectErrMsg: errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, nonAuthority),
		},
		{
			name: "Happy Flow",
			msg: &types.MsgUpdateParams{
				Authority: authority,
				Params:    types.DefaultParams(),
			},
			expectErr:    false,
			expectErrMsg: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.UpdateParams(suite.ctx, tc.msg)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestMsgUpdatePoolMultipliers() {

	pk := ed25519.GenPrivKey().PubKey()
	nonAuthority := sdk.AccAddress(pk.Address()).String()
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	suite.AddPoolInfo()

	testCases := []struct {
		name         string
		msg          *types.MsgUpdatePoolMultipliers
		expectErr    bool
		expectErrMsg error
	}{
		{
			name: "Invalid signer",
			msg: &types.MsgUpdatePoolMultipliers{
				Authority: "",
				PoolMultipliers: []types.PoolMultiplier{
					{
						PoolId:     1,
						Multiplier: sdkmath.LegacyOneDec(),
					},
				},
			},
			expectErr:    true,
			expectErrMsg: errors.New("invalid authority"),
		},
		{
			name: "Invalid Signer Address",
			msg: &types.MsgUpdatePoolMultipliers{
				Authority: nonAuthority,
				PoolMultipliers: []types.PoolMultiplier{
					{
						PoolId:     1,
						Multiplier: sdkmath.LegacyOneDec(),
					},
				},
			},
			expectErr:    true,
			expectErrMsg: errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, nonAuthority),
		},
		{
			name: "Happy Flow",
			msg: &types.MsgUpdatePoolMultipliers{
				Authority: authority,
				PoolMultipliers: []types.PoolMultiplier{
					{
						PoolId:     1,
						Multiplier: sdkmath.LegacyOneDec(),
					},
				},
			},
			expectErr:    false,
			expectErrMsg: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.UpdatePoolMultipliers(suite.ctx, tc.msg)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestAddExternalRewardDenom() {

	pk := ed25519.GenPrivKey().PubKey()
	nonAuthority := sdk.AccAddress(pk.Address()).String()
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	suite.AddPoolInfo()

	testCases := []struct {
		name         string
		msg          *types.MsgAddExternalRewardDenom
		expectErr    bool
		expectErrMsg error
	}{
		{
			name: "Invalid signer",
			msg: &types.MsgAddExternalRewardDenom{
				Authority:   "",
				RewardDenom: "reward",
				MinAmount:   sdkmath.OneInt(),
				Supported:   true,
			},
			expectErr:    true,
			expectErrMsg: errors.New("invalid authority"),
		},
		{
			name: "Invalid Signer Address",
			msg: &types.MsgAddExternalRewardDenom{
				Authority:   nonAuthority,
				RewardDenom: "reward",
				MinAmount:   sdkmath.OneInt(),
				Supported:   true,
			},
			expectErr:    true,
			expectErrMsg: errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, nonAuthority),
		},
		{
			name: "Happy Flow",
			msg: &types.MsgAddExternalRewardDenom{
				Authority:   authority,
				RewardDenom: "reward",
				MinAmount:   sdkmath.OneInt(),
				Supported:   true,
			},
			expectErr:    false,
			expectErrMsg: nil,
		},
		{
			name: "Happy Flow",
			msg: &types.MsgAddExternalRewardDenom{
				Authority:   authority,
				RewardDenom: "reward",
				MinAmount:   sdkmath.OneInt(),
				Supported:   false,
			},
			expectErr:    false,
			expectErrMsg: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.AddExternalRewardDenom(suite.ctx, tc.msg)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestAddExternalIncentive() {

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	suite.AddPoolInfo()

	testCases := []struct {
		name         string
		msg          *types.MsgAddExternalIncentive
		expectErr    bool
		expectErrMsg error
	}{
		{
			name: "Zero amount per block",
			msg: &types.MsgAddExternalIncentive{
				Sender:         authority,
				RewardDenom:    "reward",
				PoolId:         1,
				FromBlock:      1,
				ToBlock:        100,
				AmountPerBlock: sdkmath.ZeroInt(),
			},
			expectErr:    true,
			expectErrMsg: errors.New("invalid amount per block"),
		},
		{
			name: "Invalid From block",
			msg: &types.MsgAddExternalIncentive{
				Sender:         authority,
				RewardDenom:    "reward",
				PoolId:         1,
				FromBlock:      suite.ctx.BlockHeight() - 100,
				ToBlock:        100,
				AmountPerBlock: sdkmath.ZeroInt(),
			},
			expectErr:    true,
			expectErrMsg: errors.New("invalid from block"),
		},
		{
			name: "Invalid block range",
			msg: &types.MsgAddExternalIncentive{
				Sender:         authority,
				RewardDenom:    "reward",
				PoolId:         1,
				FromBlock:      100,
				ToBlock:        10,
				AmountPerBlock: sdkmath.ZeroInt(),
			},
			expectErr:    true,
			expectErrMsg: errors.New("invalid block range"),
		},
		{
			name: "Reward denom not supported",
			msg: &types.MsgAddExternalIncentive{
				Sender:         authority,
				RewardDenom:    "reward",
				PoolId:         1,
				FromBlock:      1,
				ToBlock:        100,
				AmountPerBlock: sdkmath.OneInt(),
			},
			expectErr:    true,
			expectErrMsg: errors.New("invalid reward denom"),
		},
		{
			name: "Reward denom amount low",
			msg: &types.MsgAddExternalIncentive{
				Sender:         authority,
				RewardDenom:    "reward1",
				PoolId:         1,
				FromBlock:      1,
				ToBlock:        100,
				AmountPerBlock: sdkmath.OneInt(),
			},
			expectErr:    true,
			expectErrMsg: errors.New("too small amount"),
		},
	}

	params := suite.app.MasterchefKeeper.GetParams(suite.ctx)
	params.SupportedRewardDenoms = []*types.SupportedRewardDenom{{Denom: "reward1", MinAmount: sdkmath.NewInt(10000)}}
	suite.app.MasterchefKeeper.SetParams(suite.ctx, params)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.AddExternalIncentive(suite.ctx, tc.msg)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
