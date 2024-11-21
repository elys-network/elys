package types

import (
	"cosmossdk.io/math"
	"fmt"
	"testing"

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
				Amount:  math.NewInt(100),
			},
			err: fmt.Errorf("invalid creator address"),
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
			err: fmt.Errorf("amount should be positive"),
		},
		{
			name: "nil amount",
			msg: MsgUnbond{
				Creator: sample.AccAddress(),
				Amount:  math.Int{},
			},
			err: fmt.Errorf("amount cannot be nil"),
		},
		{
			name: "zero amount",
			msg: MsgUnbond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(0),
			},
			err: fmt.Errorf("amount should be positive"),
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
	)

	want := &MsgUnbond{
		Creator: accAdress,
		Amount:  amount,
	}

	assert.Equal(t, want, got)
}
