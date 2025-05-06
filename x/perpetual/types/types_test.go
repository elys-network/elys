package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/stretchr/testify/assert"
)

func TestCheckForStopLoss(t *testing.T) {
	mtp := types.MTP{}

	mtp.Position = types.Position_LONG
	mtp.StopLossPrice = math.LegacyNewDec(10)
	assert.False(t, mtp.CheckForStopLoss(osmomath.MustNewBigDecFromStr("10.1"))) // Above StopLossPrice
	assert.True(t, mtp.CheckForStopLoss(osmomath.NewBigDec(10)))                 // Equal to StopLossPrice
	assert.True(t, mtp.CheckForStopLoss(osmomath.MustNewBigDecFromStr("9.9")))   // Below StopLossPrice

	mtp.Position = types.Position_SHORT
	mtp.StopLossPrice = math.LegacyNewDec(10)
	assert.False(t, mtp.CheckForStopLoss(osmomath.MustNewBigDecFromStr("9.9"))) // Below StopLossPrice
	assert.True(t, mtp.CheckForStopLoss(osmomath.NewBigDec(10)))                // Equal to StopLossPrice
	assert.True(t, mtp.CheckForStopLoss(osmomath.MustNewBigDecFromStr("10.1"))) // Above StopLossPrice

	assert.False(t, mtp.CheckForTakeProfit(osmomath.NewBigDec(10))) // Should always return false

	// Edge case: Very high or low StopLossPrice
	mtp.StopLossPrice = math.LegacyNewDec(1e6) // Unrealistically high stop loss
	assert.False(t, mtp.CheckForStopLoss(osmomath.NewBigDec(10)))

	mtp.StopLossPrice = math.LegacyNewDec(-1e6) // Unrealistically low stop loss
	assert.True(t, mtp.CheckForStopLoss(osmomath.NewBigDec(10)))

	// Test unknown position
	mtp.Position = -1 // Invalid position
	mtp.TakeProfitPrice = math.LegacyNewDec(10)
}

func TestCheckForTakeProfit(t *testing.T) {
	mtp := types.MTP{}

	// Test LONG position
	mtp.Position = types.Position_LONG
	mtp.TakeProfitPrice = math.LegacyNewDec(15)

	assert.False(t, mtp.CheckForTakeProfit(osmomath.MustNewBigDecFromStr("14.9"))) // Below TakeProfitPrice
	assert.True(t, mtp.CheckForTakeProfit(osmomath.NewBigDec(15)))                 // Equal to TakeProfitPrice
	assert.True(t, mtp.CheckForTakeProfit(osmomath.MustNewBigDecFromStr("15.1")))  // Above TakeProfitPrice

	// Test SHORT position
	mtp.Position = types.Position_SHORT
	mtp.TakeProfitPrice = math.LegacyNewDec(10)

	assert.False(t, mtp.CheckForTakeProfit(osmomath.MustNewBigDecFromStr("10.1"))) // Above TakeProfitPrice
	assert.True(t, mtp.CheckForTakeProfit(osmomath.NewBigDec(10)))                 // Equal to TakeProfitPrice
	assert.True(t, mtp.CheckForTakeProfit(osmomath.MustNewBigDecFromStr("9.9")))   // Below TakeProfitPrice

	// Test unknown position
	mtp.Position = -1 // Invalid position
	mtp.TakeProfitPrice = math.LegacyNewDec(10)

	assert.False(t, mtp.CheckForTakeProfit(osmomath.NewBigDec(10))) // Should always return false

	// Edge case: Very high or low TakeProfitPrice
	mtp.Position = types.Position_LONG
	mtp.TakeProfitPrice = math.LegacyNewDec(1e6) // Unrealistically high take profit price
	assert.False(t, mtp.CheckForTakeProfit(osmomath.NewBigDec(10)))

	mtp.TakeProfitPrice = math.LegacyNewDec(-1e6) // Unrealistically low take profit price
	assert.True(t, mtp.CheckForTakeProfit(osmomath.NewBigDec(10)))
}
