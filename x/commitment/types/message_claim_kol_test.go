package types_test

import (
	"testing"

	"github.com/elys-network/elys/v4/x/commitment/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v4/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgClaimKol_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgClaimKol
		err  error
	}{
		{
			name: "invalid address",
			msg:  *types.NewMsgClaimKol("invalid_address", false),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg:  *types.NewMsgClaimKol(sample.AccAddress(), false),
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
