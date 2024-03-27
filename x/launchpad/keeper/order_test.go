package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/launchpad/types"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper

	// check query when empty
	orderQuery := k.GetOrder(ctx, 1)
	require.Equal(t, orderQuery.OrderId, uint64(0))

	ordersQuery := k.GetAllOrders(ctx)
	require.Len(t, ordersQuery, 0)

	lastOrderId := k.LastOrderId(ctx)
	require.Equal(t, lastOrderId, uint64(0))

	orders := []types.Purchase{
		{
			OrderId:            1,
			OrderMaker:         "maker",
			SpendingToken:      "uusdc",
			TokenAmount:        sdk.NewInt(1000_000),
			ElysAmount:         sdk.NewInt(1000_000),
			ReturnedElysAmount: sdk.NewInt(500_000),
			BonusAmount:        sdk.NewInt(500_000),
			VestingStarted:     false,
		},
		{
			OrderId:            2,
			OrderMaker:         "maker",
			SpendingToken:      "uusdc",
			TokenAmount:        sdk.NewInt(1000_000),
			ElysAmount:         sdk.NewInt(1000_000),
			ReturnedElysAmount: sdk.NewInt(500_000),
			BonusAmount:        sdk.NewInt(500_000),
			VestingStarted:     false,
		},
		{
			OrderId:            3,
			OrderMaker:         "maker2",
			SpendingToken:      "uusdc",
			TokenAmount:        sdk.NewInt(1000_000),
			ElysAmount:         sdk.NewInt(1000_000),
			ReturnedElysAmount: sdk.NewInt(500_000),
			BonusAmount:        sdk.NewInt(500_000),
			VestingStarted:     false,
		},
	}

	for _, order := range orders {
		k.SetOrder(ctx, order)
	}

	// check queries after values are set
	for _, order := range orders {
		orderQuery := k.GetOrder(ctx, order.OrderId)
		require.Equal(t, orderQuery, order)
	}

	ordersQuery = k.GetAllOrders(ctx)
	require.Len(t, ordersQuery, 3)

	// delete order
	k.DeleteOrder(ctx, orders[0])

	// check queries after deleting an order
	ordersQuery = k.GetAllOrders(ctx)
	require.Len(t, ordersQuery, 2)

	lastOrderId = k.LastOrderId(ctx)
	require.Equal(t, lastOrderId, uint64(3))
}
