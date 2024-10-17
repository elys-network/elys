package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
)

func TestPools_InvalidRequest(t *testing.T) {
	mockAmm := new(mocks.AmmKeeper)
	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, nil, nil, nil)
	ctx := sdk.Context{}
	_, err := k.Pools(ctx, nil)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
}

func TestPools_ErrPoolDoesNotExist(t *testing.T) {

	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(23),
	})

	_, err := app.PerpetualKeeper.Pools(ctx, &types.QueryAllPoolRequest{})
	assert.Equal(t, "rpc error: code = Internal desc = pool does not exist", err.Error())
}

func TestPools_Success(t *testing.T) {

	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(1),
	})

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(2),
	})

	app.AmmKeeper.SetPool(ctx, ammtypes.Pool{
		PoolId: uint64(1),
		PoolParams: ammtypes.PoolParams{
			UseOracle: true,
		},
	})

	app.AmmKeeper.SetPool(ctx, ammtypes.Pool{
		PoolId: uint64(2),
		PoolParams: ammtypes.PoolParams{
			UseOracle: false,
		},
	})

	response, err := app.PerpetualKeeper.Pools(ctx, &types.QueryAllPoolRequest{})
	assert.Nil(t, err)
	assert.Len(t, response.Pool, 1)

}

func TestPooolQuerySingle(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(1),
	})

	_, err := app.PerpetualKeeper.Pool(ctx, &types.QueryGetPoolRequest{
		Index: uint64(1),
	})

	_, errNotFound := app.PerpetualKeeper.Pool(ctx, &types.QueryGetPoolRequest{
		Index: uint64(2),
	})

	_, errInvalidRequest := app.PerpetualKeeper.Pool(ctx, nil)

	assert.Nil(t, err)
	assert.Error(t, errNotFound)
	require.ErrorIs(t, status.Error(codes.InvalidArgument, "invalid request"), errInvalidRequest)

}
