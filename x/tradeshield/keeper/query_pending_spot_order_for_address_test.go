package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPendingSpotOrderForAddress(t *testing.T) {
	k, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	order := types.SpotOrder{
		OrderId:      1,
		OwnerAddress: "valid_address",
		OrderType:    types.SpotOrderType_LIMITBUY,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       math.LegacyNewDec(1),
		},
		OrderAmount:      sdk.NewCoin("base", math.NewInt(1)),
		OrderTargetDenom: "quote",
		Status:           types.Status_PENDING,
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

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			// order := types.PerpetualOrder{
			// 	OrderId:            0,
			// 	OwnerAddress:       fmt.Sprintf("address%d", 1),
			// 	PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
			// 	Position:           types.PerpetualPosition_LONG,
			// 	TriggerPrice:       &types.OrderPrice{Rate: sdk.NewDec(1), BaseDenom: "base", QuoteDenom: "quote"},
			// 	Collateral:         sdk.Coin{Denom: "denom", Amount: sdk.NewInt(10)},
			// 	TradingAsset:       "asset",
			// 	Leverage:           sdk.NewDec(int64(1)),
			// 	TakeProfitPrice:    sdk.NewDec(1),
			// 	PositionId:         uint64(1),
			// 	Status:             types.Status_PENDING,
			// 	StopLossPrice:      sdk.NewDec(1),
			// }
			// _ = k.AppendPendingPerpetualOrder(ctx, order)
			_ = k.AppendPendingSpotOrder(ctx, order)

			response, err := k.PendingSpotOrderForAddress(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
