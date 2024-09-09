package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateRewardsDataLifetime_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateRewardsDataLifetime
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateRewardsDataLifetime{
				Creator:             "invalid_address",
				RewardsDataLifetime: "1",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateRewardsDataLifetime{
				Creator:             sample.AccAddress(),
				RewardsDataLifetime: "1",
			},
		}, {
			name: "invalid reward lifecycle",
			msg: MsgUpdateRewardsDataLifetime{
				Creator:             sample.AccAddress(),
				RewardsDataLifetime: "abcd",
			},
			err: ErrInvalidRewardsDataLifecycle,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
