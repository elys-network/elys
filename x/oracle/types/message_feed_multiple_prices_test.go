package types_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
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
				Prices: []types.Price{
					{
						Asset:       ptypes.ATOM,
						Price:       math.LegacyOneDec(),
						Source:      "source",
						Provider:    "provider",
						Timestamp:   1,
						BlockHeight: 1,
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
