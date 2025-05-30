package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultParams(t *testing.T) {
	require.Equal(t, types.DefaultParams(), types.NewParams())
}

func TestValidateLeverageMax(t *testing.T) {
	params := types.NewParams()
	tests := []struct {
		name   string
		setter func()
		err    string
	}{
		{
			name: "success",
			setter: func() {
			},
			err: "",
		},
		{
			name: "LeverageMax is nil",
			setter: func() {
				params.LeverageMax = sdkmath.LegacyDec{}
			},
			err: "leverage max must be not nil",
		},
		{
			name: "LeverageMax < 1",
			setter: func() {
				params.LeverageMax = sdkmath.LegacyMustNewDecFromStr("0.5")
			},
			err: "leverage max must be greater than 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := params.Validate()
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateEpochLength(t *testing.T) {
	params := types.NewParams()
	tests := []struct {
		name   string
		setter func()
		err    string
	}{
		{
			name: "success",
			setter: func() {
			},
			err: "",
		},
		{
			name: "EpochLength is -va",
			setter: func() {
				params.EpochLength = -1
			},
			err: "epoch length should be positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := params.Validate()
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateSafetyFactor(t *testing.T) {
	params := types.NewParams()
	tests := []struct {
		name   string
		setter func()
		err    string
	}{
		{
			name: "success",
			setter: func() {
			},
			err: "",
		},
		{
			name: "SafetyFactor is nil",
			setter: func() {
				params.SafetyFactor = sdkmath.LegacyDec{}
			},
			err: "safety factor must be not nil",
		},
		{
			name: "SafetyFactor is 0 ",
			setter: func() {
				params.SafetyFactor = sdkmath.LegacyZeroDec()
			},
			err: "safety factor must be positive",
		}, {
			name: "SafetyFactor is < 0 ",
			setter: func() {
				params.SafetyFactor = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			err: "safety factor must be positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := params.Validate()
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidatePoolThreshold(t *testing.T) {
	params := types.NewParams()
	tests := []struct {
		name   string
		setter func()
		err    string
	}{
		{
			name: "success",
			setter: func() {
			},
			err: "",
		},
		{
			name: "PoolOpenThreshold is nil",
			setter: func() {
				params.PoolOpenThreshold = sdkmath.LegacyDec{}
			},
			err: "pool open threshold must be not nil",
		},
		{
			name: "PoolOpenThreshold is 0 ",
			setter: func() {
				params.PoolOpenThreshold = sdkmath.LegacyZeroDec()
			},
			err: "pool open threshold must be positive",
		}, {
			name: "PoolOpenThreshold is < 0 ",
			setter: func() {
				params.PoolOpenThreshold = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			err: "pool open threshold must be positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := params.Validate()
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateNumberOfBlocks(t *testing.T) {
	params := types.NewParams()
	tests := []struct {
		name   string
		setter func()
		err    string
	}{
		{
			name: "success",
			setter: func() {
			},
			err: "",
		},
		{
			name: "NumberPerBlock is 0",
			setter: func() {
				params.NumberPerBlock = 0
			},
			err: "",
		},
		{
			name: "NumberPerBlock is nil",
			setter: func() {
				params.NumberPerBlock = -1
			},
			err: "number of positions per block must be positive",
		},
		{
			name: "PoolOpenThreshold is more than page limit ",
			setter: func() {
				params.NumberPerBlock = types.MaxPageLimit + 1
			},
			err: "number of positions per block should not exceed page limit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := params.Validate()
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDuplicate(t *testing.T) {
	params := types.NewParams()
	tests := []struct {
		name   string
		setter func()
		err    string
	}{
		{
			name: "success",
			setter: func() {
				params.EnabledPools = []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			},
			err: "",
		},
		{
			name: "Duplicate",
			setter: func() {
				params.EnabledPools = []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1}
			},
			err: "array must not contain duplicate values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := params.Validate()
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
