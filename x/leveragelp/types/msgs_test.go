package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/elys-network/elys/x/leveragelp/types"

	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgClose(t *testing.T) {
	msg := types.NewMsgClose(sample.AccAddress(), 1, sdkmath.OneInt())
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
				msg.LpAmount = sdkmath.OneInt().MulRaw(-1)
			},
			errMsg: "invalid lp amount: cannot be zero or negative",
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
	msg := types.NewMsgOpen(sample.AccAddress(), "uusdc", sdkmath.OneInt(), 1, sdkmath.LegacyOneDec().MulInt64(2), sdkmath.LegacyOneDec())
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
				msg.Leverage = sdkmath.LegacyZeroDec()
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
		},
		{
			name: "leverage is < 0",
			setter: func() {
				msg.Leverage = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
		},
		{
			name: "collateral amt is 0",
			setter: func() {
				msg.Leverage = sdkmath.LegacyOneDec().MulInt64(2)
				msg.CollateralAsset = "uusdc"
				msg.CollateralAmount = sdkmath.ZeroInt()
			},
			errMsg: types.ErrInvalidCollateralAsset.Error(),
		},
		{
			name: "collateral amt is 0",
			setter: func() {
				msg.CollateralAmount = sdkmath.OneInt().MulRaw(10)
				msg.StopLossPrice = sdkmath.LegacyOneDec().MulInt64(-1)
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
				msg.Params.LeverageMax = sdkmath.LegacyOneDec().MulInt64(100)
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
		LeverageMax: sdkmath.LegacyOneDec().MulInt64(2),
	}
	msg := types.NewMsgAddPool(sample.AccAddress(), addPool)

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
				msg.Pool.LeverageMax = sdkmath.LegacyZeroDec()
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
		},
		{
			name: "leverage is < 0",
			setter: func() {
				msg.Pool.LeverageMax = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			errMsg: types.ErrLeverageTooSmall.Error(),
		},
		{
			name: "leverage is 1",
			setter: func() {
				msg.Pool.LeverageMax = sdkmath.LegacyOneDec()
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

func TestMsgUpdatePool(t *testing.T) {
	msg := types.MsgUpdatePool{
		Authority:          "",
		PoolId:             0,
		LeverageMax:        sdkmath.LegacyDec{},
		MaxLeveragelpRatio: sdkmath.LegacyDec{},
	}

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "invalid authority address",
			setter: func() {
				msg.Authority = "invalid_address"
			},
			errMsg: "invalid authority address",
		},
		{
			name: "invalid pool id",
			setter: func() {
				msg.Authority = sample.AccAddress()
				msg.PoolId = 0
			},
			errMsg: "invalid pool id",
		},
		{
			name: "invalid LeverageMax",
			setter: func() {
				msg.PoolId = 1
				msg.LeverageMax = sdkmath.LegacyZeroDec()
			},
			errMsg: "invalid leverage max",
		},
		{
			name: "invalid MaxLeveragelpRatio",
			setter: func() {
				msg.LeverageMax = sdkmath.LegacyNewDec(2)
				msg.MaxLeveragelpRatio = sdkmath.LegacyNewDec(-1)
			},
			errMsg: "invalid max leverage ratio",
		},
		{
			name: "success",
			setter: func() {
				msg.MaxLeveragelpRatio = sdkmath.LegacyMustNewDecFromStr("0.2")
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

func TestMsgUpdateEnabledPools(t *testing.T) {
	msg := types.MsgUpdateEnabledPools{
		Authority:    "",
		EnabledPools: nil,
	}

	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "invalid authority address",
			setter: func() {
				msg.Authority = "invalid_address"
			},
			errMsg: "invalid authority address",
		},
		{
			name: "repeated pool ids",
			setter: func() {
				msg.Authority = sample.AccAddress()
				msg.EnabledPools = []uint64{1, 2, 1}
			},
			errMsg: "duplicate pool id",
		},
		{
			name: "success",
			setter: func() {
				msg.EnabledPools = []uint64{1, 2}
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
