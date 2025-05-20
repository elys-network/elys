package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/elys-network/elys/v4/testutil/sample"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateStopLoss(t *testing.T) {
	msg := types.NewMsgUpdateStopLoss(sample.AccAddress(), 1, sdkmath.LegacyOneDec())
	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "success",
			setter: func() {
				msg.Creator = sample.AccAddress()
			},
			errMsg: "",
		},
		{
			name: "invalid address",
			setter: func() {
				msg.Creator = "invalid_address"
			},
			errMsg: "invalid creator address",
		},
		{
			name: "stop loss price is 0",
			setter: func() {
				msg.Creator = sample.AccAddress()
				msg.Price = sdkmath.LegacyZeroDec()
			},
			errMsg: "",
		},
		{
			name: "stop loss price is < 0",
			setter: func() {
				msg.Creator = sample.AccAddress()
				msg.Price = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			errMsg: "price is negative",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := msg.ValidateBasic()
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
