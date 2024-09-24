package types_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func TestValidateMinBorrowInterestAmount(t *testing.T) {

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
			name: "MinBorrowInterestAmount is nil",
			setter: func() {
				params.MinBorrowInterestAmount = sdk.Int{}
			},
			err: "MinBorrowInterestAmount must be not nil",
		},
		{
			name: "MinBorrowInterestAmount is -ve",
			setter: func() {
				params.MinBorrowInterestAmount = sdk.OneInt().MulRaw(-100)
			},
			err: "MinBorrowInterestAmount must be positive",
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
