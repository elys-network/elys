package keeper_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/launchpad/keeper"
	"github.com/elys-network/elys/x/launchpad/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestSetEpochInfo(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper

	k.SetEpochInfo(ctx)
	_, found := app.EpochsKeeper.GetEpochInfo(ctx, "day")
	require.True(t, found)
}

func TestSetElysVestingInfo(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	k := app.LaunchpadKeeper

	k.SetElysVestingInfo(ctx)
	params := k.GetParams(ctx)
	vestingInfo, _ := app.CommitmentKeeper.GetVestingInfo(ctx, ptypes.Elys)
	require.Equal(t, vestingInfo.BaseDenom, ptypes.Elys)
	require.Equal(t, vestingInfo.EpochIdentifier, "day")
	require.Equal(t, vestingInfo.NumEpochs, int64(params.BonusInfo.VestingDuration/86400))
	require.Equal(t, vestingInfo.NumMaxVestings, int64(10000))
	require.Equal(t, vestingInfo.VestNowFactor.String(), "0")
	require.Equal(t, vestingInfo.VestingDenom, ptypes.Elys)
}

func TestBeginBlocker(t *testing.T) {
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

	// set elys asset profile
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:       ptypes.Elys,
		Denom:           ptypes.Elys,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	// check vesting not started when the time not pass
	k.BeginBlocker(ctx)
	orders = k.GetAllOrders(ctx)
	require.Len(t, orders, 1)
	require.Equal(t, orders[0].VestingStarted, false)

	// check all vesting started when the time pass
	ctx = ctx.WithBlockTime(now.Add(time.Second * time.Duration(params.LaunchpadDuration+params.ReturnDuration+1)))
	k.BeginBlocker(ctx)
	orders = k.GetAllOrders(ctx)
	require.Len(t, orders, 1)
	require.Equal(t, orders[0].VestingStarted, true)

	commitments := app.CommitmentKeeper.GetCommitments(ctx, buyer.String())
	require.Len(t, commitments.VestingTokens, 1)
	require.Equal(t, commitments.VestingTokens[0].Denom, ptypes.Elys)
	require.Equal(t, commitments.VestingTokens[0].TotalAmount.String(), "1333333")
	require.Equal(t, commitments.VestingTokens[0].UnvestedAmount.String(), "1333333")
	require.Equal(t, commitments.VestingTokens[0].VestStartedTimestamp, int64(ctx.BlockTime().Unix()))
}
