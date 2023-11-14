package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestDepositLiquidTokens(t *testing.T) {
	app := simapp.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], edenToken)
	require.NoError(t, err)

	creator := addr[0]

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true})

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdk.NewInt(50),
			},
		},
		RewardsUnclaimed: sdk.Coins{},
		Claimed:          sdk.Coins{sdk.NewCoin(ptypes.Eden, sdk.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Deposit liquid eden to become claimed state
	keeper.DepositLiquidTokensClaimed(ctx, ptypes.Eden, sdk.NewInt(100), creator.String())

	// Check if the deposit tokens were added to commitments
	newCommitments := keeper.GetCommitments(ctx, creator.String())

	// Check if the claimed tokens were updated correctly
	claimed := newCommitments.GetClaimedForDenom(ptypes.Eden)
	require.Equal(t, sdk.NewInt(250), claimed, "claimed tokens were not updated correctly")

	// Check if the committed tokens were updated correctly
	committedToken := newCommitments.GetCommittedAmountForDenom(ptypes.Eden)
	require.Equal(t, sdk.NewInt(50), committedToken, "committed tokens were not updated correctly")

	edenCoin := app.BankKeeper.GetBalance(ctx, addr[0], ptypes.Eden)
	require.Equal(t, edenCoin.Amount, sdk.ZeroInt())
}
