package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateWasmConfig_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateWasmConfig
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateWasmConfig{
				Creator:                 "invalid_address",
				WasmMaxLabelSize:        "1",
				WasmMaxSize:             "1",
				WasmMaxProposalWasmSize: "1",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateWasmConfig{
				Creator:                 sample.AccAddress(),
				WasmMaxLabelSize:        "1",
				WasmMaxSize:             "1",
				WasmMaxProposalWasmSize: "1",
			},
		},
		{
			name: "invalid WasmMaxLabelSize",
			msg: MsgUpdateWasmConfig{
				Creator:                 sample.AccAddress(),
				WasmMaxLabelSize:        "-1",
				WasmMaxSize:             "1",
				WasmMaxProposalWasmSize: "1",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid WasmMaxSize",
			msg: MsgUpdateWasmConfig{
				Creator:                 sample.AccAddress(),
				WasmMaxLabelSize:        "1",
				WasmMaxSize:             "-1",
				WasmMaxProposalWasmSize: "1",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid WasmMaxProposalWasmSize",
			msg: MsgUpdateWasmConfig{
				Creator:                 sample.AccAddress(),
				WasmMaxLabelSize:        "1",
				WasmMaxSize:             "1",
				WasmMaxProposalWasmSize: "-1",
			},
			err: sdkerrors.ErrInvalidRequest,
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
