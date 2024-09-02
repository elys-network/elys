package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgTogglePoolEdenRewards(t *testing.T) {

	AccAddress := sample.AccAddress()

	sender, _ := sdk.AccAddressFromBech32(AccAddress)

	msg := MsgTogglePoolEdenRewards{
		Authority: AccAddress,
	}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	require.Equal(t, msg.Route(), "masterchef")
	require.Equal(t, msg.Type(), "toggle_pool_eden_rewards")
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sender})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(bz))
}

func TestMsgTogglePoolEdenRewards_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTogglePoolEdenRewards
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgTogglePoolEdenRewards{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgTogglePoolEdenRewards{
				Authority: sample.AccAddress(),
			},
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
