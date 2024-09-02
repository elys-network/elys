package types_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"testing"

	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgAddCollateral(t *testing.T) {
	msg := types.NewMsgAddCollateral(sample.AccAddress(), 1, sdkmath.OneInt())
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgAddCollateral)
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
			name: "collateral is 0",
			setter: func() {
				msg.Creator = sample.AccAddress()
				msg.Collateral = sdk.ZeroInt()
			},
			errMsg: "collateral should be positive",
		},
		{
			name: "collateral is < 0",
			setter: func() {
				msg.Collateral = sdk.OneInt().MulRaw(-1)
			},
			errMsg: "collateral should be positive",
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
