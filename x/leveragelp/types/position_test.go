package types_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/v5/x/leveragelp/types"

	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestPosition(t *testing.T) {
	creator := sample.AccAddress()
	position := types.NewPosition(creator, sdk.NewCoin("uusdc", sdkmath.OneInt()), 1)
	require.Equal(t, position.GetOwnerAddress(), sdk.MustAccAddressFromBech32(creator))
	require.Equal(t, position.GetPositionAddress(), authtypes.NewModuleAddress(fmt.Sprintf("leveragelp/%d", position.Id)))
	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "success",
			setter: func() {
				position.Id = 1
			},
			errMsg: "",
		},
		{
			name: "invalid address",
			setter: func() {
				position.Address = "invalid_address"
			},
			errMsg: "decoding bech32 failed",
		},
		{
			name: "empty address",
			setter: func() {
				position.Address = ""
			},
			errMsg: "empty address string is not allowed",
		},
		{
			name: "position id is 0",
			setter: func() {
				position.Address = creator
				position.Id = 0
			},
			errMsg: "position id cannot be 0",
		},
		{
			name: "leveraged lp amount is < 0",
			setter: func() {
				position.Id = 1
				position.LeveragedLpAmount = sdkmath.OneInt().MulRaw(-1)
			},
			errMsg: "leveraged lp amount cannot be negative",
		},
		{
			name: "Collateral is invalid",
			setter: func() {
				position.LeveragedLpAmount = sdkmath.OneInt()
				position.Collateral = sdk.Coin{"$$$$", sdkmath.OneInt()}
			},
			errMsg: "invalid collateral coin",
		},
		{
			name: "Stop loss is invalid",
			setter: func() {
				position.Collateral = sdk.Coin{"uusdc", sdkmath.OneInt()}
				position.StopLossPrice = sdkmath.LegacyOneDec().MulInt64(-10)
			},
			errMsg: "stop loss price cannot be negative",
		},
		{
			name: "liabilities is invalid",
			setter: func() {
				position.StopLossPrice = sdkmath.LegacyOneDec().MulInt64(10)
				position.Liabilities = sdkmath.OneInt().MulRaw(-1)
			},
			errMsg: "liabilities cannot be negative",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := position.Validate()
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
