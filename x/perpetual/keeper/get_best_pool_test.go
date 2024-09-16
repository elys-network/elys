package keeper_test

import (
	"errors"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func TestGetBestPool_NoPoolID(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	perpetual := app.PerpetualKeeper

	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"

	// Expect an error about the pool not existing
	_, err := perpetual.GetBestPool(ctx, collateralAsset, borrowAsset)
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
}

func TestGetBestPool_ValidPoolID(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	perpetual := app.PerpetualKeeper

	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"

	pool := ammtypes.Pool{
		PoolId:            1,
		Address:           "",
		RebalanceTreasury: "",
		PoolParams: ammtypes.PoolParams{
			UseOracle:                   true,
			ExternalLiquidityRatio:      sdk.NewDec(2),
			WeightBreakingFeeMultiplier: sdk.ZeroDec(),
			WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
			WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
			ThresholdWeightDifference:   sdk.ZeroDec(),
			SwapFee:                     sdk.ZeroDec(),
			FeeDenom:                    ptypes.BaseCurrency,
		},
		TotalShares: sdk.Coin{},
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  sdk.NewCoin(collateralAsset, sdk.NewInt(10000)),
				Weight: sdk.NewInt(10),
			},
			{
				Token:  sdk.NewCoin(borrowAsset, sdk.NewInt(10000)),
				Weight: sdk.NewInt(10),
			},
		},
		TotalWeight: sdk.ZeroInt(),
	}
	app.AmmKeeper.SetPool(ctx, pool)

	poolID, err := perpetual.GetBestPool(ctx, collateralAsset, borrowAsset)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), poolID)
}
