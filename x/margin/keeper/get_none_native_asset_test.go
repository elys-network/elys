package keeper_test

import (
	"testing"

	"github.com/elys-network/elys/x/margin/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestGetNonNativeAsset_WhenCollateralIsUSDC(t *testing.T) {
	// Create an instance of Keeper
	k := keeper.Keeper{}

	// Test case: collateral is USDC and borrow is ATOM
	result := k.GetNonNativeAsset(ptypes.USDC, ptypes.ATOM)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: both collateral and borrow are USDC
	result = k.GetNonNativeAsset(ptypes.USDC, ptypes.USDC)
	assert.Equal(t, ptypes.USDC, result)

	// Test case: collateral is USDC and borrow is some other asset (e.g., BTC)
	result = k.GetNonNativeAsset(ptypes.USDC, "BTC")
	assert.Equal(t, "BTC", result)
}

func TestGetNonNativeAsset_WhenCollateralIsNotUSDC(t *testing.T) {
	// Create an instance of Keeper
	k := keeper.Keeper{}

	// Test case: collateral is ATOM and borrow is USDC
	result := k.GetNonNativeAsset(ptypes.ATOM, ptypes.USDC)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: both collateral and borrow are ATOM
	result = k.GetNonNativeAsset(ptypes.ATOM, ptypes.ATOM)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: collateral is some other asset (e.g., BTC) and borrow is USDC
	result = k.GetNonNativeAsset("BTC", ptypes.USDC)
	assert.Equal(t, "BTC", result)
}
