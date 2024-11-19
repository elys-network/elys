package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createNPendingSpotOrder(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.SpotOrder {
	items := make([]types.SpotOrder, n)
	for i := range items {
		items[i].OrderId = keeper.AppendPendingSpotOrder(ctx, items[i])
	}
	return items
}

func TestPendingSpotOrderGet(t *testing.T) {
	keeper, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPendingSpotOrder(ctx, item.OrderId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestPendingSpotOrderRemove(t *testing.T) {
	keeper, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePendingSpotOrder(ctx, item.OrderId)
		_, found := keeper.GetPendingSpotOrder(ctx, item.OrderId)
		require.False(t, found)
	}
}

func TestPendingSpotOrderGetAll(t *testing.T) {
	keeper, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPendingSpotOrder(ctx)),
	)
}

func TestPendingSpotOrderCount(t *testing.T) {
	keeper, ctx, _, _, _ := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetPendingSpotOrderCount(ctx)-1)
}

// TestExecuteStopLossOrder
func TestExecuteStopLossOrder(t *testing.T) {
	keeper, ctx, ammKeeper, tierKeeper, _ := keepertest.TradeshieldKeeper(t)

	address := sdk.AccAddress([]byte("address"))

	tierKeeper.On("CalculateUSDValue", ctx, "base", sdkmath.NewInt(1)).Return(sdkmath.LegacyNewDec(1))
	tierKeeper.On("CalculateUSDValue", ctx, "quote", sdkmath.NewInt(1)).Return(sdkmath.LegacyNewDec(1))
	ammKeeper.On("SwapByDenom", ctx, &ammtypes.MsgSwapByDenom{
		Sender:    address.String(),
		Amount:    sdk.NewCoin("base", sdkmath.NewInt(1)),
		MinAmount: sdk.NewCoin("quote", sdkmath.ZeroInt()),
		DenomIn:   "base",
		DenomOut:  "quote",
		Recipient: address.String(),
	}).Return(&ammtypes.MsgSwapByDenomResponse{}, nil)

	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: address.String(),
		OrderId:      0,
		OrderType:    types.SpotOrderType_STOPLOSS,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdkmath.LegacyNewDec(1),
		},
		OrderTargetDenom: "quote",
		OrderAmount:      sdk.NewCoin("base", sdkmath.NewInt(1)),
	})

	order, _ := keeper.GetPendingSpotOrder(ctx, 1)

	err := keeper.ExecuteStopLossOrder(ctx, order)
	require.NoError(t, err)

	// Should remove from pending order list
	res := keeper.GetAllPendingSpotOrder(ctx)
	assert.Equal(t, res, []types.SpotOrder(nil))

	// Should remove from main storage
	_, found := keeper.GetPendingSpotOrder(ctx, 1)
	assert.False(t, found)
}

// TestExecuteLimitSellOrder
func TestExecuteLimitSellOrder(t *testing.T) {
	keeper, ctx, ammKeeper, tierKeeper, _ := keepertest.TradeshieldKeeper(t)

	address := sdk.AccAddress([]byte("address"))

	tierKeeper.On("CalculateUSDValue", ctx, "base", sdkmath.NewInt(1)).Return(sdkmath.LegacyNewDec(1))
	tierKeeper.On("CalculateUSDValue", ctx, "quote", sdkmath.NewInt(1)).Return(sdkmath.LegacyNewDec(1))
	ammKeeper.On("SwapByDenom", ctx, &ammtypes.MsgSwapByDenom{
		Sender:    address.String(),
		Amount:    sdk.NewCoin("base", sdkmath.NewInt(1)),
		MinAmount: sdk.NewCoin("quote", sdkmath.ZeroInt()),
		DenomIn:   "base",
		DenomOut:  "quote",
		Recipient: address.String(),
	}).Return(&ammtypes.MsgSwapByDenomResponse{}, nil)

	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: address.String(),
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITSELL,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdkmath.LegacyNewDec(1),
		},
		OrderTargetDenom: "quote",
		OrderAmount:      sdk.NewCoin("base", sdkmath.NewInt(1)),
	})

	order, _ := keeper.GetPendingSpotOrder(ctx, 1)

	err := keeper.ExecuteLimitSellOrder(ctx, order)
	require.NoError(t, err)

	// Should remove from pending order list
	res := keeper.GetAllPendingSpotOrder(ctx)
	assert.Equal(t, res, []types.SpotOrder(nil))

	// Should remove from main storage
	_, found := keeper.GetPendingSpotOrder(ctx, 1)
	assert.False(t, found)
}

