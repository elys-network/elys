package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdatetakeProfitPrice(t *testing.T) {
	msg := types.NewMsgUpdateTakeProfitPrice(sample.AccAddress(), 1, math.LegacyOneDec())
	//require.Equal(t, msg.GetCreator(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)})
	//msg.Creator = ""
	//require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetCreator() })
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
			name: "take profit price is 0",
			setter: func() {
				msg.Creator = sample.AccAddress()
				msg.Price = math.LegacyZeroDec()
			},
			errMsg: "",
		},
		{
			name: "take profit price is < 0",
			setter: func() {
				msg.Creator = sample.AccAddress()
				msg.Price = math.LegacyOneDec().MulInt64(-1)
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
