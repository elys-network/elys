package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/stretchr/testify/require"
)

func createNPendingPerpetualOrder(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PerpetualOrder {
	items := make([]types.PerpetualOrder, n)
	for i := range items {
		items[i] = types.PerpetualOrder{
			OrderId:            0,
			OwnerAddress:       fmt.Sprintf("address%d", i),
			PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
			Position:           types.PerpetualPosition_LONG,
			TriggerPrice:       types.TriggerPrice{Rate: sdkmath.LegacyNewDec(1), TradingAssetDenom: "base"},
			Collateral:         sdk.Coin{Denom: "denom", Amount: math.NewInt(10)},
			TradingAsset:       "asset",
			Leverage:           sdkmath.LegacyNewDec(int64(i)),
			TakeProfitPrice:    sdkmath.LegacyNewDec(1),
			PositionId:         uint64(i),
			Status:             types.Status_PENDING,
			StopLossPrice:      math.LegacyNewDec(1),
		}
		items[i].OrderId = keeper.AppendPendingPerpetualOrder(ctx, items[i])
	}
	return items
}

func TestPendingPerpetualOrderGet(t *testing.T) {
	keeper, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPendingPerpetualOrder(ctx, item.OrderId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestPendingPerpetualOrderRemove(t *testing.T) {
	keeper, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePendingPerpetualOrder(ctx, item.OrderId)
		_, found := keeper.GetPendingPerpetualOrder(ctx, item.OrderId)
		require.False(t, found)
	}
}

func TestPendingPerpetualOrderGetAll(t *testing.T) {
	keeper, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPendingPerpetualOrder(ctx)),
	)
}

func TestPendingPerpetualOrderCount(t *testing.T) {
	keeper, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetPendingPerpetualOrderCount(ctx)-1)
}

// TestExecuteLimitOpenOrder
func TestExecuteLimitOpenOrder(t *testing.T) {
	keeper, ctx, _, _, perpetualKeeper := keepertest.TradeshieldKeeper(t)

	address := sdk.AccAddress([]byte("address"))
	// Set up the mock to expect the GetAssetPrice call and specify the return values
	perpetualKeeper.On("GetAssetPrice", ctx, "quote").Return(math.LegacyOneDec(), nil)

	perpetualKeeper.On("Open", ctx, &perpetualtypes.MsgOpen{
		Creator:         address.String(),
		Position:        perpetualtypes.Position_LONG,
		Leverage:        math.LegacyNewDec(10),
		TradingAsset:    "quote",
		Collateral:      sdk.Coin{Denom: "base", Amount: sdkmath.NewInt(10)},
		TakeProfitPrice: math.LegacyZeroDec(),
		StopLossPrice:   math.LegacyZeroDec(),
	}).Return(&perpetualtypes.MsgOpenResponse{Id: 1}, nil)

	keeper.AppendPendingPerpetualOrder(ctx, types.PerpetualOrder{
		OwnerAddress:       address.String(),
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice: types.TriggerPrice{
			TradingAssetDenom: "base",
			Rate:              math.LegacyMustNewDecFromStr("10"),
		},
		Position:        types.PerpetualPosition_LONG,
		Collateral:      sdk.Coin{Denom: "base", Amount: sdkmath.NewInt(10)},
		TradingAsset:    "quote",
		Leverage:        math.LegacyNewDec(10),
		TakeProfitPrice: math.LegacyZeroDec(),
		StopLossPrice:   math.LegacyZeroDec(),
	})

	order, _ := keeper.GetPendingPerpetualOrder(ctx, 1)

	err := keeper.ExecuteLimitOpenOrder(ctx, order)
	require.NoError(t, err)

	_, found := keeper.GetPendingPerpetualOrder(ctx, 1)
	require.False(t, found)
}

