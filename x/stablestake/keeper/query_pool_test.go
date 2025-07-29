package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/stablestake/keeper"
	"github.com/elys-network/elys/v7/x/stablestake/types"
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

func (suite *KeeperTestSuite) TestMaxBondableAmount() {
	suite.SetupTest()
	tests := []struct {
		name          string
		req           *types.MaxBondableAmountRequest
		setup         func(ctx sdk.Context, k keeper.Keeper)
		expectedError error
		expectedResp  *types.MaxBondableAmountResponse
	}{
		{
			name: "valid request",
			req:  &types.MaxBondableAmountRequest{PoolId: 1},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetPool(ctx, types.Pool{Id: 1, DepositDenom: ptypes.BaseCurrency})
			},
			expectedError: nil,
			expectedResp: &types.MaxBondableAmountResponse{
				Amount: math.NewInt(1000_000),
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
			req:  &types.MaxBondableAmountRequest{PoolId: 2},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
			},
			expectedError: status.Errorf(codes.NotFound, "pool %d not found", 2),
			expectedResp:  nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.setup != nil {
				tt.setup(suite.ctx, *suite.app.StablestakeKeeper)
			}

			resp, err := suite.app.StablestakeKeeper.MaxBondableAmount(suite.ctx, tt.req)
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

func (suite *KeeperTestSuite) TestPool() {
	p := types.Pool{
		InterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     math.LegacyOneDec(),
		NetAmount:            math.ZeroInt(),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		MaxWithdrawRatio:     math.LegacyMustNewDecFromStr("0.7"),

		Id:           1,
		DepositDenom: ptypes.BaseCurrency,
	}

	poolResponse := types.PoolResponse{
		RedemptionRate:       math.LegacyZeroDec(),
		InterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     math.LegacyOneDec(),
		NetAmount:            math.ZeroInt(),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		MaxWithdrawRatio:     math.LegacyMustNewDecFromStr("0.7"),
		PoolId:               1,
		DepositDenom:         ptypes.BaseCurrency,
		BorrowRatio:          math.LegacyZeroDec(),
		TotalValue:           math.LegacyZeroDec(),
		TotalBorrow:          math.ZeroInt(),
	}
	tests := []struct {
		name          string
		req           *types.QueryGetPoolRequest
		setup         func(ctx sdk.Context, k keeper.Keeper)
		expectedError error
		expectedResp  *types.QueryGetPoolResponse
	}{
		{
			name: "valid request",
			req:  &types.QueryGetPoolRequest{PoolId: 1},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetPool(ctx, p)
			},
			expectedError: nil,
			expectedResp: &types.QueryGetPoolResponse{
				Pool: poolResponse,
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
			req:  &types.QueryGetPoolRequest{PoolId: 2},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
			},
			expectedError: types.ErrPoolNotFound,
			expectedResp: &types.QueryGetPoolResponse{
				Pool: poolResponse,
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.setup != nil {
				tt.setup(suite.ctx, *suite.app.StablestakeKeeper)
			}

			resp, err := suite.app.StablestakeKeeper.Pool(suite.ctx, tt.req)
			if tt.expectedError != nil {
				require.ErrorIs(suite.T(), err, tt.expectedError)
			} else {
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), tt.expectedResp.String(), resp.String())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAllPools() {
	p := types.Pool{
		InterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     math.LegacyOneDec(),
		NetAmount:            math.ZeroInt(),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		MaxWithdrawRatio:     math.LegacyMustNewDecFromStr("0.7"),
		Id:                   1,
		DepositDenom:         ptypes.BaseCurrency,
	}

	poolResponse := types.PoolResponse{
		RedemptionRate:       math.LegacyZeroDec(),
		InterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     math.LegacyOneDec(),
		NetAmount:            math.ZeroInt(),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		MaxWithdrawRatio:     math.LegacyMustNewDecFromStr("0.7"),
		PoolId:               1,
		DepositDenom:         ptypes.BaseCurrency,
		BorrowRatio:          math.LegacyZeroDec(),
		TotalValue:           math.LegacyZeroDec(),
		TotalBorrow:          math.ZeroInt(),
	}
	tests := []struct {
		name          string
		req           *types.QueryAllPoolRequest
		setup         func(ctx sdk.Context, k keeper.Keeper)
		expectedError error
		expectedResp  *types.QueryAllPoolResponse
	}{
		{
			name: "valid request",
			req:  &types.QueryAllPoolRequest{},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetPool(ctx, p)
			},
			expectedError: nil,
			expectedResp: &types.QueryAllPoolResponse{
				Pools: []types.PoolResponse{poolResponse},
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

			resp, err := suite.app.StablestakeKeeper.Pools(suite.ctx, tt.req)
			if tt.expectedError != nil {
				require.ErrorIs(suite.T(), err, tt.expectedError)
			} else {
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), tt.expectedResp.String(), resp.String())
			}
		})
	}
}
