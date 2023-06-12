package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func TestAbsDifferenceWithSign(t *testing.T) {
	tests := []struct {
		a        sdk.Dec
		b        sdk.Dec
		expected sdk.Dec
		sign     bool
	}{
		{sdk.NewDec(5), sdk.NewDec(3), sdk.NewDec(2), false},
		{sdk.NewDec(3), sdk.NewDec(5), sdk.NewDec(2), true},
		{sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0), false},
	}

	for _, tt := range tests {
		result, sign := types.AbsDifferenceWithSign(tt.a, tt.b)
		if !result.Equal(tt.expected) || sign != tt.sign {
			t.Errorf("AbsDifferenceWithSign(%s, %s) = (%s, %v); want (%s, %v)", tt.a, tt.b, result, sign, tt.expected, tt.sign)
		}
	}
}
