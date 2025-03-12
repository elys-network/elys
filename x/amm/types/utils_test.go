package types_test

import (
	"errors"
	"fmt"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestGetPoolShareDenom(t *testing.T) {
	tests := []struct {
		poolId   uint64
		expected string
	}{
		{1, "amm/pool/1"},
		{42, "amm/pool/42"},
		{1000, "amm/pool/1000"},
	}

	for _, tt := range tests {
		result := types.GetPoolShareDenom(tt.poolId)
		if result != tt.expected {
			t.Errorf("GetPoolShareDenom(%d) = %s; want %s", tt.poolId, result, tt.expected)
		}
	}
}

func TestEnsureDenomInPool(t *testing.T) {
	poolAssetsByDenom := map[string]types.PoolAsset{
		"abc": {Token: sdk.NewCoin("abc", sdkmath.ZeroInt())},
		"def": {Token: sdk.NewCoin("def", sdkmath.ZeroInt())},
		"ghi": {Token: sdk.NewCoin("ghi", sdkmath.ZeroInt())},
	}

	tests := []struct {
		tokensIn sdk.Coins
		err      error
	}{
		{sdk.NewCoins(sdk.NewInt64Coin("abc", 100), sdk.NewInt64Coin("def", 200)), nil},
		{sdk.NewCoins(sdk.NewInt64Coin("def", 200), sdk.NewInt64Coin("jkl", 300)), errorsmod.Wrapf(types.ErrDenomNotFoundInPool, types.InvalidInputDenomsErrFormat, "jkl")},
		{sdk.NewCoins(sdk.NewInt64Coin("xyz", 500)), errorsmod.Wrapf(types.ErrDenomNotFoundInPool, types.InvalidInputDenomsErrFormat, "xyz")},
	}

	for _, tt := range tests {
		err := types.EnsureDenomInPool(poolAssetsByDenom, tt.tokensIn)
		if !errors.Is(err, tt.err) {
			t.Errorf("EnsureDenomInPool(%v) = %v; want %v", tt.tokensIn, err, tt.err)
		}
	}
}

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
		result, sign := types.AbsDifferenceWithSign(tt.a, tt.b)
		if !result.Equal(tt.expected) || sign != tt.sign {
			t.Errorf("AbsDifferenceWithSign(%s, %s) = (%s, %v); want (%s, %v)", tt.a, tt.b, result, sign, tt.expected, tt.sign)
		}
	}
}

func TestApplyDiscount(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		swapFee  sdkmath.LegacyDec
		discount sdkmath.LegacyDec
		wantFee  sdkmath.LegacyDec
	}{
		{
			name:     "Zero discount",
			swapFee:  sdkmath.LegacyNewDecWithPrec(100, 2), // 1.00 as an example
			discount: sdkmath.LegacyZeroDec(),
			wantFee:  sdkmath.LegacyNewDecWithPrec(100, 2),
		},
		{
			name:     "Positive discount",
			swapFee:  sdkmath.LegacyNewDecWithPrec(100, 2),
			discount: sdkmath.LegacyNewDecWithPrec(10, 2), // 0.10 (10%)
			wantFee:  sdkmath.LegacyNewDecWithPrec(90, 2), // 0.90 after discount
		},
		{
			name:     "Boundary value for discount",
			swapFee:  sdkmath.LegacyNewDecWithPrec(100, 2),
			discount: sdkmath.LegacyNewDecWithPrec(9999, 4), // 0.9999 (99.99%)
			wantFee:  sdkmath.LegacyNewDecWithPrec(1, 4),    // 0.01 after discount
		},
		{
			name:     "Discount greater than swap fee",
			swapFee:  sdkmath.LegacyNewDecWithPrec(50, 2), // 0.50
			discount: sdkmath.LegacyNewDecWithPrec(75, 2), // 0.75
			wantFee:  sdkmath.LegacyNewDecWithPrec(125, 3),
		},
		{
			name:     "Zero swap fee with valid discount",
			swapFee:  sdkmath.LegacyZeroDec(),
			discount: sdkmath.LegacyNewDecWithPrec(10, 2),
			wantFee:  sdkmath.LegacyZeroDec(),
		},
		{
			name:     "Large discount",
			swapFee:  sdkmath.LegacyNewDecWithPrec(100, 2),
			discount: sdkmath.LegacyNewDecWithPrec(9000, 4), // 0.90 (90%)
			wantFee:  sdkmath.LegacyNewDecWithPrec(10, 2),   // 0.10 after discount
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fee := types.ApplyDiscount(tc.swapFee, tc.discount)
			require.Equal(t, tc.wantFee, fee)
		})
	}
}

func TestGetPoolAssetsByDenom(t *testing.T) {
	poolAssets := []types.PoolAsset{
		{
			Token:  sdk.Coin{Denom: "token1", Amount: sdkmath.NewInt(100)},
			Weight: sdkmath.NewInt(10),
		},
		{
			Token:  sdk.Coin{Denom: "token2", Amount: sdkmath.NewInt(200)},
			Weight: sdkmath.NewInt(20),
		},
	}

	// Test case 1: No duplicate pool assets
	poolAssetsByDenom, err := types.GetPoolAssetsByDenom(poolAssets)
	require.NoError(t, err)
	require.Equal(t, 2, len(poolAssetsByDenom))
	require.Equal(t, poolAssets[0], poolAssetsByDenom["token1"])
	require.Equal(t, poolAssets[1], poolAssetsByDenom["token2"])

	// Test case 2: Duplicate pool asset
	duplicatePoolAssets := []types.PoolAsset{
		{
			Token:  sdk.Coin{Denom: "token1", Amount: sdkmath.NewInt(100)},
			Weight: sdkmath.NewInt(10),
		},
		{
			Token:  sdk.Coin{Denom: "token1", Amount: sdkmath.NewInt(200)},
			Weight: sdkmath.NewInt(20),
		},
	}
	_, err = types.GetPoolAssetsByDenom(duplicatePoolAssets)
	expectedErr := fmt.Errorf(types.FormatRepeatingPoolAssetsNotAllowedErrFormat, "token1")
	require.EqualError(t, err, expectedErr.Error())
}

func TestGetPoolAssetByDenom(t *testing.T) {
	poolAssets := []types.PoolAsset{
		{
			Token:  sdk.Coin{Denom: "token1", Amount: sdkmath.NewInt(100)},
			Weight: sdkmath.NewInt(10),
		},
		{
			Token:  sdk.Coin{Denom: "token2", Amount: sdkmath.NewInt(200)},
			Weight: sdkmath.NewInt(20),
		},
	}

	// Test case 1: Existing PoolAsset
	asset, found := types.GetPoolAssetByDenom(poolAssets, "token1")
	require.True(t, found)
	require.Equal(t, poolAssets[0], asset)

	// Test case 2: Non-existing PoolAsset
	_, found = types.GetPoolAssetByDenom(poolAssets, "nonExistingToken")
	require.False(t, found)
}
