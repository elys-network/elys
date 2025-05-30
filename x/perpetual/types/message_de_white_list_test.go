package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/assert"
)

func TestNewMsgDewhitelist(t *testing.T) {

	accAdress := sample.AccAddress()
	whileListAddress := sample.AccAddress()

	got := NewMsgDewhitelist(
		accAdress,
		whileListAddress,
	)

	want := &MsgDewhitelist{
		Authority:          accAdress,
		WhitelistedAddress: whileListAddress,
	}

	assert.Equal(t, want, got)
}

func TestMsgDewhitelist_ValidateBasic(t *testing.T) {

	type Test struct {
		title string
		msg   MsgDewhitelist
		want  error
	}

	tableTest := []Test{
		{
			title: "authority - invalid address",
			msg: MsgDewhitelist{
				Authority: "invalid",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "while list - invalid address",
			msg: MsgDewhitelist{
				Authority:          sample.AccAddress(),
				WhitelistedAddress: "invalid",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "while list - invalid address",
			msg: MsgDewhitelist{
				Authority:          sample.AccAddress(),
				WhitelistedAddress: sample.AccAddress(),
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
