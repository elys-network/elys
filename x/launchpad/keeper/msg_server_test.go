package keeper_test

import (
	"context"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/launchpad/keeper"
	"github.com/elys-network/elys/x/launchpad/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.LaunchpadKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestMsgServerBuyElys(t *testing.T) {
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
	cacheCtx, _ := ctx.CacheContext()
	response, err := msgServer.BuyElys(cacheCtx, &types.MsgBuyElys{
		Sender:        buyer.String(),
		SpendingToken: ptypes.BaseCurrency,
		TokenAmount:   sdk.NewInt(1000_000),
	})
	require.NoError(t, err)
	require.Len(t, response.OrderIds, 1)
	require.Equal(t, response.OrderIds[0], uint64(1))
	orders := k.GetAllOrders(cacheCtx)
	require.Len(t, orders, 1)
	require.Equal(t, orders[0].OrderId, uint64(1))
	require.Equal(t, orders[0].OrderMaker, buyer.String())
	require.Equal(t, orders[0].SpendingToken, ptypes.BaseCurrency)
	require.Equal(t, orders[0].TokenAmount.String(), "1000000")
	require.Equal(t, orders[0].ElysAmount.String(), "1333333")
	require.Equal(t, orders[0].ReturnedElysAmount.String(), "0")
	require.Equal(t, orders[0].BonusAmount.String(), "1333333")
	require.Equal(t, orders[0].VestingStarted, false)

	cacheCtx, _ = ctx.CacheContext()
	response, err = msgServer.BuyElys(cacheCtx, &types.MsgBuyElys{
		Sender:        buyer.String(),
		SpendingToken: ptypes.BaseCurrency,
		TokenAmount:   sdk.NewInt(1000_000_000_000),
	})
	require.NoError(t, err)
	require.Len(t, response.OrderIds, 2)
	require.Equal(t, response.OrderIds[0], uint64(1))
	require.Equal(t, response.OrderIds[1], uint64(2))
	orders = k.GetAllOrders(cacheCtx)
	require.Len(t, orders, 2)
	require.Equal(t, orders[0].OrderId, uint64(1))
	require.Equal(t, orders[0].OrderMaker, buyer.String())
	require.Equal(t, orders[0].SpendingToken, ptypes.BaseCurrency)
	require.Equal(t, orders[0].TokenAmount.String(), "675000000000")
	require.Equal(t, orders[0].ElysAmount.String(), "900000000000")
	require.Equal(t, orders[0].ReturnedElysAmount.String(), "0")
	require.Equal(t, orders[0].BonusAmount.String(), "900000000000")
	require.Equal(t, orders[0].VestingStarted, false)
	require.Equal(t, orders[1].OrderId, uint64(2))
	require.Equal(t, orders[1].OrderMaker, buyer.String())
	require.Equal(t, orders[1].SpendingToken, ptypes.BaseCurrency)
	require.Equal(t, orders[1].TokenAmount.String(), "325000000000")
	require.Equal(t, orders[1].ElysAmount.String(), "433333333333")
	require.Equal(t, orders[1].ReturnedElysAmount.String(), "0")
	require.Equal(t, orders[1].BonusAmount.String(), "389999999999")
	require.Equal(t, orders[1].VestingStarted, false)
}

func TestMsgServerReturnElys(t *testing.T) {
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

	buyerOldBalance := app.BankKeeper.GetAllBalances(ctx, buyer)

	// Try returning 25%
	msgServer = keeper.NewMsgServerImpl(k)
	cacheCtx, _ := ctx.CacheContext()
	cacheCtx = cacheCtx.WithBlockTime(now.Add(time.Second * time.Duration(params.LaunchpadDuration+1)))
	_, err = msgServer.ReturnElys(cacheCtx, &types.MsgReturnElys{
		Sender:           buyer.String(),
		OrderId:          1,
		ReturnElysAmount: sdk.NewInt(333333),
	})
	require.NoError(t, err)

	buyerNewBalance := app.BankKeeper.GetAllBalances(cacheCtx, buyer)
	require.Equal(t, buyerOldBalance.AmountOf(ptypes.Elys).String(), buyerNewBalance.AmountOf(ptypes.Elys).Add(sdk.NewInt(333333)).String())
	require.Equal(t, buyerNewBalance.AmountOf(ptypes.BaseCurrency).String(), buyerOldBalance.AmountOf(ptypes.BaseCurrency).Add(sdk.NewInt(249999)).String())

	// check orders change
	orders = k.GetAllOrders(cacheCtx)
	require.Len(t, orders, 1)
	require.Equal(t, orders[0].OrderId, uint64(1))
	require.Equal(t, orders[0].OrderMaker, buyer.String())
	require.Equal(t, orders[0].SpendingToken, ptypes.BaseCurrency)
	require.Equal(t, orders[0].TokenAmount.String(), "1000000")
	require.Equal(t, orders[0].ElysAmount.String(), "1333333")
	require.Equal(t, orders[0].ReturnedElysAmount.String(), "333333")

	// Try returning more than half including previous one
	_, err = msgServer.ReturnElys(cacheCtx, &types.MsgReturnElys{
		Sender:           buyer.String(),
		OrderId:          1,
		ReturnElysAmount: sdk.NewInt(433333),
	})
	require.Error(t, err)
}

