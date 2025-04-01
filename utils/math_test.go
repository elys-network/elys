package utils_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/utils"
	"strings"
	"testing"
)

func TestMinDec(t *testing.T) {
	tests := []struct {
		name     string
		d1       math.Dec
		d2       math.Dec
		expected math.Dec
	}{
		// Positive-only cases
		{"d1 < d2", math.NewDecFromInt64(1), math.NewDecFromInt64(2), math.NewDecFromInt64(1)},
		{"d1 > d2", math.NewDecFromInt64(3), math.NewDecFromInt64(2), math.NewDecFromInt64(2)},
		{"d1 == d2", math.NewDecFromInt64(5), math.NewDecFromInt64(5), math.NewDecFromInt64(5)},

		// Negative-only cases
		{"both negative, d1 < d2", math.NewDecFromInt64(-10), math.NewDecFromInt64(-5), math.NewDecFromInt64(-10)},
		{"both negative, d1 > d2", math.NewDecFromInt64(-3), math.NewDecFromInt64(-7), math.NewDecFromInt64(-7)},
		{"equal negative", math.NewDecFromInt64(-1), math.NewDecFromInt64(-1), math.NewDecFromInt64(-1)},

		// Mixed sign
		{"d1 negative, d2 positive", math.NewDecFromInt64(-4), math.NewDecFromInt64(4), math.NewDecFromInt64(-4)},
		{"d1 positive, d2 negative", math.NewDecFromInt64(6), math.NewDecFromInt64(-6), math.NewDecFromInt64(-6)},

		// Zero + Negative
		{"d1 zero, d2 negative", utils.ZeroDec, math.NewDecFromInt64(-1), math.NewDecFromInt64(-1)},
		{"d1 negative, d2 zero", math.NewDecFromInt64(-1), utils.ZeroDec, math.NewDecFromInt64(-1)},

		// Zero + Positive (added)
		{"d1 zero, d2 positive", utils.ZeroDec, math.NewDecFromInt64(1), utils.ZeroDec},
		{"d1 positive, d2 zero", math.NewDecFromInt64(1), utils.ZeroDec, utils.ZeroDec},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.MinDec(tt.d1, tt.d2)
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %s, got %s", tt.expected.String(), result.String())
			}
		})
	}
}

func TestGetPaddedDecString(t *testing.T) {
	tests := []struct {
		name                 string
		dec                  math.Dec
		expectedFormattedStr string
		expectedPrefixLen    int // length before the decimal point (should be 32)
		expectedTotalLen     int // total string length (including dot and decimals)
		expectNegativeSign   bool
	}{
		{
			name:                 "small decimal",
			dec:                  math.NewDecWithExp(1, -34),
			expectedFormattedStr: "00000000000000000000000000000000.0000000000000000000000000000000001",
			expectedPrefixLen:    32,
			expectedTotalLen:     1 + 32 + 1 + int(utils.Precision), // optional '-' + padded + '.' + decimals
		},
		{
			name:                 "integer one",
			dec:                  utils.OneDec,
			expectedFormattedStr: "00000000000000000000000000000001.0000000000000000000000000000000000",
			expectedPrefixLen:    32,
			expectedTotalLen:     32 + 1 + int(utils.Precision),
		},
		{
			name:                 "large number",
			dec:                  math.NewDecWithExp(12345678901234567890123456789012345678, -8),
			expectedFormattedStr: "00000000000012345678901234567890.1234567800000000000000000000000000",
			expectedPrefixLen:    32,
			expectedTotalLen:     32 + 1 + int(utils.Precision),
		},
		{
			name:                 "negative decimal",
			dec:                  math.NewDecWithExp(-987654321, -6),
			expectedFormattedStr: "-00000000000000000000000000000987.6543210000000000000000000000000000",
			expectedPrefixLen:    32,
			expectedTotalLen:     1 + 32 + 1 + int(utils.Precision), // '-' + padded + '.' + decimals
			expectNegativeSign:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GetPaddedDecString(tt.dec)
			t.Logf("Formatted Result: %s", result)

			if !strings.Contains(result, ".") {
				t.Errorf("Result %s does not contain decimal point", result)
			}

			// Validate exact match
			if result != tt.expectedFormattedStr {
				t.Errorf("Expected formatted string:\n%s\ngot:\n%s", tt.expectedFormattedStr, result)
			}

			isNegative := strings.HasPrefix(result, "-")
			if isNegative != tt.expectNegativeSign {
				t.Errorf("Expected negative sign: %v, got result: %s", tt.expectNegativeSign, result)
			}

			parts := strings.Split(result, ".")
			if len(parts[0]) != tt.expectedPrefixLen {
				t.Errorf("Expected prefix length %d, got %d", tt.expectedPrefixLen, len(parts[0]))
			}
			if len(parts[1]) != int(utils.Precision) {
				t.Errorf("Expected decimal length %d, got %d", utils.Precision, len(parts[1]))
			}
			if len(result) != tt.expectedTotalLen {
				t.Errorf("Expected total length %d, got %d", tt.expectedTotalLen, len(result))
			}
		})
	}
}
