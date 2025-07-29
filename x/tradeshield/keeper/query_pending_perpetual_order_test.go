package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TradeshieldKeeperTestSuite) TestPendingPerpetualtOrder() {

	order := types.PerpetualOrder{
		OwnerAddress:       sdk.AccAddress([]byte("valid_address")).String(),
		OrderId:            1,
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice:       math.LegacyMustNewDecFromStr("10"),
		Position:           types.PerpetualPosition_LONG,
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(10),
		TakeProfitPrice:    math.LegacyNewDec(10),
		StopLossPrice:      math.LegacyZeroDec(),
		PoolId:             1,
	}

	tests := []struct {
		desc     string
		request  *types.QueryGetPendingPerpetualOrderRequest
		response *types.QueryGetPendingPerpetualOrderResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryGetPendingPerpetualOrderRequest{
				OwnerAddress: sdk.AccAddress([]byte("valid_address")).String(),
				PoolId:       1,
				OrderId:      1,
			},
			response: &types.QueryGetPendingPerpetualOrderResponse{
				PendingPerpetualOrder: types.PerpetualOrderExtraInfo{
					PerpetualOrder:     &order,
					LiquidationPrice:   math.LegacyZeroDec(),
					FundingRate:        math.LegacyZeroDec(),
					BorrowInterestRate: math.LegacyZeroDec(),
				},
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	_ = suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, order)

	for _, tc := range tests {
		suite.Run(tc.desc, func() {

			response, err := suite.app.TradeshieldKeeper.PendingPerpetualOrder(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response, response)
			}
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) TestPendingPerpetualOrderAll() {

	order := types.PerpetualOrder{
		OwnerAddress:       sdk.AccAddress([]byte("valid_address")).String(),
		OrderId:            1,
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice:       math.LegacyMustNewDecFromStr("10"),
		Position:           types.PerpetualPosition_LONG,
		Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
		Leverage:           math.LegacyNewDec(10),
		TakeProfitPrice:    math.LegacyNewDec(10),
		StopLossPrice:      math.LegacyZeroDec(),
		PoolId:             1,
	}

	order2 := order
	order2.OrderId = 2

	tests := []struct {
		desc     string
		request  *types.QueryAllPendingPerpetualOrderRequest
		response *types.QueryAllPendingPerpetualOrderResponse
		err      error
		setup    func()
	}{
		{
			desc: "valid request, one order",
			request: &types.QueryAllPendingPerpetualOrderRequest{
				Pagination: &query.PageRequest{},
			},
			response: &types.QueryAllPendingPerpetualOrderResponse{
				PendingPerpetualOrder: []types.PerpetualOrderExtraInfo{
					{
						PerpetualOrder:     &order,
						LiquidationPrice:   math.LegacyZeroDec(),
						FundingRate:        math.LegacyZeroDec(),
						BorrowInterestRate: math.LegacyZeroDec(),
					},
				},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
			err: nil,
			setup: func() {
				_ = suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, order)
			},
		},
		{
			desc:    "valid request, multiple orders",
			request: &types.QueryAllPendingPerpetualOrderRequest{},
			response: &types.QueryAllPendingPerpetualOrderResponse{
				PendingPerpetualOrder: []types.PerpetualOrderExtraInfo{
					{
						PerpetualOrder:     &order,
						LiquidationPrice:   math.LegacyZeroDec(),
						FundingRate:        math.LegacyZeroDec(),
						BorrowInterestRate: math.LegacyZeroDec(),
					},
					{
						PerpetualOrder:     &order2,
						LiquidationPrice:   math.LegacyZeroDec(),
						FundingRate:        math.LegacyZeroDec(),
						BorrowInterestRate: math.LegacyZeroDec(),
					},
				},
				Pagination: &query.PageResponse{
					Total: 2,
				},
			},
			err: nil,
			setup: func() {
				_ = suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, order2)
			},
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
			setup:   func() {},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			tc.setup()
			response, err := suite.app.TradeshieldKeeper.PendingPerpetualOrderAll(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response, response)
			}
		})
	}
}
