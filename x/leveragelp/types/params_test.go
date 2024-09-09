package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestDefaultParams(t *testing.T) {
	require.Equal(t, types.DefaultParams(), types.NewParams())
	output, err := yaml.Marshal(types.DefaultParams())
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams().String(), string(output))
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
				params.LeverageMax = sdk.Dec{}
			},
			err: "leverage max must be not nil",
		},
		{
			name: "LeverageMax < 1",
			setter: func() {
				params.LeverageMax = sdk.MustNewDecFromStr("0.5")
			},
			err: "leverage max must be greater than 1",
		}, {
			name: "LeverageMax is 100",
			setter: func() {
				params.LeverageMax = sdk.OneDec().MulInt64(100)
			},
			err: "",
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
				params.SafetyFactor = sdk.Dec{}
			},
			err: "safety factor must be not nil",
		},
		{
			name: "SafetyFactor is 0 ",
			setter: func() {
				params.SafetyFactor = sdk.ZeroDec()
			},
			err: "safety factor must be positive",
		}, {
			name: "SafetyFactor is < 0 ",
			setter: func() {
				params.SafetyFactor = sdk.OneDec().MulInt64(-1)
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
				params.PoolOpenThreshold = sdk.Dec{}
			},
			err: "pool open threshold must be not nil",
		},
		{
			name: "PoolOpenThreshold is 0 ",
			setter: func() {
				params.PoolOpenThreshold = sdk.ZeroDec()
			},
			err: "pool open threshold must be positive",
		}, {
			name: "PoolOpenThreshold is < 0 ",
			setter: func() {
				params.PoolOpenThreshold = sdk.OneDec().MulInt64(-1)
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
