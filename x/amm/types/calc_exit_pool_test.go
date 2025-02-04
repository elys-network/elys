package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/amm/types/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCalcExitValueWithSlippage(t *testing.T) {
	ctx := sdk.Context{}

	// Define test cases
	testCases := []struct {
		name           string
		setupMock      func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper)
		pool           types.Pool
		exitingShares  sdkmath.Int
		tokenOutDenom  string
		expectedValue  sdkmath.LegacyDec
		expectedErrMsg string
	}{
		{
			"successful exit value calculation",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(sdkmath.LegacyNewDec(10))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenB").Return(sdkmath.LegacyNewDec(5))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenA").Return(sdkmath.NewInt(1000))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenB").Return(sdkmath.NewInt(2000))
			},
			types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
				},
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(100)),
			},
			sdkmath.NewInt(10),
			"tokenA",
			sdkmath.LegacyNewDec(1660),
			"",
		},
		{
			"total shares is zero",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(sdkmath.LegacyNewDec(10))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenB").Return(sdkmath.LegacyNewDec(5))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenA").Return(sdkmath.NewInt(1000))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenB").Return(sdkmath.NewInt(2000))
			},
			types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
				},
				TotalShares: sdk.NewCoin("shares", sdkmath.ZeroInt()),
			},
			sdkmath.NewInt(10),
			"tokenA",
			sdkmath.LegacyZeroDec(),
			"amount too low",
		},
		{
			"exiting shares greater than total shares",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(sdkmath.LegacyNewDec(10))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenB").Return(sdkmath.LegacyNewDec(5))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenA").Return(sdkmath.NewInt(1000))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenB").Return(sdkmath.NewInt(2000))
			},
			types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
				},
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(10)),
			},
			sdkmath.NewInt(100),
			"tokenA",
			sdkmath.LegacyZeroDec(),
			"shares is larger than the max amount",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oracleKeeper := mocks.NewOracleKeeper(t)
			accKeeper := mocks.NewAccountedPoolKeeper(t)
			tc.setupMock(oracleKeeper, accKeeper)

			value, _, err := types.CalcExitValueWithSlippage(ctx, oracleKeeper, accKeeper, tc.pool, tc.exitingShares, tc.tokenOutDenom, sdkmath.LegacyOneDec(), true)
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedValue, value)
			}

			oracleKeeper.AssertExpectations(t)
		})
	}
}

func TestCalcExitPool(t *testing.T) {
	ctx := sdk.Context{}

	// Define test cases
	testCases := []struct {
		name           string
		setupMock      func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper)
		pool           types.Pool
		exitingShares  sdkmath.Int
		tokenOutDenom  string
		params         types.Params
		expectedCoins  sdk.Coins
		expectedBonus  sdkmath.LegacyDec
		expectedErrMsg string
	}{
		{
			"successful exit with oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(sdkmath.LegacyNewDec(10))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenA").Return(sdkmath.NewInt(1000))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenB").Return(sdkmath.LegacyNewDec(5))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenB").Return(sdkmath.NewInt(2000))
			},
			types.Pool{
				PoolParams: types.PoolParams{UseOracle: true, SwapFee: sdkmath.LegacyZeroDec()},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
				},
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(100)),
			},
			sdkmath.NewInt(10),
			"tokenA",
			types.Params{
				WeightBreakingFeeMultiplier: sdkmath.LegacyMustNewDecFromStr("0.0005"),
				WeightBreakingFeePortion:    sdkmath.LegacyMustNewDecFromStr("0.5"),
				ThresholdWeightDifference:   sdkmath.LegacyMustNewDecFromStr("0.2"),
				WeightBreakingFeeExponent:   sdkmath.LegacyMustNewDecFromStr("0.5"),
			},
			sdk.Coins{sdk.NewCoin("tokenA", sdkmath.NewInt(190))},
			sdkmath.LegacyZeroDec(),
			"",
		},
		{
			"exiting shares greater than total shares",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
			},
			types.Pool{
				PoolParams:  types.PoolParams{UseOracle: true},
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(10)),
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
				},
			},
			sdkmath.NewInt(20),
			"tokenA",
			types.Params{
				WeightBreakingFeePortion:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				ThresholdWeightDifference: sdkmath.LegacyMustNewDecFromStr("0.2"),
				WeightBreakingFeeExponent: sdkmath.LegacyMustNewDecFromStr("0.5"),
			},
			sdk.Coins{},
			sdkmath.LegacyZeroDec(),
			"shares is larger than the max amount",
		},
		{
			"exiting shares greater than total shares",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(sdkmath.LegacyNewDec(0))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenA").Return(sdkmath.NewInt(1000))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenB").Return(sdkmath.NewInt(2000))
			},
			types.Pool{
				PoolParams:  types.PoolParams{UseOracle: true},
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(100)),
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000)), Weight: sdkmath.NewInt(1), ExternalLiquidityRatio: sdkmath.LegacyOneDec()},
				},
			},
			sdkmath.NewInt(10),
			"tokenA",
			types.Params{
				WeightBreakingFeePortion:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				ThresholdWeightDifference: sdkmath.LegacyMustNewDecFromStr("0.2"),
				WeightBreakingFeeExponent: sdkmath.LegacyMustNewDecFromStr("0.5"),
			},
			sdk.Coins{},
			sdkmath.LegacyZeroDec(),
			"token price not set",
		},
		{
			"successful exit without oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {},
			types.Pool{
				PoolParams: types.PoolParams{UseOracle: false},
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000)), Weight: sdkmath.NewInt(1)},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000)), Weight: sdkmath.NewInt(1)},
				},
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(100)),
			},
			sdkmath.NewInt(10),
			"",
			types.Params{
				WeightBreakingFeePortion:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				ThresholdWeightDifference: sdkmath.LegacyMustNewDecFromStr("0.2"),
				WeightBreakingFeeExponent: sdkmath.LegacyMustNewDecFromStr("0.5"),
			},
			sdk.Coins{sdk.NewCoin("tokenA", sdkmath.NewInt(100)), sdk.NewCoin("tokenB", sdkmath.NewInt(200))},
			sdkmath.LegacyZeroDec(),
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oracleKeeper := mocks.NewOracleKeeper(t)
			accKeeper := mocks.NewAccountedPoolKeeper(t)
			tc.setupMock(oracleKeeper, accKeeper)

			exitCoins, weightBalanceBonus, _, _, err := types.CalcExitPool(ctx, oracleKeeper, tc.pool, accKeeper, tc.exitingShares, tc.tokenOutDenom, tc.params, true)
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedCoins, exitCoins)
				require.Equal(t, tc.expectedBonus, weightBalanceBonus)
			}

			oracleKeeper.AssertExpectations(t)
			accKeeper.AssertExpectations(t)
		})
	}
}
