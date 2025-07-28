package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/v7/x/stablestake/keeper"
	"github.com/elys-network/elys/v7/x/stablestake/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *KeeperTestSuite) TestBorrowRatio() {
	tests := []struct {
		name          string
		req           *types.QueryBorrowRatioRequest
		setup         func(ctx sdk.Context, k keeper.Keeper)
		expectedError error
		expectedResp  *types.QueryBorrowRatioResponse
	}{
		{
			name: "valid request",
			req: &types.QueryBorrowRatioRequest{
				PoolId: 1,
			},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				pool := types.Pool{
					NetAmount:    sdkmath.NewInt(1000),
					DepositDenom: "token",
					Id:           1,
				}
				k.SetPool(ctx, pool)
				// bootstrap balances
				err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("token", sdkmath.NewInt(500))))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin("token", sdkmath.NewInt(500))))
				suite.Require().NoError(err)
			},
			expectedError: nil,
			expectedResp: &types.QueryBorrowRatioResponse{
				NetAmount:   sdkmath.NewInt(1000),
				TotalBorrow: sdkmath.NewInt(500),
				BorrowRatio: sdkmath.LegacyNewDecWithPrec(5, 1), // 0.5
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
			name: "zero total value",
			req: &types.QueryBorrowRatioRequest{
				PoolId: 1,
			},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				pool := types.Pool{
					NetAmount:    sdkmath.ZeroInt(),
					DepositDenom: "token",
					Id:           1,
				}
				k.SetPool(ctx, pool)
			},
			expectedError: nil,
			expectedResp: &types.QueryBorrowRatioResponse{
				NetAmount:   sdkmath.ZeroInt(),
				TotalBorrow: sdkmath.NewInt(-500),
				BorrowRatio: sdkmath.LegacyZeroDec(),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.setup != nil {
				tt.setup(suite.ctx, *suite.app.StablestakeKeeper)
			}

			resp, err := suite.app.StablestakeKeeper.BorrowRatio(suite.ctx, tt.req)
			if tt.expectedError != nil {
				require.ErrorIs(suite.T(), err, tt.expectedError)
			} else {
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), tt.expectedResp, resp)
			}
		})
	}
}
