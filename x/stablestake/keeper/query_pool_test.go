package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *KeeperTestSuite) TestAmmPool() {
	p := types.AmmPool{
		Id:               1,
		TotalLiabilities: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)},
	}
	tests := []struct {
		name          string
		req           *types.QueryAmmPoolRequest
		setup         func(ctx sdk.Context, k keeper.Keeper)
		expectedError error
		expectedResp  *types.QueryAmmPoolResponse
	}{
		{
			name: "valid request",
			req:  &types.QueryAmmPoolRequest{Id: 1},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetAmmPool(ctx, p)
			},
			expectedError: nil,
			expectedResp: &types.QueryAmmPoolResponse{
				AmmPool: p,
			},
		},
		{
			name:          "invalid request",
			req:           nil,
			setup:         func(ctx sdk.Context, k keeper.Keeper) {},
			expectedError: status.Error(codes.InvalidArgument, "invalid request"),
			expectedResp:  nil,
		},
		{
			name: "pool not exists",
			req:  &types.QueryAmmPoolRequest{Id: 2},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
			},
			expectedError: nil,
			expectedResp: &types.QueryAmmPoolResponse{
				AmmPool: types.AmmPool{
					Id:               2,
					TotalLiabilities: sdk.Coins{},
				},
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.setup != nil {
				tt.setup(suite.ctx, *suite.app.StablestakeKeeper)
			}

			resp, err := suite.app.StablestakeKeeper.AmmPool(suite.ctx, tt.req)
			if tt.expectedError != nil {
				require.ErrorIs(suite.T(), err, tt.expectedError)
			} else {
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), tt.expectedResp, resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAllAmmPools() {
	p := types.AmmPool{
		Id:               1,
		TotalLiabilities: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)},
	}
	tests := []struct {
		name          string
		req           *types.QueryAllAmmPoolsRequest
		setup         func(ctx sdk.Context, k keeper.Keeper)
		expectedError error
		expectedResp  *types.QueryAllAmmPoolsResponse
	}{
		{
			name: "valid request",
			req:  &types.QueryAllAmmPoolsRequest{},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetAmmPool(ctx, p)
			},
			expectedError: nil,
			expectedResp: &types.QueryAllAmmPoolsResponse{
				AmmPools: []types.AmmPool{p},
			},
		},
		{
			name:          "invalid request",
			req:           nil,
			setup:         func(ctx sdk.Context, k keeper.Keeper) {},
			expectedError: status.Error(codes.InvalidArgument, "invalid request"),
			expectedResp:  nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.setup != nil {
				tt.setup(suite.ctx, *suite.app.StablestakeKeeper)
			}

			resp, err := suite.app.StablestakeKeeper.AllAmmPools(suite.ctx, tt.req)
			if tt.expectedError != nil {
				require.ErrorIs(suite.T(), err, tt.expectedError)
			} else {
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), tt.expectedResp, resp)
			}
		})
	}
}
