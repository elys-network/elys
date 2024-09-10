package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/x/clock/types"
)

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name     string
		params   types.Params
		expError bool
	}{
		{"default", types.DefaultParams(), false},
		{
			"valid: no contracts, enough gas",
			types.NewParams([]string(nil), 100_000),
			false,
		},
		{
			"invalid: address malformed",
			types.NewParams([]string{"invalid address"}, 100_000),
			true,
		},
		{
			"invalid: not enough gas",
			types.NewParams([]string(nil), 1),
			true,
		},
		{
			"duplicate addresses",
			types.NewParams([]string{"cosmos10duudma7ef9849ee42zhe5q4t4fmk0z99uuh92", "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5", "cosmos10duudma7ef9849ee42zhe5q4t4fmk0z99uuh92"}, 100_001),
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}
