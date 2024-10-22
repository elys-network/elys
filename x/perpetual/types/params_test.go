package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func TestValidateMaxLimitOrder(t *testing.T) {

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
			name: "MaxLimitOrder is nil",
			setter: func() {
				params.MaxLimitOrder = sdk.OneInt().MulRaw(-100).Int64()
			},
			err: "MaxLimitOrder should not be -ve",
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
