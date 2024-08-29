package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateStopLoss(t *testing.T) {
	msg := types.NewMsgUpdateStopLoss(sample.AccAddress(), 1, sdk.OneDec())
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgUpdateStopLoss)
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)))
	msg.Creator = ""
	require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })
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
				msg.Price = sdk.ZeroDec()
			},
			errMsg: "",
		},
		{
			name: "stop loss price is < 0",
			setter: func() {
				msg.Creator = sample.AccAddress()
				msg.Price = sdk.OneDec().MulInt64(-1)
			},
			errMsg: "stop loss price cannot be negative",
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
