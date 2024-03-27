package keeper_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/launchpad/keeper"
	"github.com/elys-network/elys/x/launchpad/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func SetupStableCoinPrices(app *simapp.ElysApp, ctx sdk.Context) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	app.OracleKeeper.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.BaseCurrency,
		Display: "USDC",
		Decimal: 6,
	})
	app.OracleKeeper.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdk.NewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
}

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

func TestBuyElysEst(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper
	SetupStableCoinPrices(app, ctx)

	coins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000_000)}

	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, coins)
	require.NoError(t, err)

	response, err := k.BuyElysEst(ctx, &types.QueryBuyElysEstRequest{
		SpendingToken: ptypes.BaseCurrency,
		TokenAmount:   sdk.NewInt(1000_000),
	})
	require.NoError(t, err)
	require.Equal(t, response.BonusAmount.String(), "1333333")
	require.Equal(t, response.ElysAmount.String(), "1333333")
	require.Len(t, response.Orders, 1)
	require.Equal(t, response.Orders[0].OrderId, uint64(1))
	require.Equal(t, response.Orders[0].OrderMaker, "")
	require.Equal(t, response.Orders[0].SpendingToken, ptypes.BaseCurrency)
	require.Equal(t, response.Orders[0].TokenAmount.String(), "1000000")
	require.Equal(t, response.Orders[0].ElysAmount.String(), "1333333")
	require.Equal(t, response.Orders[0].ReturnedElysAmount.String(), "0")
	require.Equal(t, response.Orders[0].BonusAmount.String(), "1333333")
	require.Equal(t, response.Orders[0].VestingStarted, false)

	response, err = k.BuyElysEst(ctx, &types.QueryBuyElysEstRequest{
		SpendingToken: ptypes.BaseCurrency,
		TokenAmount:   sdk.NewInt(1000_000_000_000),
	})
	require.NoError(t, err)
	require.Equal(t, response.BonusAmount.String(), "1289999999999")
	require.Equal(t, response.ElysAmount.String(), "1333333333333")
	require.Len(t, response.Orders, 2)
	require.Equal(t, response.Orders[0].OrderId, uint64(1))
	require.Equal(t, response.Orders[0].OrderMaker, "")
	require.Equal(t, response.Orders[0].SpendingToken, ptypes.BaseCurrency)
	require.Equal(t, response.Orders[0].TokenAmount.String(), "675000000000")
	require.Equal(t, response.Orders[0].ElysAmount.String(), "900000000000")
	require.Equal(t, response.Orders[0].ReturnedElysAmount.String(), "0")
	require.Equal(t, response.Orders[0].BonusAmount.String(), "900000000000")
	require.Equal(t, response.Orders[0].VestingStarted, false)
	require.Equal(t, response.Orders[1].OrderId, uint64(2))
	require.Equal(t, response.Orders[1].OrderMaker, "")
	require.Equal(t, response.Orders[1].SpendingToken, ptypes.BaseCurrency)
	require.Equal(t, response.Orders[1].TokenAmount.String(), "325000000000")
	require.Equal(t, response.Orders[1].ElysAmount.String(), "433333333333")
	require.Equal(t, response.Orders[1].ReturnedElysAmount.String(), "0")
	require.Equal(t, response.Orders[1].BonusAmount.String(), "389999999999")
	require.Equal(t, response.Orders[1].VestingStarted, false)
}

func TestReturnElysEst(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper
	SetupStableCoinPrices(app, ctx)

	coins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000_000_000_000_000)}
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, coins)
	require.NoError(t, err)

	buyer := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	coins = sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1100_000_000_000)}
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, buyer, coins)
	require.NoError(t, err)

	now := time.Now()
	ctx = ctx.WithBlockTime(now)
	params := k.GetParams(ctx)
	params.LaunchpadStarttime = uint64(now.Unix())
	k.SetParams(ctx, params)

	msgServer := keeper.NewMsgServerImpl(k)
	response, err := msgServer.BuyElys(ctx, &types.MsgBuyElys{
		Sender:        buyer.String(),
		SpendingToken: ptypes.BaseCurrency,
		TokenAmount:   sdk.NewInt(1000_000),
	})
	require.NoError(t, err)
	require.Len(t, response.OrderIds, 1)
	require.Equal(t, response.OrderIds[0], uint64(1))
	orders := k.GetAllOrders(ctx)
	require.Len(t, orders, 1)
	require.Equal(t, orders[0].OrderId, uint64(1))
	require.Equal(t, orders[0].OrderMaker, buyer.String())
	require.Equal(t, orders[0].SpendingToken, ptypes.BaseCurrency)
	require.Equal(t, orders[0].TokenAmount.String(), "1000000")
	require.Equal(t, orders[0].ElysAmount.String(), "1333333")
	require.Equal(t, orders[0].ReturnedElysAmount.String(), "0")

	resp, err := k.ReturnElysEst(ctx, &types.QueryReturnElysEstRequest{
		OrderId:    1,
		ElysAmount: sdk.NewInt(333333),
	})
	require.NoError(t, err)
	require.Equal(t, resp.Amount.String(), "249999")
}
