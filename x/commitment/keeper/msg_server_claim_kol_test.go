package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestClaimKol(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)
	addr := simapp.AddTestAddrs(app, ctx, 3, sdkmath.NewInt(1000000))
	creator := addr[0]

	commitmentkeeper.KolWallet = addr[2].String()
	airdropAddress, _ := sdk.AccAddressFromBech32(commitmentkeeper.AirdropWallet)

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(2000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, airdropAddress, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(2000))))
	require.NoError(t, err)

	keeper.SetKol(ctx, types.KolList{
		Address:  creator.String(),
		Amount:   sdkmath.NewInt(100),
		Claimed:  false,
		Refunded: false,
	})

	claimKolMsg := &types.MsgClaimKol{
		ClaimAddress: creator.String(),
		Refund:       false,
	}

	_, err = msgServer.ClaimKol(ctx, claimKolMsg)
	require.NoError(t, err)

	// Test for elys + eden rewards
	balances := app.BankKeeper.GetAllBalances(ctx, creator)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1001000))), balances)

	walletBalances := app.BankKeeper.GetAllBalances(ctx, airdropAddress)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1001000))), walletBalances)

	// Wrong block height
	// params := keeper.GetParams(ctx)
	// params.StartAirdropClaimHeight = 100
	// params.EndAirdropClaimHeight = 200
	// keeper.SetParams(ctx, params)

	// keeper.SetAtomStaker(ctx, types.AtomStaker{
	// 	Address: addr[1].String(),
	// 	Amount:  sdkmath.NewInt(100),
	// })

	// claimAirdropMsg = &types.MsgClaimAirdrop{
	// 	ClaimAddress: addr[1].String(),
	// }

	// ctx = ctx.WithBlockHeight(50)
	// _, err = msgServer.ClaimAirdrop(ctx, claimAirdropMsg)
	// require.True(t, types.ErrAirdropNotStarted.Is(err), "error should be invalid denom")

	// ctx = ctx.WithBlockHeight(250)
	// _, err = msgServer.ClaimAirdrop(ctx, claimAirdropMsg)
	// require.True(t, types.ErrAirdropEnded.Is(err), "error should be invalid denom")
}
