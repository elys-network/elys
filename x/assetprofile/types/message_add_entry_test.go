package types

import (
	"errors"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
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
				Denom:     "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
			},
		},
		{
			name: "invalid decimal",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  3,
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
				Denom:     "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
			},
		},
		{
			name: "valid decimal",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  12,
				BaseDenom: "uusdc",
				Denom:     "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
			},
		},
		{
			name: "valid decimal",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  18,
				BaseDenom: "uusdc",
				Denom:     "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
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
		{
			name: "invalid base denom - single char",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  18,
				BaseDenom: "a",
			},
			err: ErrInvalidBaseDenom,
		},
		{
			name: "invalid base denom - too long",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  18,
				BaseDenom: "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A652180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
			},
			err: ErrInvalidBaseDenom,
		},
		{
			name: "invalid denom",
			msg: MsgAddEntry{
				Creator:   sample.AccAddress(),
				Decimals:  18,
				BaseDenom: "uusdc",
				Denom:     "",
			},
			err: errors.New("invalid denom"),
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
