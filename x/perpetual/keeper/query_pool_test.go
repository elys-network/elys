package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	simapp "github.com/elys-network/elys/app"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
)

func TestPools_InvalidRequest(t *testing.T) {
	mockAmm := new(mocks.AmmKeeper)
	k := keeper.NewKeeper(nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, nil, nil)
	ctx := sdk.Context{}
	_, err := k.Pools(ctx, nil)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
}

func TestPools_ErrPoolDoesNotExist(t *testing.T) {

	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	app.PerpetualKeeper.SetPool(ctx, types.Pool{
		AmmPoolId: uint64(23),
	})

	_, err := app.PerpetualKeeper.Pools(ctx, &types.QueryAllPoolRequest{})
	assert.Equal(t, "rpc error: code = Internal desc = pool does not exist", err.Error())
}

func TestPools_Success(t *testing.T) {

	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

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

func TestPoolQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.PerpetualKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPoolResponse(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPoolRequest
		response *types.QueryGetPoolResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPoolRequest{
				Index: msgs[0].AmmPoolId,
			},
			response: &types.QueryGetPoolResponse{Pool: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPoolRequest{
				Index: msgs[1].AmmPoolId,
			},
			response: &types.QueryGetPoolResponse{Pool: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPoolRequest{
				Index: (uint64)(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Pool(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
