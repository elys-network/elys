package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/v5/x/oracle/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgFeedPrice_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgFeedPrice
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgFeedPrice{
				Provider: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgFeedPrice{
				Provider: sample.AccAddress(),
				FeedPrice: types.FeedPrice{
					Price:  sdkmath.LegacyMustNewDecFromStr("100"),
					Asset:  ptypes.ATOM,
					Source: "source",
				},
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
