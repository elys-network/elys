package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"

	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUncommitTokens(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(false)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := "testdenom"
	initialCommitted := sdkmath.NewInt(100)
	uncommitAmount := sdkmath.NewInt(100)

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
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: denom, CommitEnabled: true, WithdrawEnabled: true})

	// Add coins on commitment module
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(denom, initialCommitted)})
	require.NoError(t, err)

	// Call the UncommitTokens function
	err = keeper.UncommitTokens(ctx, addr[0], denom, uncommitAmount, false)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments := keeper.GetCommitments(ctx, sdk.MustAccAddressFromBech32(creator))

	// Check if the committed tokens have the expected values
	assert.Equal(t, len(commitments.CommittedTokens), 0, "Incorrect creator")

	rewardUnclaimed := sdk.NewCoins(sdk.NewCoin(denom, uncommitAmount))

	edenCoin := app.BankKeeper.GetBalance(ctx, addr[0], denom)
	require.Equal(t, sdk.Coins{edenCoin}, rewardUnclaimed)
}