func TestIsEnabledToken(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper

	params := k.GetParams(ctx)
	require.Len(t, params.SpendingTokens, 3)
	require.Equal(t, params.SpendingTokens[0], ptypes.BaseCurrency)

	isEnabled := k.IsEnabledToken(ctx, ptypes.BaseCurrency)
	require.Equal(t, isEnabled, true)

	isEnabled = k.IsEnabledToken(ctx, ptypes.Elys)
	require.Equal(t, isEnabled, false)

	isEnabled = k.IsEnabledToken(ctx, "uatom")
	require.Equal(t, isEnabled, true)
}

func TestMsgServerDepositElysToken(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper
	SetupStableCoinPrices(app, ctx)

	admin := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	coins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000_000_000_000_000)}
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, admin, coins)
	require.NoError(t, err)

	msgServer := keeper.NewMsgServerImpl(k)
	cacheCtx, _ := ctx.CacheContext()
	_, err = msgServer.DepositElysToken(cacheCtx, &types.MsgDepositElysToken{
		Sender: admin.String(),
		Coin:   sdk.NewInt64Coin(ptypes.Elys, 1000_000),
	})
	require.NoError(t, err)

	resp, err := k.ModuleBalances(cacheCtx, &types.QueryModuleBalancesRequest{})
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(resp.Coins).String(), sdk.NewInt64Coin(ptypes.Elys, 1000_000).String())
}

func TestMsgServerWithdrawRaised(t *testing.T) {
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
	admin := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	coins = sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1100_000_000_000)}
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, buyer, coins)
	require.NoError(t, err)

	now := time.Now()
	ctx = ctx.WithBlockTime(now)
	params := k.GetParams(ctx)
	params.LaunchpadStarttime = uint64(now.Unix())
	params.WithdrawAddress = admin.String()
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

	adminOldBalance := app.BankKeeper.GetAllBalances(ctx, admin)

	// Try withdrawing 10% with invalid account
	msgServer = keeper.NewMsgServerImpl(k)
	cacheCtx, _ := ctx.CacheContext()
	cacheCtx = cacheCtx.WithBlockTime(now.Add(time.Second * time.Duration(params.LaunchpadDuration+1)))
	_, err = msgServer.WithdrawRaised(cacheCtx, &types.MsgWithdrawRaised{
		Sender: buyer.String(),
		Coins:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000)},
	})
	require.Error(t, err)

	// Try withdrawing 10% with admin account
	msgServer = keeper.NewMsgServerImpl(k)
	cacheCtx, _ = ctx.CacheContext()
	cacheCtx = cacheCtx.WithBlockTime(now.Add(time.Second * time.Duration(params.LaunchpadDuration+1)))
	_, err = msgServer.WithdrawRaised(cacheCtx, &types.MsgWithdrawRaised{
		Sender: admin.String(),
		Coins:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000)},
	})
	require.NoError(t, err)

	adminNewBalance := app.BankKeeper.GetAllBalances(cacheCtx, admin)
	require.Equal(t, adminNewBalance.AmountOf(ptypes.BaseCurrency).String(), adminOldBalance.AmountOf(ptypes.BaseCurrency).Add(sdk.NewInt(100_000)).String())

	// check params withdrawn amount change
	params = k.GetParams(cacheCtx)
	require.Len(t, orders, 1)
	require.Equal(t, sdk.Coins(params.WithdrawnAmount).String(), sdk.NewInt64Coin(ptypes.BaseCurrency, 100_000).String())

	// Try withdrawing more than half including previous one
	_, err = msgServer.WithdrawRaised(cacheCtx, &types.MsgWithdrawRaised{
		Sender: admin.String(),
		Coins:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000)},
	})
	require.Error(t, err)

	// Withdraw full after return period is over
	cacheCtx = cacheCtx.WithBlockTime(now.Add(time.Second * time.Duration(params.LaunchpadDuration+params.ReturnDuration+1)))
	_, err = msgServer.WithdrawRaised(cacheCtx, &types.MsgWithdrawRaised{
		Sender: admin.String(),
		Coins:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 500_000)},
	})
	require.NoError(t, err)

	adminNewBalance = app.BankKeeper.GetAllBalances(cacheCtx, admin)
	require.Equal(t, adminNewBalance.AmountOf(ptypes.BaseCurrency).String(), adminOldBalance.AmountOf(ptypes.BaseCurrency).Add(sdk.NewInt(600_000)).String())
}
