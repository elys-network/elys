package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TradeshieldKeeperTestSuite) TestPendingSpotOrderForAddress() {

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

	order2 := types.SpotOrder{
		OrderId:      2,
		OwnerAddress: "valid_address",
		OrderType:    types.SpotOrderType_LIMITBUY,
		LegacyOrderPriceV1: types.LegacyOrderPriceV1{
			Rate: math.LegacyNewDec(1),
		},
		OrderPrice:       math.LegacyNewDec(1),
		OrderAmount:      sdk.NewCoin("base", math.NewInt(1)),
		OrderTargetDenom: "quote",
		Status:           types.Status_EXECUTED,
	}

	tests := []struct {
		desc     string
		request  *types.QueryPendingSpotOrderForAddressRequest
		response *types.QueryPendingSpotOrderForAddressResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryPendingSpotOrderForAddressRequest{
				Address: "valid_address",
				Status:  types.Status_ALL,
			},
			response: &types.QueryPendingSpotOrderForAddressResponse{
				PendingSpotOrders: []types.SpotOrder{order, order2},
			},
			err: nil,
		},
		{
			desc: "valid request",
			request: &types.QueryPendingSpotOrderForAddressRequest{
				Address: "valid_address",
				Status:  types.Status_PENDING,
			},
			response: &types.QueryPendingSpotOrderForAddressResponse{
				PendingSpotOrders: []types.SpotOrder{order},
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
	_ = suite.app.TradeshieldKeeper.AppendPendingSpotOrder(suite.ctx, order2)

	for _, tc := range tests {
		suite.Run(tc.desc, func() {

			response, err := suite.app.TradeshieldKeeper.PendingSpotOrderForAddress(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response, response)
			}
		})
	}
}
