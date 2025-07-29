package types_test

import (
	"github.com/elys-network/elys/v7/x/clob/types"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestUpdateTotalOpenInterest(t *testing.T) {
	tradeSize := math.LegacyNewDec(10)
	initialOI := math.LegacyNewDec(100) // Example starting OI

	testCases := []struct {
		name             string
		initialOI        math.LegacyDec
		buyerBefore      math.LegacyDec
		sellerBefore     math.LegacyDec
		tradeSize        math.LegacyDec
		expectedFinalOI  math.LegacyDec // Expected market.TotalOpen after call
		expectPanic      bool
		expectedPanicMsg string
	}{
		// --- Panic Cases ---
		{
			name:             "Panic on zero trade size",
			initialOI:        initialOI,
			buyerBefore:      math.LegacyZeroDec(),
			sellerBefore:     math.LegacyZeroDec(),
			tradeSize:        math.LegacyZeroDec(), // Invalid
			expectedFinalOI:  initialOI,            // Should not change
			expectPanic:      true,
			expectedPanicMsg: "trade size cannot be 0 or negative",
		},
		{
			name:             "Panic on negative trade size",
			initialOI:        initialOI,
			buyerBefore:      math.LegacyZeroDec(),
			sellerBefore:     math.LegacyZeroDec(),
			tradeSize:        math.LegacyNewDec(-5), // Invalid
			expectedFinalOI:  initialOI,             // Should not change
			expectPanic:      true,
			expectedPanicMsg: "trade size cannot be 0 or negative",
		},

		// --- OI Increase Cases (+tradeSize) ---
		{
			name:            "OI Increase: Buyer Opens (+), Seller Opens (-)",
			initialOI:       initialOI,
			buyerBefore:     math.LegacyZeroDec(),
			sellerBefore:    math.LegacyZeroDec(),
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI.Add(tradeSize), // 100 + 10 = 110
			expectPanic:     false,
		},
		{
			name:            "OI Increase: Buyer Increases Long (+), Seller Increases Short (-)",
			initialOI:       initialOI,
			buyerBefore:     math.LegacyNewDec(20),  // +20 -> +30
			sellerBefore:    math.LegacyNewDec(-30), // -30 -> -40
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI.Add(tradeSize), // 100 + 10 = 110
			expectPanic:     false,
		},
		{
			name:            "OI Increase: Buyer Opens (+), Seller Increases Short (-)",
			initialOI:       math.LegacyNewDec(50),  // Start lower for clarity
			buyerBefore:     math.LegacyZeroDec(),   // 0 -> +10
			sellerBefore:    math.LegacyNewDec(-30), // -30 -> -40
			tradeSize:       tradeSize,
			expectedFinalOI: math.LegacyNewDec(60), // 50 + 10 = 60
			expectPanic:     false,
		},
		{
			name:            "OI Increase: Buyer Increases Long (+), Seller Opens (-)",
			initialOI:       math.LegacyNewDec(50),
			buyerBefore:     math.LegacyNewDec(25), // +25 -> +35
			sellerBefore:    math.LegacyZeroDec(),  // 0 -> -10
			tradeSize:       tradeSize,
			expectedFinalOI: math.LegacyNewDec(60), // 50 + 10 = 60
			expectPanic:     false,
		},

		// --- OI Decrease Cases (-tradeSize) ---
		{
			name:            "OI Decrease: Buyer Closes Short (+), Seller Closes Long (-)",
			initialOI:       initialOI,
			buyerBefore:     tradeSize.Neg(), // -10 -> 0
			sellerBefore:    tradeSize,       // +10 -> 0
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI.Sub(tradeSize), // 100 - 10 = 90
			expectPanic:     false,
		},
		{
			name:            "OI Decrease: Buyer Decreases Short (+), Seller Decreases Long (-)",
			initialOI:       initialOI,
			buyerBefore:     math.LegacyNewDec(-25), // -25 -> -15
			sellerBefore:    math.LegacyNewDec(35),  // +35 -> +25
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI.Sub(tradeSize), // 100 - 10 = 90
			expectPanic:     false,
		},
		{
			name:            "OI Decrease: Buyer Closes Short (+), Seller Decreases Long (-)",
			initialOI:       initialOI,
			buyerBefore:     tradeSize.Neg(),       // -10 -> 0
			sellerBefore:    math.LegacyNewDec(35), // +35 -> +25
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI.Sub(tradeSize), // 100 - 10 = 90
			expectPanic:     false,
		},
		{
			name:            "OI Decrease: Buyer Decreases Short (+), Seller Closes Long (-)",
			initialOI:       initialOI,
			buyerBefore:     math.LegacyNewDec(-25), // -25 -> -15
			sellerBefore:    tradeSize,              // +10 -> 0
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI.Sub(tradeSize), // 100 - 10 = 90
			expectPanic:     false,
		},

		// --- OI Unchanged Cases (0 change) ---
		{
			name:            "OI Unchanged: Buyer Opens (+), Seller Closes Long (-)",
			initialOI:       initialOI,
			buyerBefore:     math.LegacyZeroDec(), // 0 -> +10
			sellerBefore:    tradeSize,            // +10 -> 0
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI, // 100 + 0 = 100
			expectPanic:     false,
		},
		{
			name:            "OI Unchanged: Buyer Increases Long (+), Seller Decreases Long (-)",
			initialOI:       initialOI,
			buyerBefore:     math.LegacyNewDec(20), // +20 -> +30
			sellerBefore:    math.LegacyNewDec(15), // +15 -> +5
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI, // 100 + 0 = 100
			expectPanic:     false,
		},
		{
			name:            "OI Unchanged: Buyer Decreases Short (+), Seller Opens (-)",
			initialOI:       initialOI,
			buyerBefore:     math.LegacyNewDec(-20), // -20 -> -10
			sellerBefore:    math.LegacyZeroDec(),   // 0 -> -10
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI, // 100 + 0 = 100
			expectPanic:     false,
		},
		{
			name:            "OI Unchanged: Buyer Closes Short (+), Seller Increases Short (-)",
			initialOI:       initialOI,
			buyerBefore:     tradeSize.Neg(),        // -10 -> 0
			sellerBefore:    math.LegacyNewDec(-15), // -15 -> -25
			tradeSize:       tradeSize,
			expectedFinalOI: initialOI, // 100 + 0 = 100
			expectPanic:     false,
		},
		// Flip cases test implicit unchanged OI based on magnitude comparison logic
		{
			name:            "OI Unchanged: Buyer Flips Short -> Long (+5), Seller Increases Short (-40)",
			initialOI:       initialOI,
			buyerBefore:     math.LegacyNewDec(-10),
			sellerBefore:    math.LegacyNewDec(-25),
			tradeSize:       math.LegacyNewDec(15), // Buyer: -10 -> +5, Seller: -25 -> -40
			expectedFinalOI: initialOI,             // Buyer exposure increases, Seller exposure increases -> OI should increase? Let's recheck logic.
			// Buyer abs change: 5 vs 10 -> Decreased. NO -> Increased (0->5 is increase from crossing zero) abs(5)>abs(-10) is FALSE. abs(5)<abs(-10) is TRUE -> Buyer exposure DECREASED (by closing)
			// Seller abs change: 40 vs 25 -> Increased. abs(40)>abs(25) is TRUE -> Seller exposure INCREASED
			// One increase, one decrease -> OI UNCHANGED. Test case is correct.
			expectPanic: false,
		},
		{
			name:         "OI Unchanged: Seller Flips Long -> Short (-5), Buyer Increases Long (+35)",
			initialOI:    initialOI,
			buyerBefore:  math.LegacyNewDec(20),
			sellerBefore: math.LegacyNewDec(10),
			tradeSize:    math.LegacyNewDec(15), // Buyer: +20 -> +35, Seller: +10 -> -5
			// Buyer abs change: 35 vs 20 -> Increased.
			// Seller abs change: 5 vs 10 -> Decreased.
			// One increase, one decrease -> OI UNCHANGED. Test case is correct.
			expectedFinalOI: initialOI,
			expectPanic:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a market instance for each test
			market := types.PerpetualMarket{
				TotalOpen: tc.initialOI,
				// Initialize other fields if they affect the method, otherwise not needed
			}

			// Call the function with error handling
			err := market.UpdateTotalOpenInterest(tc.buyerBefore, tc.sellerBefore, tc.tradeSize)

			if tc.expectPanic {
				// Now we expect error instead of panic
				require.Error(t, err, "Expected error for invalid input")
				require.Contains(t, err.Error(), "invalid trade size", "Error message should mention invalid trade size")
				// Verify OI did not change after error
				require.True(t, tc.initialOI.Equal(market.TotalOpen), "TotalOpen changed despite error. Initial: %s, Final: %s", tc.initialOI, market.TotalOpen)
			} else {
				require.NoError(t, err, "Function call returned unexpected error")
				// Verify final OI matches expectation
				require.True(t, tc.expectedFinalOI.Equal(market.TotalOpen), "Final TotalOpen mismatch. Expected: %s, Got: %s", tc.expectedFinalOI, market.TotalOpen)
			}
		})
	}
}
func TestUpdateTotalOpenInterest_InvalidTradeSize(t *testing.T) {
	invalidTradeSizes := []math.LegacyDec{
		math.LegacyZeroDec(),
		math.LegacyNewDec(-5),
	}
	for _, size := range invalidTradeSizes {
		t.Run("error on tradeSize "+size.String(), func(t *testing.T) {
			market := &types.PerpetualMarket{
				TotalOpen: math.LegacyZeroDec(),
			}
			err := market.UpdateTotalOpenInterest(math.LegacyNewDec(5), math.LegacyNewDec(10), size)
			require.Error(t, err, "expected error for trade size %s", size.String())
			require.Contains(t, err.Error(), "invalid trade size", "Error should mention invalid trade size")
		})
	}
}
