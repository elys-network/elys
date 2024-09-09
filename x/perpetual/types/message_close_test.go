package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
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

func TestMsgClose_Route(t *testing.T) {
	msg := MsgClose{}
	assert.Equal(t, "perpetual", msg.Route())
}

func TestMsgClose_Type(t *testing.T) {
	msg := MsgClose{}
	assert.Equal(t, "close", msg.Type())
}

func TestMsgClose_GetSigners(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgClose{Creator: accAdress}

	creator, err := sdk.AccAddressFromBech32(accAdress)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestMsgClose_GetSignBytes(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgClose{Creator: accAdress}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	b := sdk.MustSortJSON(bz)

	assert.Equal(t, b, msg.GetSignBytes())
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