// TestExecuteLimitBuyOrder
func TestExecuteLimitBuyOrder(t *testing.T) {
	keeper, ctx, ammKeeper, tierKeeper, _ := keepertest.TradeshieldKeeper(t)

	address := sdk.AccAddress([]byte("address"))

	tierKeeper.On("CalculateUSDValue", ctx, "base", sdkmath.NewInt(1)).Return(sdkmath.LegacyNewDec(1))
	tierKeeper.On("CalculateUSDValue", ctx, "quote", sdkmath.NewInt(1)).Return(sdkmath.LegacyNewDec(1))
	ammKeeper.On("SwapByDenom", ctx, &ammtypes.MsgSwapByDenom{
		Sender:    address.String(),
		Amount:    sdk.NewCoin("base", sdkmath.NewInt(1)),
		MinAmount: sdk.NewCoin("quote", sdkmath.ZeroInt()),
		DenomIn:   "base",
		DenomOut:  "quote",
		Recipient: address.String(),
	}).Return(&ammtypes.MsgSwapByDenomResponse{}, nil)

	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: address.String(),
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITBUY,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdkmath.LegacyNewDec(1),
		},
		OrderTargetDenom: "quote",
		OrderAmount:      sdk.NewCoin("base", sdkmath.NewInt(1)),
	})

	order, _ := keeper.GetPendingSpotOrder(ctx, 1)

	err := keeper.ExecuteLimitBuyOrder(ctx, order)
	require.NoError(t, err)

	// Should remove from pending order list
	res := keeper.GetAllPendingSpotOrder(ctx)
	assert.Equal(t, res, []types.SpotOrder(nil))

	// Should remove from main storage
	_, found := keeper.GetPendingSpotOrder(ctx, 1)
	assert.False(t, found)
}

// TestExecuteMarketBuyOrder
func TestExecuteMarketBuyOrder(t *testing.T) {
	keeper, ctx, ammKeeper, _, _ := keepertest.TradeshieldKeeper(t)

	address := sdk.AccAddress([]byte("address"))

	ammKeeper.On("SwapByDenom", ctx, &ammtypes.MsgSwapByDenom{
		Sender:    address.String(),
		Amount:    sdk.NewCoin("base", sdkmath.NewInt(1)),
		MinAmount: sdk.NewCoin("quote", sdkmath.ZeroInt()),
		DenomIn:   "base",
		DenomOut:  "quote",
		Recipient: address.String(),
	}).Return(&ammtypes.MsgSwapByDenomResponse{}, nil)

	order := types.SpotOrder{
		OwnerAddress: address.String(),
		OrderId:      0,
		OrderType:    types.SpotOrderType_MARKETBUY,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdkmath.LegacyNewDec(1),
		},
		OrderTargetDenom: "quote",
		OrderAmount:      sdk.NewCoin("base", sdkmath.NewInt(1)),
	}

	err := keeper.ExecuteMarketBuyOrder(ctx, order)
	require.NoError(t, err)

	// Should remove from pending order list
	res := keeper.GetAllPendingSpotOrder(ctx)
	assert.Equal(t, res, []types.SpotOrder(nil))

	// Should remove from main storage
	_, found := keeper.GetPendingSpotOrder(ctx, 1)
	assert.False(t, found)
}
