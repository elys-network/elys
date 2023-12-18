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

func TestGetFirstValidPool_NoPoolID(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	margin := app.MarginKeeper

	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"

	_, err := margin.GetFirstValidPool(ctx, collateralAsset, borrowAsset)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
}

func TestGetFirstValidPool_ValidPoolID(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	margin := app.MarginKeeper

	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"

	pool := ammtypes.Pool{
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
	app.AmmKeeper.SetPool(ctx, pool)

	poolID, err := margin.GetFirstValidPool(ctx, collateralAsset, borrowAsset)

	// Expect no error and the first pool ID to be returned
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), poolID)
}
