package types_test

import (
	"testing"

	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestGetTradingAsset_WhenCollateralIsBaseCurrency(t *testing.T) {
	// Test case: collateral is base currency and borrow is ATOM
	result := types.GetTradingAsset(ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: both collateral and borrow are base currency
	result = types.GetTradingAsset(ptypes.BaseCurrency, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.Equal(t, ptypes.BaseCurrency, result)

	// Test case: collateral is base currency and borrow is some other asset (e.g., BTC)
	result = types.GetTradingAsset(ptypes.BaseCurrency, "BTC", ptypes.BaseCurrency)
	assert.Equal(t, "BTC", result)
}

func TestGetTradingAsset_WhenCollateralIsNotBaseCurrency(t *testing.T) {
	// Test case: collateral is ATOM and borrow is base currency
	result := types.GetTradingAsset(ptypes.ATOM, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: both collateral and borrow are ATOM
	result = types.GetTradingAsset(ptypes.ATOM, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Equal(t, ptypes.ATOM, result)

	// Test case: collateral is some other asset (e.g., BTC) and borrow is base currency
	result = types.GetTradingAsset("BTC", ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.Equal(t, "BTC", result)
}
