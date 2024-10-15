package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPendingPerpetualOrderForAddress(t *testing.T) {
	k, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	order := types.PerpetualOrder{
		OrderId:            1,
		OwnerAddress:       "valid_address",
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		Position:           types.PerpetualPosition_LONG,
		TriggerPrice:       &types.OrderPrice{Rate: sdk.NewDec(1), BaseDenom: "base", QuoteDenom: "quote"},
		Collateral:         sdk.Coin{Denom: "denom", Amount: sdk.NewInt(10)},
		TradingAsset:       "asset",
		Leverage:           sdk.NewDec(int64(1)),
		TakeProfitPrice:    sdk.NewDec(1),
		PositionId:         uint64(1),
		Status:             types.Status_PENDING,
		StopLossPrice:      sdk.NewDec(1),
		PositionSize:       sdk.NewCoin("denom", sdk.NewInt(10)),
		LiquidationPrice:   sdk.NewDec(1),
		FundingRate:        sdk.NewDec(1),
		BorrowInterestRate: sdk.NewDec(1),
	}

	order2 := types.PerpetualOrder{
		OrderId:            2,
		OwnerAddress:       "valid_address",
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		Position:           types.PerpetualPosition_LONG,
		TriggerPrice:       &types.OrderPrice{Rate: sdk.NewDec(2), BaseDenom: "base", QuoteDenom: "quote"},
		Collateral:         sdk.Coin{Denom: "denom", Amount: sdk.NewInt(10)},
		TradingAsset:       "asset",
		Leverage:           sdk.NewDec(int64(1)),
		TakeProfitPrice:    sdk.NewDec(1),
		PositionId:         uint64(1),
		Status:             types.Status_EXECUTED,
		StopLossPrice:      sdk.NewDec(1),
		PositionSize:       sdk.NewCoin("denom", sdk.NewInt(10)),
		LiquidationPrice:   sdk.NewDec(1),
		FundingRate:        sdk.NewDec(1),
		BorrowInterestRate: sdk.NewDec(1),
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
				Address: "valid_address",
				Status:  types.Status_ALL,
			},
			response: &types.QueryPendingPerpetualOrderForAddressResponse{
				PendingPerpetualOrders: []types.PerpetualOrder{order, order2},
			},
			err: nil,
		},
		{
			desc: "valid request",
			request: &types.QueryPendingPerpetualOrderForAddressRequest{
				Address: "valid_address",
				Status:  types.Status_EXECUTED,
			},
			response: &types.QueryPendingPerpetualOrderForAddressResponse{
				PendingPerpetualOrders: []types.PerpetualOrder{order2},
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	_ = k.AppendPendingPerpetualOrder(ctx, order)
	_ = k.AppendPendingPerpetualOrder(ctx, order2)

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.PendingPerpetualOrderForAddress(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
