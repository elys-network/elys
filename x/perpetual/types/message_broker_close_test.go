package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/assert"
)

func TestNewMsgBrokerClose(t *testing.T) {

	accAdress := sample.AccAddress()
	amount := math.NewInt(200)
	owner := sample.AccAddress()

	got := NewMsgBrokerClose(
		accAdress,
		2,
		amount,
		owner,
	)

	want := &MsgBrokerClose{
		Creator: accAdress,
		Id:      2,
		Amount:  amount,
		Owner:   owner,
	}

	assert.Equal(t, want, got)
}

func TestMsgBrokerClose_Route(t *testing.T) {
	msg := MsgBrokerClose{}
	assert.Equal(t, "perpetual", msg.Route())
}

func TestMsgBrokerClose_Type(t *testing.T) {
	msg := MsgBrokerClose{}
	assert.Equal(t, "broker_close", msg.Type())
}

func TestMsgBrokerClose_GetSigners(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgBrokerClose{Creator: accAdress}

	creator, err := sdk.AccAddressFromBech32(accAdress)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestMsgBrokerClose_GetSignBytes(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgBrokerClose{Creator: accAdress}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	b := sdk.MustSortJSON(bz)

	assert.Equal(t, b, msg.GetSignBytes())
}

func TestMsgBrokerClose_ValidateBasic(t *testing.T) {

	type Test struct {
		title string
		msg   MsgBrokerClose
		want  error
	}

	tableTest := []Test{
		{
			title: "invalid address",
			msg: MsgBrokerClose{
				Creator: "invalid",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "invalid owner",
			msg: MsgBrokerClose{
				Creator: sample.AccAddress(),
				Owner:   "invalid",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "invalid amount - is nil",
			msg: MsgBrokerClose{
				Creator: sample.AccAddress(),
				Owner:   sample.AccAddress(),
			},
			want: ErrInvalidAmount,
		},
		{
			title: "invalid amount - is negative",
			msg: MsgBrokerClose{
				Creator: sample.AccAddress(),
				Owner:   sample.AccAddress(),
				Amount:  math.NewInt(-200),
			},
			want: ErrInvalidAmount,
		},
		{
			title: "successful ",
			msg: MsgBrokerClose{
				Creator: sample.AccAddress(),
				Owner:   sample.AccAddress(),
				Amount:  math.NewInt(200),
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
