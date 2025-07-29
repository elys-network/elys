package keeper_test

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/x/masterchef/types"
)

func (suite *MasterchefKeeperTestSuite) TestTogglePoolEdenRewards() {

	pk := ed25519.GenPrivKey().PubKey()
	nonAuthority := sdk.AccAddress(pk.Address()).String()
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	suite.AddPoolInfo()

	testCases := []struct {
		name         string
		msg          *types.MsgTogglePoolEdenRewards
		expectErr    bool
		expectErrMsg error
	}{
		{
			name: "Invalid signer",
			msg: &types.MsgTogglePoolEdenRewards{
				Authority: "",
				PoolId:    2,
			},
			expectErr:    true,
			expectErrMsg: errors.New("invalid authority"),
		},
		{
			name: "Invalid Signer Address",
			msg: &types.MsgTogglePoolEdenRewards{
				Authority: nonAuthority,
				PoolId:    2,
			},
			expectErr:    true,
			expectErrMsg: errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, nonAuthority),
		},
		{
			name: "Pool not found",
			msg: &types.MsgTogglePoolEdenRewards{
				Authority: authority,
				PoolId:    3,
			},
			expectErr:    true,
			expectErrMsg: types.ErrPoolNotFound,
		},
		{
			name: "Happy Flow",
			msg: &types.MsgTogglePoolEdenRewards{
				Authority: authority,
				PoolId:    2,
				Enable:    true,
			},
			expectErr:    false,
			expectErrMsg: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.TogglePoolEdenRewards(suite.ctx, tc.msg)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
