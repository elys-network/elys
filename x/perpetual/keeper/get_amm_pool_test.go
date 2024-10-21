package keeper_test

import (
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func TestGetAmmPool_PoolNotFound(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	perpetual := app.PerpetualKeeper

	borrowAsset := "testAsset"

	_, err := perpetual.GetAmmPool(ctx, 1, borrowAsset)

	// Expect no error and the first pool ID to be returned
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
}

func TestGetAmmPool_PoolFound(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	perpetual := app.PerpetualKeeper

	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"

	expectedPool := ammtypes.Pool{
		PoolId:            1,
		Address:           "",
		RebalanceTreasury: "",
		PoolParams: ammtypes.PoolParams{
			UseOracle:                   false,
			ExternalLiquidityRatio:      sdkmath.LegacyNewDec(2),
			WeightBreakingFeeMultiplier: sdkmath.LegacyZeroDec(),
			WeightBreakingFeeExponent:   sdkmath.LegacyNewDecWithPrec(25, 1), // 2.5
			WeightRecoveryFeePortion:    sdkmath.LegacyNewDecWithPrec(10, 2), // 10%
			ThresholdWeightDifference:   sdkmath.LegacyZeroDec(),
			SwapFee:                     sdkmath.LegacyZeroDec(),
			FeeDenom:                    ptypes.BaseCurrency,
		},
		TotalShares: sdk.Coin{},
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  sdk.NewCoin(collateralAsset, sdkmath.NewInt(10000)),
				Weight: sdkmath.NewInt(10),
			},
			{
				Token:  sdk.NewCoin(borrowAsset, sdkmath.NewInt(10000)),
				Weight: sdkmath.NewInt(10),
			},
		},
		TotalWeight: sdkmath.ZeroInt(),
	}
	app.AmmKeeper.SetPool(ctx, expectedPool)

	pool, err := perpetual.GetAmmPool(ctx, 1, borrowAsset)

	// Expect no error and the correct pool to be returned
	assert.Nil(t, err)
	assert.Equal(t, expectedPool.PoolId, pool.PoolId)
}
