package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/amm/types/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetTokenARate(t *testing.T) {
	ctx := sdk.Context{}
	accKeeper := mocks.NewAccountedPoolKeeper(t)

	// Define test cases
	testCases := []struct {
		name           string
		setupMock      func(oracleKeeper *mocks.OracleKeeper)
		pool           *types.Pool
		tokenA         string
		tokenB         string
		expectedRate   elystypes.Dec34
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
			elystypes.FourDec34().Quo(elystypes.ThreeDec34()),
			"",
		},
		{
			"oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(elystypes.NewDec34FromInt64(10), uint64(0))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenB").Return(elystypes.NewDec34FromInt64(5), uint64(0))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			elystypes.TwoDec34(),
			"",
		},
		{
			"token price not set for tokenA",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "unknownToken").Return(elystypes.ZeroDec34(), uint64(0))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"unknownToken",
			"tokenB",
			elystypes.ZeroDec34(),
			"token price not set: unknownToken",
		},
		{
			"token price not set for tokenB",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(elystypes.NewDec34FromInt64(5), uint64(0))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "unknownToken").Return(elystypes.ZeroDec34(), uint64(0))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"unknownToken",
			elystypes.ZeroDec34(),
			"token price not set: unknownToken",
		},
		{
			"Success with oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(elystypes.NewDec34FromInt64(5), uint64(0))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenB").Return(elystypes.NewDec34FromInt64(2), uint64(0))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			elystypes.FiveDec34().Quo(elystypes.TwoDec34()),
			"",
		},
		{
			"Success with oracle pricing",
			func(oracleKeeper *mocks.OracleKeeper) {
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(elystypes.NewDec34FromInt64(5), uint64(0))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenB").Return(elystypes.NewDec34FromInt64(2), uint64(0))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			elystypes.FiveDec34().Quo(elystypes.TwoDec34()),
			"",
		},
		{
			"Success with oracle pricing with price less than 1",
			func(oracleKeeper *mocks.OracleKeeper) {
				// for 6 decimal tokens
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenA").Return(elystypes.NewDec34FromString("0.0000002"), uint64(0))
				oracleKeeper.On("GetAssetPriceFromDenom", mock.Anything, "tokenB").Return(elystypes.OneDec34(), uint64(0))
			},
			&types.Pool{
				PoolParams: types.PoolParams{UseOracle: true},
			},
			"tokenA",
			"tokenB",
			elystypes.NewDec34FromString("0.0000002"),
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oracleKeeper := mocks.NewOracleKeeper(t)
			tc.setupMock(oracleKeeper)

			rate, err := tc.pool.GetTokenARate(ctx, oracleKeeper, tc.pool, tc.tokenA, tc.tokenB, accKeeper)
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedRate.ToLegacyDec().String(), rate.ToLegacyDec().String())
			}
		})
	}
}
