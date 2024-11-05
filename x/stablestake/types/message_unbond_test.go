package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgUnbond_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUnbond
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUnbond{
				Creator: "invalid_address",
				Amount:  sdk.NewInt(100),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUnbond{
				Creator: sample.AccAddress(),
				Amount:  sdk.NewInt(100),
			},
		},
		{
			name: "negative amount",
			msg: MsgUnbond{
				Creator: sample.AccAddress(),
				Amount:  sdk.NewInt(-100),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "zero amount",
			msg: MsgUnbond{
				Creator: sample.AccAddress(),
				Amount:  sdk.NewInt(0),
			},
			err: sdkerrors.ErrInvalidRequest,
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

func TestNewMsgUnbond(t *testing.T) {

	accAdress := sample.AccAddress()
	amount := math.NewInt(200)

	got := NewMsgUnbond(
		accAdress,
		amount,
	)

	want := &MsgUnbond{
		Creator: accAdress,
		Amount:  amount,
	}

	assert.Equal(t, want, got)
}

func TestMsgUnbond_Route(t *testing.T) {
	msg := MsgUnbond{}
	assert.Equal(t, "stablestake", msg.Route())
}

func TestMsgUnbond_Type(t *testing.T) {
	msg := MsgUnbond{}
	assert.Equal(t, "unbond", msg.Type())
}

func TestMsgUnbond_GetSigners(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgUnbond{Creator: accAdress}

	creator, err := sdk.AccAddressFromBech32(accAdress)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestMsgUnbond_GetSignBytes(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgUnbond{Creator: accAdress}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	b := sdk.MustSortJSON(bz)

	assert.Equal(t, b, msg.GetSignBytes())
}
