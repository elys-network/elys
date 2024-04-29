package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"

	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUncommitTokens(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := "testdenom"
	initialCommitted := sdk.NewInt(100)
	uncommitAmount := sdk.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: initialCommitted,
	}

	initialCommitments := types.Commitments{
		Creator:         creator,
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: denom, CommitEnabled: true, WithdrawEnabled: true})

	// Add coins on commitment module
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(denom, initialCommitted)})
	require.NoError(t, err)

	// Call the UncommitTokens function
	err = keeper.UncommitTokens(ctx, addr[0], denom, uncommitAmount)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments := keeper.GetCommitments(ctx, creator)

	// Check if the committed tokens have the expected values
	assert.Equal(t, creator, commitments.Creator, "Incorrect creator")
	assert.Len(t, commitments.CommittedTokens, 1, "Incorrect number of committed tokens")
	assert.Equal(t, denom, commitments.CommittedTokens[0].Denom, "Incorrect denom")
	assert.Equal(t, sdk.ZeroInt(), commitments.CommittedTokens[0].Amount, "Incorrect amount")

	rewardUnclaimed := sdk.NewCoins(sdk.NewCoin(denom, uncommitAmount))

	edenCoin := app.BankKeeper.GetBalance(ctx, addr[0], denom)
	require.Equal(t, sdk.Coins{edenCoin}, rewardUnclaimed)
}
