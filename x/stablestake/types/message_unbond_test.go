package types

import (
	"errors"
	"testing"

	"cosmossdk.io/math"

	"github.com/elys-network/elys/v4/testutil/sample"
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
				Amount:  math.NewInt(100),
			},
			err: errors.New("invalid creator address"),
		}, {
			name: "valid address",
			msg: MsgUnbond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(100),
			},
		},
		{
			name: "negative amount",
			msg: MsgUnbond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(-100),
			},
			err: errors.New("amount should be positive"),
		},
		{
			name: "nil amount",
			msg: MsgUnbond{
				Creator: sample.AccAddress(),
				Amount:  math.Int{},
			},
			err: errors.New("amount cannot be nil"),
		},
		{
			name: "zero amount",
			msg: MsgUnbond{
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

func TestNewMsgUnbond(t *testing.T) {

	accAdress := sample.AccAddress()
	amount := math.NewInt(200)

	got := NewMsgUnbond(
		accAdress,
		amount,
		UsdcPoolId,
	)

	want := &MsgUnbond{
		Creator: accAdress,
		Amount:  amount,
		PoolId:  UsdcPoolId,
	}

	assert.Equal(t, want, got)
}
