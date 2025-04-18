package types_test

import (
	"github.com/elys-network/elys/x/clob/types"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

// TestUpdateTotalOpenInterest_ExtendedMergedCases covers all scenarios including fractional values and various flip/increase/reduction combinations.
func TestUpdateTotalOpenInterest_ExtendedMergedCases(t *testing.T) {
	cases := []struct {
		name          string
		buyerBefore   sdkmath.LegacyDec
		sellerBefore  sdkmath.LegacyDec
		tradeSize     sdkmath.LegacyDec
		expectedDelta sdkmath.LegacyDec // expected net change in Market OI
	}{
		// Original cases:
		{
			name:          "Case 1: Both new positions",
			buyerBefore:   sdkmath.LegacyZeroDec(),
			sellerBefore:  sdkmath.LegacyZeroDec(),
			tradeSize:     sdkmath.LegacyNewDec(10),
			expectedDelta: sdkmath.LegacyNewDec(10),
		},
		{
			name:          "Case 2: Buyer adds long, seller opens short",
			buyerBefore:   sdkmath.LegacyNewDec(5),
			sellerBefore:  sdkmath.LegacyZeroDec(),
			tradeSize:     sdkmath.LegacyNewDec(5),
			expectedDelta: sdkmath.LegacyNewDec(5),
		},
		{
			name:          "Case 3: Buyer opens long, seller adds short",
			buyerBefore:   sdkmath.LegacyZeroDec(),
			sellerBefore:  sdkmath.LegacyNewDec(-3),
			tradeSize:     sdkmath.LegacyNewDec(7),
			expectedDelta: sdkmath.LegacyNewDec(7),
		},
		{
			name:          "Case 4: Buyer adds long, seller reduces long (both positive)",
			buyerBefore:   sdkmath.LegacyNewDec(5),
			sellerBefore:  sdkmath.LegacyNewDec(10),
			tradeSize:     sdkmath.LegacyNewDec(5),
			expectedDelta: sdkmath.LegacyZeroDec(),
		},
		{
			name:          "Case 5: Both negative positions offset",
			buyerBefore:   sdkmath.LegacyNewDec(-8),
			sellerBefore:  sdkmath.LegacyNewDec(-2),
			tradeSize:     sdkmath.LegacyNewDec(2),
			expectedDelta: sdkmath.LegacyZeroDec(),
		},
		{
			name:          "Case 6: Both closing positions",
			buyerBefore:   sdkmath.LegacyNewDec(-4),
			sellerBefore:  sdkmath.LegacyNewDec(4),
			tradeSize:     sdkmath.LegacyNewDec(4),
			expectedDelta: sdkmath.LegacyNewDec(-4),
		},
		{
			name:          "Case 7: Buyer flips from short to long, seller opens short",
			buyerBefore:   sdkmath.LegacyNewDec(-5),
			sellerBefore:  sdkmath.LegacyZeroDec(),
			tradeSize:     sdkmath.LegacyNewDec(10),
			expectedDelta: sdkmath.LegacyNewDec(5),
		},
		{
			name:          "Case 8: Both flip positions (net negative)",
			buyerBefore:   sdkmath.LegacyNewDec(-10),
			sellerBefore:  sdkmath.LegacyNewDec(10),
			tradeSize:     sdkmath.LegacyNewDec(15),
			expectedDelta: sdkmath.LegacyNewDec(-5),
		},
		{
			name:          "Case 9: Both increase existing positions",
			buyerBefore:   sdkmath.LegacyNewDec(5),
			sellerBefore:  sdkmath.LegacyNewDec(-5),
			tradeSize:     sdkmath.LegacyNewDec(10),
			expectedDelta: sdkmath.LegacyNewDec(10),
		},
		{
			name:          "Case 10: Both reduce positions",
			buyerBefore:   sdkmath.LegacyNewDec(-5),
			sellerBefore:  sdkmath.LegacyNewDec(5),
			tradeSize:     sdkmath.LegacyNewDec(2),
			expectedDelta: sdkmath.LegacyNewDec(-2),
		},
		{
			name:          "Case 11: Seller flips from long to short, buyer goes flat→long",
			buyerBefore:   sdkmath.LegacyZeroDec(),
			sellerBefore:  sdkmath.LegacyNewDec(10),
			tradeSize:     sdkmath.LegacyNewDec(15),
			expectedDelta: sdkmath.LegacyNewDec(5),
		},
		{
			name:          "Case 12: Both positive; buyer 0, seller reduces to 0",
			buyerBefore:   sdkmath.LegacyZeroDec(),
			sellerBefore:  sdkmath.LegacyNewDec(10),
			tradeSize:     sdkmath.LegacyNewDec(10),
			expectedDelta: sdkmath.LegacyZeroDec(),
		},
		{
			name:          "Case 13: Full mixed flip (net zero)",
			buyerBefore:   sdkmath.LegacyNewDec(-5),
			sellerBefore:  sdkmath.LegacyNewDec(10),
			tradeSize:     sdkmath.LegacyNewDec(15),
			expectedDelta: sdkmath.LegacyZeroDec(),
		},
		{
			name:          "Case 14: Mixed: buyer negative, seller zero",
			buyerBefore:   sdkmath.LegacyNewDec(-5),
			sellerBefore:  sdkmath.LegacyZeroDec(),
			tradeSize:     sdkmath.LegacyNewDec(2),
			expectedDelta: sdkmath.LegacyZeroDec(),
		},
		{
			name:          "Case 15: Mixed: buyer zero, seller positive",
			buyerBefore:   sdkmath.LegacyZeroDec(),
			sellerBefore:  sdkmath.LegacyNewDec(10),
			tradeSize:     sdkmath.LegacyNewDec(5),
			expectedDelta: sdkmath.LegacyZeroDec(),
		},
		// New additional cases:
		{
			name:         "Case 16: Fractional test",
			buyerBefore:  sdkmath.LegacyMustNewDecFromStr("2.5"),
			sellerBefore: sdkmath.LegacyMustNewDecFromStr("-3.75"),
			tradeSize:    sdkmath.LegacyMustNewDecFromStr("1.25"),
			// buyerAfter=3.75, sellerAfter=-5; OI: before=2.5+3.75=6.25, after=3.75+5=8.75, delta=(8.75-6.25)/2 = 1.25.
			expectedDelta: sdkmath.LegacyMustNewDecFromStr("1.25"),
		},
		{
			name:         "Case 17: Seller flip only (net positive)",
			buyerBefore:  sdkmath.LegacyNewDec(10),
			sellerBefore: sdkmath.LegacyNewDec(5),
			tradeSize:    sdkmath.LegacyNewDec(8),
			// buyerAfter = 18, sellerAfter = 5-8 = -3; OI before = 10+5=15, after=18+3=21, delta = (21-15)/2 = 3.
			expectedDelta: sdkmath.LegacyNewDec(3),
		},
		{
			name:         "Case 18: Opening new on one side & seller flip",
			buyerBefore:  sdkmath.LegacyZeroDec(),
			sellerBefore: sdkmath.LegacyNewDec(4),
			tradeSize:    sdkmath.LegacyNewDec(6),
			// buyerAfter = 6, sellerAfter = 4-6 = -2; OI: before=0+4=4, after=6+2=8, delta = (8-4)/2 = 2.
			expectedDelta: sdkmath.LegacyNewDec(2),
		},
		{
			name:         "Case 19: Increasing & seller flip (net negative)",
			buyerBefore:  sdkmath.LegacyNewDec(-4),
			sellerBefore: sdkmath.LegacyNewDec(2),
			tradeSize:    sdkmath.LegacyNewDec(3),
			// buyerAfter = -4+3 = -1, sellerAfter = 2-3 = -1; OI: before=4+2=6, after=1+1=2, delta = (2-6)/2 = -2.
			expectedDelta: sdkmath.LegacyNewDec(-2),
		},
		{
			name:         "Case 20: Reducing & seller flip (net negative)",
			buyerBefore:  sdkmath.LegacyNewDec(-8),
			sellerBefore: sdkmath.LegacyNewDec(5),
			tradeSize:    sdkmath.LegacyNewDec(7),
			// buyerAfter = -8+7 = -1, sellerAfter = 5-7 = -2; OI: before=8+5=13, after=1+2=3, delta = (3-13)/2 = -5.
			expectedDelta: sdkmath.LegacyNewDec(-5),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			market := &types.PerpetualMarket{
				TotalOpen: sdkmath.LegacyZeroDec(),
			}
			market.UpdateTotalOpenInterest(tc.buyerBefore, tc.sellerBefore, tc.tradeSize)
			require.True(t, market.TotalOpen.Equal(tc.expectedDelta),
				"[%s] expected TotalOpen: %s, got: %s",
				tc.name, tc.expectedDelta.String(), market.TotalOpen.String())
		})
	}
}
func TestUpdateTotalOpenInterest_InvalidTradeSize(t *testing.T) {
	invalidTradeSizes := []sdkmath.LegacyDec{
		sdkmath.LegacyZeroDec(),
		sdkmath.LegacyNewDec(-5),
	}
	for _, size := range invalidTradeSizes {
		t.Run("panic on tradeSize "+size.String(), func(t *testing.T) {
			market := &types.PerpetualMarket{
				TotalOpen: sdkmath.LegacyZeroDec(),
			}
			require.Panics(t, func() {
				market.UpdateTotalOpenInterest(sdkmath.LegacyNewDec(5), sdkmath.LegacyNewDec(10), size)
			}, "expected panic for trade size %s", size.String())
		})
	}
}
