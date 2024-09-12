package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"testing"

	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgClose(t *testing.T) {
	msg := types.NewMsgClose(sample.AccAddress(), 1, sdk.OneInt())
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgClose)
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
			name: "lp is < 0",
			setter: func() {
				msg.Creator = sample.AccAddress()
				msg.LpAmount = sdk.OneInt().MulRaw(-1)
			},
			errMsg: "invalid lp amount: cannot be negative",
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

func TestMsgOpen(t *testing.T) {
	msg := types.NewMsgOpen(sample.AccAddress(), "uusdc", sdk.OneInt(), 1, sdk.OneDec().MulInt64(2), sdk.OneDec())
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgOpen)
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
			name: "leverage is 0",
			setter: func() {
				msg.Creator = sample.AccAddress()
				msg.Leverage = sdk.ZeroDec()
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
		},
		{
			name: "leverage is < 0",
			setter: func() {
				msg.Leverage = sdk.OneDec().MulInt64(-1)
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
		},
		{
			name: "collateral amt is 0",
			setter: func() {
				msg.Leverage = sdk.OneDec().MulInt64(2)
				msg.CollateralAsset = "uusdc"
				msg.CollateralAmount = sdk.ZeroInt()
			},
			errMsg: types.ErrInvalidCollateralAsset.Error(),
		},
		{
			name: "collateral amt is 0",
			setter: func() {
				msg.CollateralAmount = sdk.OneInt().MulRaw(10)
				msg.StopLossPrice = sdk.OneDec().MulInt64(-1)
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

func TestMsgUpdateParams(t *testing.T) {
	params := types.DefaultParams()
	msg := types.NewMsgUpdateParams(sample.AccAddress(), &params)
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgUpdateParams)
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Authority)})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)))
	msg.Authority = ""
	require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "success",
			setter: func() {
				msg.Authority = sample.AccAddress()
			},
			errMsg: "",
		},
		{
			name: "invalid address",
			setter: func() {
				msg.Authority = "invalid_address"
			},
			errMsg: "invalid creator address",
		},
		{
			name: "invalid params",
			setter: func() {
				msg.Authority = sample.AccAddress()
				msg.Params.LeverageMax = sdk.OneDec().MulInt64(100)
			},
			errMsg: "invalid params",
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

func TestMsgWhitelistAddress(t *testing.T) {
	msg := types.NewMsgWhitelist(sample.AccAddress(), sample.AccAddress())
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgWhitelist)
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Authority)})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)))
	msg.Authority = ""
	require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "success",
			setter: func() {
				msg.Authority = sample.AccAddress()
			},
			errMsg: "",
		},
		{
			name: "invalid authority address",
			setter: func() {
				msg.Authority = "invalid_address"
			},
			errMsg: "invalid creator address",
		},
		{
			name: "invalid whitelist address",
			setter: func() {
				msg.Authority = sample.AccAddress()
				msg.WhitelistedAddress = "invalid_address"
			},
			errMsg: "invalid whitelist address",
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

func TestMsgAddPool(t *testing.T) {
	addPool := types.AddPool{
		AmmPoolId:   1,
		Enabled:     false,
		Closed:      false,
		LeverageMax: sdk.OneDec().MulInt64(2),
	}
	msg := types.NewMsgAddPool(sample.AccAddress(), addPool)
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgAddPool)
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Authority)})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)))
	msg.Authority = ""
	require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "success",
			setter: func() {
				msg.Authority = sample.AccAddress()
			},
			errMsg: "",
		},
		{
			name: "invalid address",
			setter: func() {
				msg.Authority = "invalid_address"
			},
			errMsg: "invalid creator address",
		},
		{
			name: "leverage is 0",
			setter: func() {
				msg.Authority = sample.AccAddress()
				msg.Pool.LeverageMax = sdk.ZeroDec()
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
		},
		{
			name: "leverage is < 0",
			setter: func() {
				msg.Pool.LeverageMax = sdk.OneDec().MulInt64(-1)
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
		},
		{
			name: "leverage is 1",
			setter: func() {
				msg.Pool.LeverageMax = sdk.OneDec()
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
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

func TestMsgRemovePool(t *testing.T) {
	msg := types.NewMsgRemovePool(sample.AccAddress(), 1)
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgRemovePool)
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Authority)})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)))
	msg.Authority = ""
	require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "success",
			setter: func() {
				msg.Authority = sample.AccAddress()
			},
			errMsg: "",
		},
		{
			name: "invalid address",
			setter: func() {
				msg.Authority = "invalid_address"
			},
			errMsg: "invalid creator address",
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

func TestMsgUpdatePool(t *testing.T) {
	updatePool := types.UpdatePool{
		PoolId:  1,
		Enabled: false,
		Closed:  false,
	}
	msg := types.NewMsgUpdatePool(sample.AccAddress(), updatePool)
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgUpdatePool)
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Authority)})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)))
	msg.Authority = ""
	require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "success",
			setter: func() {
				msg.Authority = sample.AccAddress()
			},
			errMsg: "",
		},
		{
			name: "invalid address",
			setter: func() {
				msg.Authority = "invalid_address"
			},
			errMsg: "invalid creator address",
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

func TestMsgDewhitelistAddress(t *testing.T) {
	msg := types.NewMsgDewhitelist(sample.AccAddress(), sample.AccAddress())
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgDewhitelist)
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Authority)})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)))
	msg.Authority = ""
	require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "success",
			setter: func() {
				msg.Authority = sample.AccAddress()
			},
			errMsg: "",
		},
		{
			name: "invalid authority address",
			setter: func() {
				msg.Authority = "invalid_address"
			},
			errMsg: "invalid creator address",
		},
		{
			name: "invalid whitelist address",
			setter: func() {
				msg.Authority = sample.AccAddress()
				msg.WhitelistedAddress = "invalid_address"
			},
			errMsg: "invalid whitelisted address",
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

func TestMsgClaimRewards(t *testing.T) {
	msg := types.NewMsgClaimRewards(sample.AccAddress(), []uint64{1})
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), types.TypeMsgClaimRewards)
	require.Equal(t, msg.GetSigners(), []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)})
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)))
	msg.Sender = ""
	require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "invalid sender address",
			setter: func() {
				msg.Sender = "invalid_address"
			},
			errMsg: "invalid sender address",
		},
		{
			name: "empty pool ids",
			setter: func() {
				msg.Sender = sample.AccAddress()
				msg.Ids = []uint64{}
			},
			errMsg: "empty ids",
		},
		{
			name: "repeated pool ids",
			setter: func() {
				msg.Ids = []uint64{1, 2, 1}
			},
			errMsg: "duplicate pool id",
		},
		{
			name: "success",
			setter: func() {
				msg.Ids = []uint64{1, 2}
			},
			errMsg: "",
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
