package keeper_test

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/masterchef/types"
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
			expectErrMsg: fmt.Errorf("invalid authority"),
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
			expectErrMsg: fmt.Errorf("invalid authority"),
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
