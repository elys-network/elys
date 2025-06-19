package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TradeshieldKeeperTestSuite) TestPendingSpotOrder() {

	order := types.SpotOrder{
		OrderId:      1,
		OwnerAddress: "valid_address",
		OrderType:    types.SpotOrderType_LIMITBUY,
		LegacyOrderPriceV1: types.LegacyOrderPriceV1{
			Rate: math.LegacyNewDec(1),
		},
		OrderPrice:       math.LegacyNewDec(1),
		OrderAmount:      sdk.NewCoin("base", math.NewInt(1)),
		OrderTargetDenom: "quote",
		Status:           types.Status_PENDING,
	}

	tests := []struct {
		desc     string
		request  *types.QueryGetPendingSpotOrderRequest
		response *types.QueryGetPendingSpotOrderResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryGetPendingSpotOrderRequest{
				Id: 1,
			},
			response: &types.QueryGetPendingSpotOrderResponse{
				PendingSpotOrder: order,
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	_ = suite.app.TradeshieldKeeper.AppendPendingSpotOrder(suite.ctx, order)

	for _, tc := range tests {
		suite.Run(tc.desc, func() {

			response, err := suite.app.TradeshieldKeeper.PendingSpotOrder(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response, response)
			}
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) TestPendingSpotOrderAll() {

	order := types.SpotOrder{
		OrderId:      1,
		OwnerAddress: "valid_address",
		OrderType:    types.SpotOrderType_LIMITBUY,
		LegacyOrderPriceV1: types.LegacyOrderPriceV1{
			Rate: math.LegacyNewDec(1),
		},
		OrderPrice:       math.LegacyNewDec(1),
		OrderAmount:      sdk.NewCoin("base", math.NewInt(1)),
		OrderTargetDenom: "quote",
		Status:           types.Status_PENDING,
	}

	order2 := order
	order2.OrderId = 2

	tests := []struct {
		desc     string
		request  *types.QueryAllPendingSpotOrderRequest
		response *types.QueryAllPendingSpotOrderResponse
		err      error
		setup    func()
	}{
		{
			desc: "valid request, one order",
			request: &types.QueryAllPendingSpotOrderRequest{
				Pagination: &query.PageRequest{},
			},
			response: &types.QueryAllPendingSpotOrderResponse{
				PendingSpotOrder: []types.SpotOrder{order},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
			err: nil,
			setup: func() {
				_ = suite.app.TradeshieldKeeper.AppendPendingSpotOrder(suite.ctx, order)
			},
		},
		{
			desc:    "valid request, multiple orders",
			request: &types.QueryAllPendingSpotOrderRequest{},
			response: &types.QueryAllPendingSpotOrderResponse{
				PendingSpotOrder: []types.SpotOrder{order, order2},
				Pagination: &query.PageResponse{
					Total: 2,
				},
			},
			err: nil,
			setup: func() {
				_ = suite.app.TradeshieldKeeper.AppendPendingSpotOrder(suite.ctx, order2)
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
			response, err := suite.app.TradeshieldKeeper.PendingSpotOrderAll(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response, response)
			}
		})
	}
}
