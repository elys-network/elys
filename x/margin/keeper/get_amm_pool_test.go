package keeper_test

import (
	"errors"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestGetAmmPool_PoolNotFound(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	margin := app.MarginKeeper

	borrowAsset := "testAsset"

	_, err := margin.GetAmmPool(ctx, 1, borrowAsset)

	// Expect no error and the first pool ID to be returned
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
}

func TestGetAmmPool_PoolFound(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	margin := app.MarginKeeper

	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"

	expectedPool := ammtypes.Pool{
		PoolId:            1,
		Address:           "",
		RebalanceTreasury: "",
		PoolParams: ammtypes.PoolParams{
			UseOracle:                   false,
			ExternalLiquidityRatio:      sdk.NewDec(2),
			WeightBreakingFeeMultiplier: sdk.ZeroDec(),
			WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
			LpFeePortion:                sdk.ZeroDec(),
			StakingFeePortion:           sdk.ZeroDec(),
			WeightRecoveryFeePortion:    sdk.ZeroDec(),
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
	app.AmmKeeper.SetPool(ctx, expectedPool)

	pool, err := margin.GetAmmPool(ctx, 1, borrowAsset)

	// Expect no error and the correct pool to be returned
	assert.Nil(t, err)
	assert.Equal(t, expectedPool.PoolId, pool.PoolId)
}
