package types_test

import (
	"testing"

	"github.com/elys-network/elys/v7/x/commitment/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgClaimAirdrop_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgClaimAirdrop
		err  error
	}{
		{
			name: "invalid address",
			msg:  *types.NewMsgClaimAirdrop("invalid_address"),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg:  *types.NewMsgClaimAirdrop(sample.AccAddress()),
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
