package keeper_test

import (
	"testing"

	ptypes "github.com/elys-network/elys/v6/x/parameter/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	simapp "github.com/elys-network/elys/v6/app"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/perpetual/keeper"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func TestPools_InvalidRequest(t *testing.T) {
	k := keeper.NewKeeper(nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, nil, nil, nil)
	ctx := sdk.Context{}
	_, err := k.Pools(ctx, nil)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
}

func TestPools_ErrPoolDoesNotExist(t *testing.T) {

	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(23),
	})

	_, err := app.PerpetualKeeper.Pools(ctx, &types.QueryAllPoolRequest{})
	assert.Equal(t, "rpc error: code = Internal desc = perpetual pool does not exist", err.Error())
}

func (suite *PerpetualKeeperTestSuite) TestPools_Success() {
	suite.ResetSuite()
	suite.SetupCoinPrices()
	app := suite.app
	ctx := suite.ctx

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(1),
		PoolAssetsLong: []types.PoolAsset{
			{
				AssetDenom: ptypes.ATOM,
			},
			{
				AssetDenom: ptypes.BaseCurrency,
			},
		},
	})

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(2),
		PoolAssetsLong: []types.PoolAsset{
			{
				AssetDenom: ptypes.ATOM,
			},
			{
				AssetDenom: ptypes.BaseCurrency,
			},
		},
	})

	app.AmmKeeper.SetPool(ctx, ammtypes.Pool{
		PoolId:  uint64(1),
		Address: ammtypes.NewPoolAddress(1).String(),
		PoolParams: ammtypes.PoolParams{
			UseOracle: true,
		},
		PoolAssets: make([]ammtypes.PoolAsset, 2),
	})

	app.AmmKeeper.SetPool(ctx, ammtypes.Pool{
		PoolId:  uint64(2),
		Address: ammtypes.NewPoolAddress(2).String(),
		PoolParams: ammtypes.PoolParams{
			UseOracle: false,
		},
		PoolAssets: make([]ammtypes.PoolAsset, 2),
	})

	response, err := app.PerpetualKeeper.Pools(ctx, &types.QueryAllPoolRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(response.Pool, 1)

}

func (suite *PerpetualKeeperTestSuite) TestPooolQuerySingle() {
	suite.ResetSuite()
	suite.SetupCoinPrices()
	app := suite.app
	ctx := suite.ctx

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(1),
		PoolAssetsLong: []types.PoolAsset{
			{
				AssetDenom: ptypes.ATOM,
			},
			{
				AssetDenom: ptypes.BaseCurrency,
			},
		},
	})

	_, err := app.PerpetualKeeper.Pool(ctx, &types.QueryGetPoolRequest{
		Index: uint64(1),
	})

	_, errNotFound := app.PerpetualKeeper.Pool(ctx, &types.QueryGetPoolRequest{
		Index: uint64(2),
	})

	_, errInvalidRequest := app.PerpetualKeeper.Pool(ctx, nil)

	suite.Require().NoError(err)
	suite.Require().Error(errNotFound)
	suite.Require().ErrorIs(errInvalidRequest, status.Error(codes.InvalidArgument, "invalid request"))

}
