package types_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestDefaultParams(t *testing.T) {
	require.Equal(t, types.DefaultParams(), types.NewParams(math.NewInt(10_000_000), 86400*7, false))
	output, err := yaml.Marshal(types.DefaultParams())
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams().String(), string(output))
}

func TestParamsValidation(t *testing.T) {
	params := types.DefaultParams()
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
			name: "PoolCreationFee is nil",
			setter: func() {
				params.PoolCreationFee = sdk.Int{}
			},
			err: "pool creation fee must not be empty",
		},
		{
			name: "PoolCreationFee < 0",
			setter: func() {
				params.PoolCreationFee = sdk.NewInt(1).MulRaw(-1)
			},
			err: "pool creation fee must be positive",
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
