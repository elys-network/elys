package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
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

func TestMsgDewhitelist_Route(t *testing.T) {
	msg := MsgDewhitelist{}
	assert.Equal(t, "perpetual", msg.Route())
}

func TestMsgDewhitelist_Type(t *testing.T) {
	msg := MsgDewhitelist{}
	assert.Equal(t, "dewhitelist", msg.Type())
}

func TestMsgDewhitelist_GetSigners(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgDewhitelist{Authority: accAdress}

	creator, err := sdk.AccAddressFromBech32(accAdress)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestMsgDewhitelist_GetSignBytes(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgDewhitelist{Authority: accAdress}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	b := sdk.MustSortJSON(bz)

	assert.Equal(t, b, msg.GetSignBytes())
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
