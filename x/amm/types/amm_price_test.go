package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/amm/types/mocks"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"
)

func TestGetTokenARate(t *testing.T) {
	ctx := sdk.Context{}

	// Define test cases
	testCases := []struct {
		name           string
		setupMock      func(oracleKeeper *mocks.OracleKeeper)
		pool           *types.Pool
		tokenA         string
		tokenB         string
		expectedRate   osmomath.BigDec
		expectedErrMsg string
	}{
		{
			"balancer pricing",
			func(oracleKeeper *mocks.OracleKeeper) {},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: false},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1500)), Weight: sdkmath.NewInt(1)},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000)), Weight: sdkmath.NewInt(1)},
				},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(4).Quo(osmomath.NewBigDec(3)),
			"",
		},
		{
			"oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(10))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDec(5))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(2),
			"",
		},
		{
			"token price not set for tokenA",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "unknownToken").Return(osmomath.ZeroBigDec())
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"unknownToken",
			"tokenB",
			osmomath.ZeroBigDec(),
			"token price not set or zero for token: unknownToken",
		},
		{
			"token price not set for tokenB",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(5))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "unknownToken").Return(osmomath.ZeroBigDec())
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"unknownToken",
			osmomath.ZeroBigDec(),
			"token price not set or zero for token: unknownToken",
		},
		{
			"Success with oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(5))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDec(2))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(5).Quo(osmomath.NewBigDec(2)),
			"",
		},
		{
			"Success with oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(5))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDec(2))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(5).Quo(osmomath.NewBigDec(2)),
			"",
		},
		{
			"Success with oracle pricing with price less than 1",
			func(oracleKeeper *mocks.OracleKeeper) {
				// for 6 decimal tokens
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.MustNewBigDecFromStr("0.0000002"))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDec(1))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			osmomath.MustNewBigDecFromStr("0.0000002"),
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oracleKeeper := mocks.NewOracleKeeper(t)
			tc.setupMock(oracleKeeper)

			rate, err := tc.pool.GetTokenARate(ctx, oracleKeeper, tc.tokenA, tc.tokenB)
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedRate, rate)
			}
		})
	}
}

var decimal6, _ = sdkmath.NewIntFromString("1000000")
var decimal18, _ = sdkmath.NewIntFromString("1000000000000000000")

func TestGetTokenARateNormalized(t *testing.T) {
	ctx := sdk.Context{}

	// Define test cases
	testCases := []struct {
		name           string
		setupMock      func(oracleKeeper *mocks.OracleKeeper)
		pool           *types.Pool
		tokenA         string
		tokenB         string
		expectedRate   osmomath.BigDec
		expectedErrMsg string
	}{
		{
			"balancer pricing with equal decimals",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenA").Return(oracletypes.AssetInfo{Decimal: 6}, true)
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenB").Return(oracletypes.AssetInfo{Decimal: 6}, true)
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: false},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1500).Mul(decimal6)), Weight: sdkmath.NewInt(1)},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000).Mul(decimal6)), Weight: sdkmath.NewInt(1)},
				},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(4).Quo(osmomath.NewBigDec(3)),
			"",
		},
		{
			"balancer pricing with different decimals - tokenB has more",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenA").Return(oracletypes.AssetInfo{Decimal: 6}, true)
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenB").Return(oracletypes.AssetInfo{Decimal: 18}, true)
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: false},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1500).Mul(decimal6)), Weight: sdkmath.NewInt(1)},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000).Mul(decimal18)), Weight: sdkmath.NewInt(1)},
				},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(4).Quo(osmomath.NewBigDec(3)),
			"",
		},
		{
			"balancer pricing with different decimals - tokenA has more",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenA").Return(oracletypes.AssetInfo{Decimal: 18}, true)
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenB").Return(oracletypes.AssetInfo{Decimal: 6}, true)
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: false},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1500).Mul(decimal18)), Weight: sdkmath.NewInt(1)},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000).Mul(decimal6)), Weight: sdkmath.NewInt(1)},
				},
			},
			"tokenA",
			"tokenB",
			osmomath.MustNewBigDecFromStr("1.333333333333333333333333000000000000"),
			"",
		},
		{
			"oracle pricing with equal decimals",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDecWithPrec(10, 6))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDecWithPrec(5, 6))
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenA").Return(oracletypes.AssetInfo{Decimal: 6}, true)
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenB").Return(oracletypes.AssetInfo{Decimal: 6}, true)
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(2),
			"",
		},
		{
			"oracle pricing with different decimals - tokenB has more",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDecWithPrec(10, 6))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDecWithPrec(5, 18))
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenA").Return(oracletypes.AssetInfo{Decimal: 6}, true)
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenB").Return(oracletypes.AssetInfo{Decimal: 18}, true)
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(2),
			"",
		},
		{
			"oracle pricing with different decimals - tokenA has more",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDecWithPrec(10, 18))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDecWithPrec(5, 6))
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenA").Return(oracletypes.AssetInfo{Decimal: 18}, true)
				oracleKeeper.On("GetAssetInfo", mock.Anything, "tokenB").Return(oracletypes.AssetInfo{Decimal: 6}, true)
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			osmomath.NewBigDec(2),
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oracleKeeper := mocks.NewOracleKeeper(t)
			tc.setupMock(oracleKeeper)

			rate, err := tc.pool.GetTokenARateNormalized(ctx, oracleKeeper, tc.tokenA, tc.tokenB)
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
				require.True(t, rate.IsZero())
			} else {
				require.NoError(t, err)
				require.True(t, rate.Equal(tc.expectedRate), "expected %s, got %s", tc.expectedRate.String(), rate.String())
			}
		})
	}
}
