package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/amm/types/mocks"
	"github.com/osmosis-labs/osmosis/osmomath"
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
		expectedValue  osmomath.BigDec
		expectedErrMsg string
	}{
		{
			"successful exit value calculation",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(10))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDec(5))
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
			osmomath.NewBigDec(1660),
			"",
		},
		{
			"total shares is zero",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(10))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDec(5))
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
			osmomath.ZeroBigDec(),
			"amount too low",
		},
		{
			"exiting shares greater than total shares",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(10))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDec(5))
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
			osmomath.ZeroBigDec(),
			"shares is larger than the max amount",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oracleKeeper := mocks.NewOracleKeeper(t)
			accKeeper := mocks.NewAccountedPoolKeeper(t)
			tc.setupMock(oracleKeeper, accKeeper)

			value, _, _, err := tc.pool.CalcExitValueWithSlippage(ctx, oracleKeeper, accKeeper, tc.pool, tc.exitingShares, tc.tokenOutDenom, osmomath.OneBigDec(), true, types.DefaultParams())
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
		expectedBonus  osmomath.BigDec
		expectedErrMsg string
	}{
		{
			"successful exit with oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(10))
				accKeeper.On("GetAccountedBalance", mock.Anything, mock.Anything, "tokenA").Return(sdkmath.NewInt(1000))
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenB").Return(osmomath.NewBigDec(5))
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
				MinSlippage:                 sdkmath.LegacyMustNewDecFromStr("0.001"),
			},
			sdk.Coins{sdk.NewCoin("tokenA", sdkmath.NewInt(190))},
			osmomath.ZeroBigDec(),
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
				MinSlippage:               sdkmath.LegacyMustNewDecFromStr("0.001"),
			},
			sdk.Coins{},
			osmomath.ZeroBigDec(),
			"shares is larger than the max amount",
		},
		{
			"exiting shares greater than total shares",
			func(oracleKeeper *mocks.OracleKeeper, accKeeper *mocks.AccountedPoolKeeper) {
				oracleKeeper.On("GetDenomPrice", mock.Anything, "tokenA").Return(osmomath.NewBigDec(0))
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
				MinSlippage:               sdkmath.LegacyMustNewDecFromStr("0.001"),
			},
			sdk.Coins{},
			osmomath.ZeroBigDec(),
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
				MinSlippage:               sdkmath.LegacyMustNewDecFromStr("0.001"),
			},
			sdk.Coins{sdk.NewCoin("tokenA", sdkmath.NewInt(100)), sdk.NewCoin("tokenB", sdkmath.NewInt(200))},
			osmomath.ZeroBigDec(),
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oracleKeeper := mocks.NewOracleKeeper(t)
			accKeeper := mocks.NewAccountedPoolKeeper(t)
			tc.setupMock(oracleKeeper, accKeeper)

			exitCoins, weightBalanceBonus, _, _, _, _, err := tc.pool.CalcExitPool(ctx, oracleKeeper, tc.pool, accKeeper, tc.exitingShares, tc.tokenOutDenom, tc.params, osmomath.ZeroBigDec(), true)
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