// TestExecuteLimitCloseOrder
func TestExecuteLimitCloseOrder(t *testing.T) {
	keeper, ctx, _, _, perpetualKeeper := keepertest.TradeshieldKeeper(t)

	address := sdk.AccAddress([]byte("address"))

	// Set up the mock to expect the GetAssetPrice call and specify the return values
	perpetualKeeper.On("GetAssetPrice", ctx, "quote").Return(math.LegacyOneDec(), nil)
	perpetualKeeper.On("Close", ctx, &perpetualtypes.MsgClose{
		Creator: address.String(),
		Id:      1,
		Amount:  sdkmath.ZeroInt(),
	}).Return(&perpetualtypes.MsgCloseResponse{Id: 1, Amount: sdkmath.ZeroInt()}, nil)

	keeper.AppendPendingPerpetualOrder(ctx, types.PerpetualOrder{
		OwnerAddress:       address.String(),
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		TriggerPrice: types.TriggerPrice{
			TradingAssetDenom: "base",
			Rate:              math.LegacyNewDec(0),
		},
		TradingAsset: "quote",
		Position:     types.PerpetualPosition_LONG,
		PositionId:   1,
	})

	order, _ := keeper.GetPendingPerpetualOrder(ctx, 1)

	err := keeper.ExecuteLimitCloseOrder(ctx, order)
	require.NoError(t, err)

	_, found := keeper.GetPendingPerpetualOrder(ctx, 1)
	require.False(t, found)
}

// TestExecuteMarketOpenOrder
func TestExecuteMarketOpenOrder(t *testing.T) {
	keeper, ctx, _, _, perpetualKeeper := keepertest.TradeshieldKeeper(t)

	address := sdk.AccAddress([]byte("address"))

	perpetualKeeper.On("Open", ctx, &perpetualtypes.MsgOpen{
		Creator:         address.String(),
		Position:        perpetualtypes.Position_LONG,
		Leverage:        math.LegacyNewDec(10),
		TradingAsset:    "quote",
		Collateral:      sdk.Coin{Denom: "base", Amount: sdkmath.NewInt(10)},
		TakeProfitPrice: math.LegacyZeroDec(),
		StopLossPrice:   math.LegacyZeroDec(),
	}).Return(&perpetualtypes.MsgOpenResponse{Id: 1}, nil)

	keeper.AppendPendingPerpetualOrder(ctx, types.PerpetualOrder{
		OwnerAddress:       address.String(),
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice: types.TriggerPrice{
			TradingAssetDenom: "base",
			Rate:              math.LegacyNewDec(1),
		},
		Position:        types.PerpetualPosition_LONG,
		Collateral:      sdk.Coin{Denom: "base", Amount: sdkmath.NewInt(10)},
		TradingAsset:    "quote",
		Leverage:        math.LegacyNewDec(10),
		TakeProfitPrice: math.LegacyZeroDec(),
		StopLossPrice:   math.LegacyZeroDec(),
	})

	order, _ := keeper.GetPendingPerpetualOrder(ctx, 1)

	err := keeper.ExecuteMarketOpenOrder(ctx, order)
	require.NoError(t, err)

	_, found := keeper.GetPendingPerpetualOrder(ctx, 1)
	require.False(t, found)
}

// TestExecuteMarketCloseOrder
func TestExecuteMarketCloseOrder(t *testing.T) {
	keeper, ctx, _, _, perpetualKeeper := keepertest.TradeshieldKeeper(t)

	address := sdk.AccAddress([]byte("address"))

	perpetualKeeper.On("Close", ctx, &perpetualtypes.MsgClose{
		Creator: address.String(),
		Id:      1,
		Amount:  sdkmath.ZeroInt(),
	}).Return(&perpetualtypes.MsgCloseResponse{Id: 1, Amount: sdkmath.ZeroInt()}, nil)

	keeper.AppendPendingPerpetualOrder(ctx, types.PerpetualOrder{
		OwnerAddress:       address.String(),
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		TriggerPrice: types.TriggerPrice{
			TradingAssetDenom: "base",
			Rate:              math.LegacyNewDec(1),
		},
		Position:   types.PerpetualPosition_LONG,
		PositionId: 1,
	})

	order, _ := keeper.GetPendingPerpetualOrder(ctx, 1)

	err := keeper.ExecuteMarketCloseOrder(ctx, order)
	require.NoError(t, err)

	_, found := keeper.GetPendingPerpetualOrder(ctx, 1)
	require.False(t, found)
}
