package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitTokens(t *testing.T) {
	// Create a test context and keeper
	keeper, ctx := keepertest.CommitmentKeeper(t)
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Define the test data
	creator := "test_creator"
	denom := "test_denom"
	initialUncommitted := sdk.NewInt(500)
	commitAmount := sdk.NewInt(100)

	// Set up initial commitments object with sufficient uncommitted tokens
	uncommittedTokens := types.UncommittedTokens{
		Denom:  denom,
		Amount: initialUncommitted,
	}

	initialCommitments := types.Commitments{
		Creator:           creator,
		UncommittedTokens: []*types.UncommittedTokens{&uncommittedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Call the CommitTokens function
	msg := types.MsgCommitTokens{
		Creator: creator,
		Amount:  commitAmount,
		Denom:   denom,
	}
	_, err := msgServer.CommitTokens(ctx, &msg)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments, found := keeper.GetCommitments(ctx, creator)
	assert.True(t, found, "Commitments not found in the store")

	// Check if the committed tokens have the expected values
	assert.Equal(t, creator, commitments.Creator, "Incorrect creator")
	assert.Len(t, commitments.CommittedTokens, 1, "Incorrect number of committed tokens")
	assert.Equal(t, denom, commitments.CommittedTokens[0].Denom, "Incorrect denom")
	assert.Equal(t, commitAmount, commitments.CommittedTokens[0].Amount, "Incorrect amount")
}
