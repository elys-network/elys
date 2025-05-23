package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdatePool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     MsgUpdatePool
		wantErr bool
	}{
		{
			name: "valid message",
			msg: MsgUpdatePool{
				Authority:            sample.AccAddress(),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
				PoolId:               1,
			},
			wantErr: false,
		},
		{
			name: "invalid authority address",
			msg: MsgUpdatePool{
				Authority:            "invalid_address",
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
				PoolId:               1,
			},
			wantErr: true,
		},
		{
			name: "InterestRateMax less than InterestRateMin",
			msg: MsgUpdatePool{
				Authority:            sample.AccAddress(),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
				PoolId:               1,
			},
			wantErr: true,
		},
		{
			name: "InterestRateMax is nil",
			msg: MsgUpdatePool{
				Authority:            sample.AccAddress(),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
				PoolId:               1,
			},
			wantErr: true,
		},
		{
			name: "InterestRateMin is nil",
			msg: MsgUpdatePool{
				Authority:            sample.AccAddress(),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
				PoolId:               1,
			},
			wantErr: true,
		},
		{
			name: "InterestRateIncrease is nil",
			msg: MsgUpdatePool{
				Authority:            sample.AccAddress(),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
				PoolId:               1,
			},
			wantErr: true,
		},
		{
			name: "InterestRateDecrease is nil",
			msg: MsgUpdatePool{
				Authority:            sample.AccAddress(),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
				PoolId:               1,
			},
			wantErr: true,
		},
		{
			name: "HealthGainFactor is nil",
			msg: MsgUpdatePool{
				Authority:            sample.AccAddress(),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
				PoolId:               1,
			},
			wantErr: true,
		},
		{
			name: "MaxLeverageRatio is nil",
			msg: MsgUpdatePool{
				Authority:            sample.AccAddress(),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				PoolId:               1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgUpdatePool(t *testing.T) {

	accAdress := sample.AccAddress()

	got := NewMsgUpdatePool(
		accAdress,
		math.LegacyMustNewDecFromStr("0.07"),
		math.LegacyMustNewDecFromStr("0.05"),
		math.LegacyMustNewDecFromStr("0.01"),
		math.LegacyMustNewDecFromStr("0.02"),
		math.LegacyMustNewDecFromStr("0.01"),
		math.LegacyMustNewDecFromStr("0.1"),
		1,
	)

	want := &MsgUpdatePool{
		Authority:            accAdress,
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.07"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.05"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.02"),
		HealthGainFactor:     math.LegacyMustNewDecFromStr("0.01"),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.1"),
		PoolId:               1,
	}

	assert.Equal(t, want, got)
}
