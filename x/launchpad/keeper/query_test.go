package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/launchpad/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestBonus(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper

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

	response, err := k.Bonus(ctx, &types.QueryBonusRequest{User: "maker"})
	require.NoError(t, err)
	require.Equal(t, response.TotalBonus.String(), sdk.NewInt(1000_000).String())
}

func TestQueryOrders(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper

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

	response, err := k.Orders(ctx, &types.QueryOrdersRequest{User: "maker"})
	require.NoError(t, err)
	require.Len(t, response.Purchases, 2)
}

func TestQueryModuleBalances(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper

	coins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000_000)}

	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, coins)
	require.NoError(t, err)

	response, err := k.ModuleBalances(ctx, &types.QueryModuleBalancesRequest{})
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(response.Coins).String(), coins.String())
}

// TODO:
// func (k Keeper) BuyElysEst(goCtx context.Context, req *types.QueryBuyElysEstRequest) (*types.QueryBuyElysEstResponse, error)
// func (k Keeper) ReturnElysEst(goCtx context.Context, req *types.QueryReturnElysEstRequest) (*types.QueryReturnElysEstResponse, error)
