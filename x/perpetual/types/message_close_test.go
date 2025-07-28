package types

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/assert"
)

func TestNewMsgClose(t *testing.T) {

	accAdress := sample.AccAddress()
	got := NewMsgClose(
		accAdress,
		2,
		math.NewInt(20),
		1,
		math.LegacyOneDec(),
	)

	want := &MsgClose{
		Creator:      accAdress,
		Id:           2,
		Amount:       math.NewInt(20),
		PoolId:       1,
		ClosingRatio: math.LegacyOneDec(),
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
			title: "invalid closing ratio",
			msg: MsgClose{
				Creator:      sample.AccAddress(),
				Amount:       math.NewInt(20),
				ClosingRatio: math.LegacyOneDec().MulInt64(-1),
			},
			want: errors.New("ClosingRatio is nil"),
		},
		{
			title: "invalid closing ratio",
			msg: MsgClose{
				Creator:      sample.AccAddress(),
				Amount:       math.NewInt(0),
				ClosingRatio: math.LegacyZeroDec(),
			},
			want: errors.New("closing ratio and amount both cannot be zero"),
		},
		{
			title: "invalid pool id",
			msg: MsgClose{
				Creator:      sample.AccAddress(),
				Amount:       math.NewInt(20),
				PoolId:       0,
				ClosingRatio: math.LegacyOneDec(),
			},
			want: errors.New("invalid pool id"),
		},
		{
			title: "successful",
			msg: MsgClose{
				Creator:      sample.AccAddress(),
				Amount:       math.NewInt(20),
				PoolId:       1,
				ClosingRatio: math.LegacyOneDec(),
			},
			want: nil,
		},
	}

	for _, test := range tableTest {
		t.Run(test.title, func(t *testing.T) {
			got := test.msg.ValidateBasic()

			if got != nil {
				assert.Contains(t, got.Error(), test.want.Error())
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}

}
