package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/v6/x/oracle/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgFeedMultiplePrices_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgFeedMultiplePrices
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgFeedMultiplePrices{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgFeedMultiplePrices{
				Creator: sample.AccAddress(),
				FeedPrices: []types.FeedPrice{
					{
						Asset:  ptypes.ATOM,
						Price:  math.LegacyOneDec(),
						Source: "source",
					},
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
