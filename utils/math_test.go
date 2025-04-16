package utils

import (
	"testing"

	sdkmath "cosmossdk.io/math"
)

func TestAbsDifferenceWithSign(t *testing.T) {
	tests := []struct {
		a        sdkmath.LegacyDec
		b        sdkmath.LegacyDec
		expected sdkmath.LegacyDec
		sign     bool
	}{
		{sdkmath.LegacyNewDec(5), sdkmath.LegacyNewDec(3), sdkmath.LegacyNewDec(2), false},
		{sdkmath.LegacyNewDec(3), sdkmath.LegacyNewDec(5), sdkmath.LegacyNewDec(2), true},
		{sdkmath.LegacyNewDec(0), sdkmath.LegacyNewDec(0), sdkmath.LegacyNewDec(0), false},
	}

	for _, tt := range tests {
		result, sign := AbsDifferenceWithSign(tt.a, tt.b)
		if !result.Equal(tt.expected) || sign != tt.sign {
			t.Errorf("AbsDifferenceWithSign(%s, %s) = (%s, %v); want (%s, %v)", tt.a, tt.b, result, sign, tt.expected, tt.sign)
		}
	}
}
