package types

import (
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/stretchr/testify/assert"
)

func TestNewMsgClose(t *testing.T) {

	accAdress := sample.AccAddress()
	got := NewMsgClose(
		accAdress,
		2,
		math.NewInt(20),
	)

	want := &MsgClose{
		Creator: accAdress,
		Id:      2,
		Amount:  math.NewInt(20),
	}

	assert.Equal(t, want, got)
}

func TestMsgClose_ValidateBasic(t *testing.T) {

	type Test struct {
		title string
		msg   MsgClose
		want  error
	}

	tableTest := []Test{
		{
			title: "invalid address",
			msg: MsgClose{
				Creator: "invalid",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "invalid amount",
			msg: MsgClose{
				Creator: sample.AccAddress(),
			},
			want: ErrInvalidAmount,
		},
		{
			title: "invalid amount - negative",
			msg: MsgClose{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(-20),
			},
			want: ErrInvalidAmount,
		},
		{
			title: "successful",
			msg: MsgClose{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(20),
			},
			want: nil,
		},
	}

	for _, test := range tableTest {
		t.Run(test.title, func(t *testing.T) {
			got := test.msg.ValidateBasic()

			if got != nil {
				assert.ErrorIs(t, got, test.want)
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}

}
