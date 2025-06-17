package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/testutil/sample"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Update test for close order once feature is complete
func (suite *TradeshieldKeeperTestSuite) TestPendingPerpetualOrderForAddress() {
	// _, _, _ = suite.SetPerpetualPool(1)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
	})
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.BaseCurrency,
		Denom:       ptypes.BaseCurrency,
		Decimals:    6,
		DisplayName: "USDC",
	})

	positionId := uint64(0)
	ownerAddress := sample.AccAddress()

	order := types.PerpetualOrderExtraInfo{
		PerpetualOrder: &types.PerpetualOrder{
			OrderId:            1,
			OwnerAddress:       ownerAddress,
			PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
			Position:           types.PerpetualPosition_LONG,
			TriggerPrice:       math.LegacyNewDec(1),
			Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
			Leverage:           math.LegacyNewDec(int64(2)),
			TakeProfitPrice:    math.LegacyNewDec(10),
			PositionId:         positionId,
			Status:             types.Status_PENDING,
			StopLossPrice:      math.LegacyNewDec(1),
			PoolId:             1,
		},
		LiquidationPrice:   math.LegacyZeroDec(),
		FundingRate:        math.LegacyZeroDec(),
		BorrowInterestRate: math.LegacyZeroDec(),
	}

	order2 := types.PerpetualOrderExtraInfo{
		PerpetualOrder: &types.PerpetualOrder{
			OrderId:            2,
			OwnerAddress:       "dummy_address",
			PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
			Position:           types.PerpetualPosition_LONG,
			TriggerPrice:       math.LegacyNewDec(2),
			Collateral:         sdk.Coin{Denom: "uatom", Amount: math.NewInt(10)},
			Leverage:           math.LegacyNewDec(int64(2)),
			TakeProfitPrice:    math.LegacyNewDec(10),
			PositionId:         positionId,
			Status:             types.Status_EXECUTED,
			StopLossPrice:      math.LegacyNewDec(1),
		},
		LiquidationPrice:   math.LegacyZeroDec(),
		FundingRate:        math.LegacyZeroDec(),
		BorrowInterestRate: math.LegacyZeroDec(),
	}

	tests := []struct {
		desc     string
		request  *types.QueryPendingPerpetualOrderForAddressRequest
		response *types.QueryPendingPerpetualOrderForAddressResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryPendingPerpetualOrderForAddressRequest{
				Address: ownerAddress,
				Status:  types.Status_ALL,
			},
			response: &types.QueryPendingPerpetualOrderForAddressResponse{
				PendingPerpetualOrders: []types.PerpetualOrderExtraInfo{order},
			},
			err: nil,
		},
		{
			desc: "valid request",
			request: &types.QueryPendingPerpetualOrderForAddressRequest{
				Address: "dummy_address",
				Status:  types.Status_EXECUTED,
			},
			response: &types.QueryPendingPerpetualOrderForAddressResponse{
				PendingPerpetualOrders: []types.PerpetualOrderExtraInfo{order2},
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	_ = suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, *order.PerpetualOrder)
	_ = suite.app.TradeshieldKeeper.AppendPendingPerpetualOrder(suite.ctx, *order2.PerpetualOrder)

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			response, err := suite.app.TradeshieldKeeper.PendingPerpetualOrderForAddress(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response, response)
			}
		})
	}
}
