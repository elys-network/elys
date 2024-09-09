package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/assert"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func TestNewMsgUpdateParams(t *testing.T) {

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	params := &Params{}
	got := NewMsgUpdateParams(authority, params)

	want := &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}

	assert.Equal(t, want, got)
}

func TestMsgUpdateParams_Route(t *testing.T) {
	msg := MsgUpdateParams{}
	assert.Equal(t, "perpetual", msg.Route())
}

func TestMsgUpdateParams_Type(t *testing.T) {
	msg := MsgUpdateParams{}
	assert.Equal(t, "update_params", msg.Type())
}

func TestMsgUpdateParams_GetSigners(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgUpdateParams{Authority: accAdress}

	creator, err := sdk.AccAddressFromBech32(accAdress)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestMsgUpdateParams_GetSignBytes(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgUpdateParams{Authority: accAdress}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	b := sdk.MustSortJSON(bz)

	assert.Equal(t, b, msg.GetSignBytes())
}

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {

	type Test struct {
		title string
		msg   MsgUpdateParams
		want  error
	}

	tableTest := []Test{
		{
			title: "invalid address",
			msg: MsgUpdateParams{
				Authority: "invalid address",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "successful",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
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
