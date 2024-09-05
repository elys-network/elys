package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgAddEntry_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddEntry
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddEntry{
				Creator:   "invalid_address",
				Decimals:  6,
				BaseDenom: "uusdc",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  6,
				BaseDenom: "uusdc",
			},
		},
		{
			name: "invalid decimal",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  4,
				BaseDenom: "uusdc",
			},
			err: ErrDecimalsInvalid,
		},
		{
			name: "invalid decimal",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  19,
				BaseDenom: "uusdc",
			},
			err: ErrDecimalsInvalid,
		},
		{
			name: "valid decimal",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  6,
				BaseDenom: "uusdc",
			},
		},
		{
			name: "valid decimal",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  12,
				BaseDenom: "uusdc",
			},
		},
		{
			name: "valid decimal",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  18,
				BaseDenom: "uusdc",
			},
		},
		{
			name: "invalid base denom",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  18,
				BaseDenom: "",
			},
			err: ErrInvalidBaseDenom,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
