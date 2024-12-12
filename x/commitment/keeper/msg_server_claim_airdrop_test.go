package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestAirdropClaim(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}

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
		Creator: creator.String(),
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
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000))), balances)

	walletBalances := app.BankKeeper.GetAllBalances(ctx, airdropAddress)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000))), walletBalances)

	require.Equal(t, sdkmath.NewInt(100), newCommitments.GetClaimedForDenom(ptypes.Eden))

	// Wrong block height
	params := keeper.GetParams(ctx)
	params.StartAirdropClaimHeight = 100
	params.EndAirdropClaimHeight = 200
	keeper.SetParams(ctx, params)

	// ctx = ctx.WithBlockHeight(50)
	// _, err = msgServer.ClaimAirdrop(ctx, claimAirdropMsg)
	// require.Error(t, err)
}
