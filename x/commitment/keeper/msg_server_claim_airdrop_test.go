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

func TestAirdropClaim(t *testing.T) {
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

	commitmentkeeper.AirdropWallet = addr[2].String()
	airdropAddress, _ := sdk.AccAddressFromBech32(commitmentkeeper.AirdropWallet)

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(2000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, airdropAddress, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(2000))))
	require.NoError(t, err)

	keeper.SetAtomStaker(ctx, types.AtomStaker{
		Address: creator.String(),
		Amount:  sdkmath.NewInt(100),
	})

	keeper.SetNFTHodler(ctx, types.NftHolder{
		Address: creator.String(),
		Amount:  sdkmath.NewInt(500),
	})

	keeper.SetCadet(ctx, types.Cadet{
		Address: creator.String(),
		Amount:  sdkmath.NewInt(250),
	})

	keeper.SetGovernor(ctx, types.Governor{
		Address: creator.String(),
		Amount:  sdkmath.NewInt(250),
	})

	claimAirdropMsg := &types.MsgClaimAirdrop{
		ClaimAddress: creator.String(),
	}

	_, err = msgServer.ClaimAirdrop(ctx, claimAirdropMsg)
	require.NoError(t, err)

	newCommitments := keeper.GetCommitments(ctx, creator)
	// check if the eden was increased tokens were updated correctly
	require.Equal(t, sdkmath.NewInt(100), newCommitments.GetClaimedForDenom(ptypes.Eden))

	// Try to claim airdrop again
	_, err = msgServer.ClaimAirdrop(ctx, claimAirdropMsg)
	require.Error(t, err, "should throw an error when trying claim airdrop again")
	require.True(t, types.ErrAirdropAlreadyClaimed.Is(err), "error should be airdrop already claimed")

	// Test for elys + eden rewards
	balances := app.BankKeeper.GetAllBalances(ctx, creator)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1001000))), balances)

	walletBalances := app.BankKeeper.GetAllBalances(ctx, airdropAddress)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1001000))), walletBalances)

	require.Equal(t, sdkmath.NewInt(100), newCommitments.GetClaimedForDenom(ptypes.Eden))

	// Wrong block height
	params.StartAirdropClaimHeight = 100
	params.EndAirdropClaimHeight = 200
	keeper.SetParams(ctx, params)

	keeper.SetAtomStaker(ctx, types.AtomStaker{
		Address: addr[1].String(),
		Amount:  sdkmath.NewInt(100),
	})

	claimAirdropMsg = &types.MsgClaimAirdrop{
		ClaimAddress: addr[1].String(),
	}

	ctx = ctx.WithBlockHeight(50)
	_, err = msgServer.ClaimAirdrop(ctx, claimAirdropMsg)
	require.True(t, types.ErrAirdropNotStarted.Is(err), "error should be invalid denom")

	ctx = ctx.WithBlockHeight(250)
	_, err = msgServer.ClaimAirdrop(ctx, claimAirdropMsg)
	require.True(t, types.ErrAirdropEnded.Is(err), "error should be invalid denom")
}
