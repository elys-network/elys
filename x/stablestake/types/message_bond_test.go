package types

import (
	"errors"
	"testing"

	"cosmossdk.io/math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgBond_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgBond
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgBond{
				Creator: "invalid_address",
				Amount:  math.NewInt(100),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgBond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(100),
			},
		},
		{
			name: "negative amount",
			msg: MsgBond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(-100),
			},
			err: errors.New("amount should be positive"),
		},
		{
			name: "nil amount",
			msg: MsgBond{
				Creator: sample.AccAddress(),
				Amount:  math.Int{},
			},
			err: errors.New("amount cannot be nil"),
		},
		{
			name: "zero amount",
			msg: MsgBond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(0),
			},
			err: errors.New("amount should be positive"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewMsgBond(t *testing.T) {

	accAdress := sample.AccAddress()
	amount := math.NewInt(200)

	got := NewMsgBond(
		accAdress,
		amount,
		UsdcPoolId,
	)

	want := &MsgBond{
		Creator: accAdress,
		Amount:  amount,
		PoolId:  UsdcPoolId,
	}

	assert.Equal(t, want, got)
}
