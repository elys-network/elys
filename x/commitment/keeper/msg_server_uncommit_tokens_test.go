package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"

	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUncommitTokens(t *testing.T) {
	app := simapp.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := "test_denom"
	initialUncommitted := sdk.NewInt(400)
	initialCommitted := sdk.NewInt(100)
	uncommitAmount := sdk.NewInt(100)

	// Set up initial commitments object with sufficient uncommitted & committed tokens
	uncommittedTokens := types.UncommittedTokens{
		Denom:  denom,
		Amount: initialUncommitted,
	}

	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: initialCommitted,
	}

	initialCommitments := types.Commitments{
		Creator:           creator,
		UncommittedTokens: []*types.UncommittedTokens{&uncommittedTokens},
		CommittedTokens:   []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: denom, CommitEnabled: true})

	// Call the UncommitTokens function
	msg := types.MsgUncommitTokens{
		Creator: creator,
		Amount:  uncommitAmount,
		Denom:   denom,
	}
	_, err := msgServer.UncommitTokens(ctx, &msg)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments, found := keeper.GetCommitments(ctx, creator)
	assert.True(t, found, "Commitments not found in the store")

	// Check if the committed tokens have the expected values
	assert.Equal(t, creator, commitments.Creator, "Incorrect creator")
	assert.Len(t, commitments.CommittedTokens, 1, "Incorrect number of committed tokens")
	assert.Equal(t, denom, commitments.CommittedTokens[0].Denom, "Incorrect denom")
	assert.Equal(t, sdk.ZeroInt(), commitments.CommittedTokens[0].Amount, "Incorrect amount")

	uncommittedToken := sdk.NewCoins(sdk.NewCoin(denom, uncommitAmount))

	edenCoin := app.BankKeeper.GetBalance(ctx, addr[0], denom)
	require.Equal(t, sdk.Coins{edenCoin}, uncommittedToken)
}
