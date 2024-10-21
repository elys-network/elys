package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func TestGetBestPool_NoPoolID(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	perpetual := app.PerpetualKeeper

	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"

	// Expect an error about the pool not existing
	_, err := perpetual.GetBestPool(ctx, collateralAsset, borrowAsset)
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
}

func TestGetBestPool_ValidPoolID(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	perpetual := app.PerpetualKeeper

	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"

	pool := ammtypes.Pool{
		PoolId:            1,
		Address:           "",
		RebalanceTreasury: "",
		PoolParams: ammtypes.PoolParams{
			UseOracle:                   true,
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
	app.AmmKeeper.SetPool(ctx, pool)

	poolID, err := perpetual.GetBestPool(ctx, collateralAsset, borrowAsset)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), poolID)
}