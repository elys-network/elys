package keeper_test

import (
	"testing"

	"github.com/elys-network/elys/x/margin/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestGetTradingAsset_WhenCollateralIsBaseCurrency(t *testing.T) {
	// Create an instance of Keeper
	k := keeper.Keeper{}

	// Test case: collateral is base currency and borrow is ATOM
	result := k.GetTradingAsset(ptypes.BaseCurrency, ptypes.ATOM)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: both collateral and borrow are base currency
	result = k.GetTradingAsset(ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.Equal(t, ptypes.BaseCurrency, result)

	// Test case: collateral is base currency and borrow is some other asset (e.g., BTC)
	result = k.GetTradingAsset(ptypes.BaseCurrency, "BTC")
	assert.Equal(t, "BTC", result)
}

func TestGetTradingAsset_WhenCollateralIsNotBaseCurrency(t *testing.T) {
	// Create an instance of Keeper
	k := keeper.Keeper{}

	// Test case: collateral is ATOM and borrow is base currency
	result := k.GetTradingAsset(ptypes.ATOM, ptypes.BaseCurrency)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: both collateral and borrow are ATOM
	result = k.GetTradingAsset(ptypes.ATOM, ptypes.ATOM)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: collateral is some other asset (e.g., BTC) and borrow is base currency
	result = k.GetTradingAsset("BTC", ptypes.BaseCurrency)
	assert.Equal(t, "BTC", result)
}
