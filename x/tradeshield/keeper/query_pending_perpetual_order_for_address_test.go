package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPendingPerpetualOrderForAddress(t *testing.T) {
	k, ctx, _, _, perpetualKeeper := keepertest.TradeshieldKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	ownerAddress := "valid_address"
	ownerAddressAcc := sdk.AccAddress(ownerAddress)
	positionId := uint64(1)
	poolId := uint64(123)

	perpetualKeeper.On("GetMTP", ctx, ownerAddressAcc, positionId).Return(perpetualtypes.MTP{
		AmmPoolId: poolId,
	}, nil)
	perpetualKeeper.On("GetPool", ctx, poolId).Return(perpetualtypes.Pool{
		FundingRate:        math.LegacyNewDec(1),
		BorrowInterestRate: math.LegacyNewDec(1),
	}, true)
	perpetualKeeper.On("HandleCloseEstimation", ctx, &perpetualtypes.QueryCloseEstimationRequest{
		Address:    ownerAddress,
		PositionId: positionId,
	}).Return(&perpetualtypes.QueryCloseEstimationResponse{
		Position:         perpetualtypes.Position_LONG,
		PositionSize:     sdk.NewCoin("denom", math.NewInt(10)),
		PriceImpact:      math.LegacyNewDec(1),
		LiquidationPrice: math.LegacyNewDec(1),
	}, nil)

	order := types.PerpetualOrderExtraInfo{
		PerpetualOrder: &types.PerpetualOrder{
			OrderId:            1,
			OwnerAddress:       "valid_address",
			PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
			Position:           types.PerpetualPosition_LONG,
			TriggerPrice:       &types.TriggerPrice{Rate: math.LegacyNewDec(1), TradingAssetDenom: "base"},
			Collateral:         sdk.Coin{Denom: "denom", Amount: math.NewInt(10)},
			TradingAsset:       "asset",
			Leverage:           math.LegacyNewDec(int64(1)),
			TakeProfitPrice:    math.LegacyNewDec(1),
			PositionId:         positionId,
			Status:             types.Status_PENDING,
			StopLossPrice:      math.LegacyNewDec(1),
		},
		PositionSize:       sdk.NewCoin("denom", math.NewInt(10)),
		LiquidationPrice:   math.LegacyNewDec(1),
		FundingRate:        math.LegacyNewDec(1),
		BorrowInterestRate: math.LegacyNewDec(1),
	}

	order2 := types.PerpetualOrderExtraInfo{
		PerpetualOrder: &types.PerpetualOrder{
			OrderId:            2,
			OwnerAddress:       "valid_address",
			PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
			Position:           types.PerpetualPosition_LONG,
			TriggerPrice:       &types.TriggerPrice{Rate: math.LegacyNewDec(2), TradingAssetDenom: "base"},
			Collateral:         sdk.Coin{Denom: "denom", Amount: math.NewInt(10)},
			TradingAsset:       "asset",
			Leverage:           math.LegacyNewDec(int64(1)),
			TakeProfitPrice:    math.LegacyNewDec(1),
			PositionId:         positionId,
			Status:             types.Status_EXECUTED,
			StopLossPrice:      math.LegacyNewDec(1),
		},
		PositionSize:       sdk.NewCoin("denom", math.NewInt(10)),
		LiquidationPrice:   math.LegacyNewDec(1),
		FundingRate:        math.LegacyNewDec(1),
		BorrowInterestRate: math.LegacyNewDec(1),
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
				PendingPerpetualOrders: []types.PerpetualOrderExtraInfo{order, order2},
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

	_ = k.AppendPendingPerpetualOrder(ctx, *order.PerpetualOrder)
	_ = k.AppendPendingPerpetualOrder(ctx, *order2.PerpetualOrder)

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
