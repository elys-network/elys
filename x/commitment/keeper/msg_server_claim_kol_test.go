package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v5/app"
	commitmentkeeper "github.com/elys-network/elys/v5/x/commitment/keeper"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestClaimKol(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	params := keeper.GetParams(ctx)
	params.EnableClaim = true
	keeper.SetParams(ctx, params)

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)
	addr := simapp.AddTestAddrs(app, ctx, 3, sdkmath.NewInt(1000000))
	creator := addr[0]

	commitmentkeeper.KolWallet = addr[2].String()
	kolAddress, _ := sdk.AccAddressFromBech32(commitmentkeeper.KolWallet)

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(2000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, kolAddress, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(2000))))
	require.NoError(t, err)

	keeper.SetKol(ctx, types.KolList{
		Address:  creator.String(),
		Amount:   sdkmath.NewInt(1000),
		Claimed:  false,
		Refunded: false,
	})

	claimKolMsg := &types.MsgClaimKol{
		ClaimAddress: creator.String(),
		Refund:       false,
	}

	_, err = msgServer.ClaimKol(ctx, claimKolMsg)
	require.NoError(t, err)

	// should be 12.5% of total amount and claimed should be true
	balances := app.BankKeeper.GetAllBalances(ctx, creator)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000125))), balances)

	walletBalances := app.BankKeeper.GetAllBalances(ctx, kolAddress)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1001875))), walletBalances)

	_, err = msgServer.ClaimKol(ctx, claimKolMsg)
	require.True(t, types.ErrKolAlreadyClaimed.Is(err), "error should be invalid denom")

	// Wrong block height
	params.StartKolClaimHeight = 100
	params.EndKolClaimHeight = 200
	keeper.SetParams(ctx, params)

	keeper.SetKol(ctx, types.KolList{
		Address:  addr[1].String(),
		Amount:   sdkmath.NewInt(1000),
		Claimed:  false,
		Refunded: false,
	})

	claimKolMsg = &types.MsgClaimKol{
		ClaimAddress: addr[1].String(),
		Refund:       false,
	}

	ctx = ctx.WithBlockHeight(50)
	_, err = msgServer.ClaimKol(ctx, claimKolMsg)
	require.True(t, types.ErrAirdropNotStarted.Is(err), "error should be airdop not started")

	ctx = ctx.WithBlockHeight(120)
	claimKolMsg.Refund = true
	_, err = msgServer.ClaimKol(ctx, claimKolMsg)
	require.NoError(t, err)

	kol := keeper.GetKol(ctx, addr[1])
	require.True(t, kol.Refunded)

	_, err = msgServer.ClaimKol(ctx, claimKolMsg)
	require.True(t, types.ErrKolRefunded.Is(err), "error should be invalid denom")
}
