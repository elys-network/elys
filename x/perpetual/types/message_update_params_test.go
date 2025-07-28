package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/testutil/sample"
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

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {

	params := DefaultParams()
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
				Params:    &params,
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
