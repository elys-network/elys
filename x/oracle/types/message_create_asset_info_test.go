package types_test

import (
	"github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateAssetInfo_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgCreateAssetInfo
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgCreateAssetInfo{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgCreateAssetInfo{
				Creator:    sample.AccAddress(),
				Denom:      ptypes.ATOM,
				Display:    "ATOM",
				ElysTicker: "ATOM",
				BandTicker: "ATOM",
				Decimal:    6,
			},
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
