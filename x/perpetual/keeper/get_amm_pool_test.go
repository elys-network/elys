package keeper_test

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func (suie *PerpetualKeeperTestSuite) TestGetAmmPool_PoolNotFound() {
	borrowAsset := "testAsset"

	_, err := suie.app.PerpetualKeeper.GetAmmPool(suie.ctx, 1, borrowAsset)

	// Expect no error and the first pool ID to be returned
	assert.True(suie.T(), errors.Is(err, types.ErrPoolDoesNotExist))
}

func (suie *PerpetualKeeperTestSuite) TestGetAmmPool_PoolFound() {
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
	err := suie.app.AmmKeeper.SetPool(suie.ctx, expectedPool)
	suie.Require().NoError(err)

	pool, err := suie.app.PerpetualKeeper.GetAmmPool(suie.ctx, 1, borrowAsset)

	// Expect no error and the correct pool to be returned
	assert.Nil(suie.T(), err)
	assert.Equal(suie.T(), expectedPool.PoolId, pool.PoolId)
}
