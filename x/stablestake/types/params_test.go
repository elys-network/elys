package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	params := types.DefaultParams()
	tests := []struct {
		name      string
		err       string
		runBefore func()
	}{
		{
			name: "interest rate max is nil",
			err:  "InterestRateMax",
			runBefore: func() {
				params.LegacyInterestRateMax = math.LegacyDec{}
			},
		},
		{
			name: "total value is nil",
			err:  "TotalValue is nil",
			runBefore: func() {
				params = types.DefaultParams()
				params.TotalValue = math.Int{}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.runBefore()
			err := params.Validate()
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
