package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetGetRemoveCommitments(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	commitments := types.Commitments{
		Creator: addr.String(),
	}

	// Test SetCommitments
	keeper.SetCommitments(ctx, commitments)

	// Test GetCommitments
	retrievedCommitments := keeper.GetCommitments(ctx, addr)
	assert.Equal(t, commitments, retrievedCommitments)

	// Test RemoveCommitments
	keeper.RemoveCommitments(ctx, addr)

	// Test that commitments are removed
	commitments = keeper.GetCommitments(ctx, addr)
	assert.True(t, commitments.IsEmpty())
}

// TestKeeper_GetAllCommitments tests the GetAllCommitments function
func TestKeeper_GetAllCommitments(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	commitments := []*types.Commitments{
		{
			Creator: addr1.String(),
		},
	}

	// Set the commitments
	for _, commitment := range commitments {
		keeper.SetCommitments(ctx, *commitment)
	}

	// Test GetAllCommitments
	retrievedCommitments := keeper.GetAllCommitments(ctx)
	assert.Equal(t, commitments, retrievedCommitments)
}

// TestKeeper_GetAllLegacyCommitments tests the GetAllLegacyCommitments function
func TestKeeper_GetAllLegacyCommitments(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	commitments := []*types.Commitments{
		{
			Creator: addr1.String(),
		},
	}
	// Set the commitments
	for _, commitment := range commitments {
		keeper.SetLegacyCommitments(ctx, *commitment)
	}

	// Test GetAllLegacyCommitments
	retrievedCommitments := keeper.GetAllLegacyCommitments(ctx)
	assert.Equal(t, commitments, retrievedCommitments)
}

// TestKeeper_DeleteLegacyCommitments tests the DeleteLegacyCommitments function
func TestKeeper_DeleteLegacyCommitments(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	commitments := types.Commitments{
		Creator: addr.String(),
	}

	// Set the commitments
	keeper.SetLegacyCommitments(ctx, commitments)

	// Test DeleteLegacyCommitments
	keeper.DeleteLegacyCommitments(ctx, addr.String())

	// Test that commitments are removed
	found := keeper.HasLegacyCommitments(ctx, addr.String())
	assert.False(t, found)
}

// TestKeeper_IterateCommitments tests the IterateCommitments function
func TestKeeper_IterateCommitments(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	commitments := []*types.Commitments{
		{
			Creator: addr1.String(),
		},
	}

	// Set the commitments
	for _, commitment := range commitments {
		keeper.SetCommitments(ctx, *commitment)
	}

	// Test IterateCommitments
	var retrievedCommitments []*types.Commitments
	keeper.IterateCommitments(ctx, func(commitment types.Commitments) bool {
		retrievedCommitments = append(retrievedCommitments, &commitment)
		return false
	})
	assert.Equal(t, commitments, retrievedCommitments)
}

// TestKeeper_IterateCommitments tests the IterateCommitments function with handlerFn returning true
func TestKeeper_IterateCommitmentsWithHandlerFn(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	commitments := []*types.Commitments{
		{
			Creator: addr1.String(),
		},
	}

	// Set the commitments
	for _, commitment := range commitments {
		keeper.SetCommitments(ctx, *commitment)
	}

	// Test IterateCommitments
	var retrievedCommitments []*types.Commitments
	keeper.IterateCommitments(ctx, func(commitment types.Commitments) bool {
		retrievedCommitments = append(retrievedCommitments, &commitment)
		return true
	})
	assert.Equal(t, commitments[:1], retrievedCommitments)
}

// TestKeeper_TotalNumberOfCommitments tests the TotalNumberOfCommitments function
func TestKeeper_TotalNumberOfCommitments(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	commitments := []*types.Commitments{
		{
			Creator: addr1.String(),
		},
	}

	// Set the commitments
	for _, commitment := range commitments {
		keeper.SetCommitments(ctx, *commitment)
	}

	// Test TotalNumberOfCommitments
	total := keeper.TotalNumberOfCommitments(ctx)
	assert.Equal(t, len(commitments), int(total))
}

// TestKeeper_BurnEdenBoost tests the BurnEdenBoost function
func TestKeeper_BurnEdenBoost(t *testing.T) {
	// define a test matrix that will cover all the use cases
	tests := []struct {
		name            string
		claimedAmount   math.Int
		committedTokens math.Int
		deductAmount    math.Int
		expectedError   bool
	}{
		{
			name:            "deduct amount is zero",
			claimedAmount:   sdk.NewInt(100),
			committedTokens: sdk.NewInt(100),
			deductAmount:    sdk.NewInt(0),
			expectedError:   false,
		},
		{
			name:            "deduct amount is greater than claimed amount",
			claimedAmount:   sdk.NewInt(100),
			committedTokens: sdk.NewInt(100),
			deductAmount:    sdk.NewInt(200),
			expectedError:   false,
		},
		{
			name:            "deduct amount is greater than claimed amount with no committed tokens",
			claimedAmount:   sdk.NewInt(100),
			committedTokens: sdk.NewInt(0),
			deductAmount:    sdk.NewInt(200),
			expectedError:   false,
		},
		{
			name:            "deduct amount is less than claimed amount",
			claimedAmount:   sdk.NewInt(100),
			committedTokens: sdk.NewInt(100),
			deductAmount:    sdk.NewInt(50),
			expectedError:   false,
		},
		{
			name:            "deduct amount is equal to claimed amount",
			claimedAmount:   sdk.NewInt(100),
			committedTokens: sdk.NewInt(100),
			deductAmount:    sdk.NewInt(100),
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := keepertest.CommitmentKeeper(t)

			addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			commitments := types.Commitments{
				Creator: addr.String(),
			}

			// Set the commitments
			keeper.SetCommitments(ctx, commitments)

			// Set the claimed amount
			commitments.AddClaimed(sdk.NewCoin("denom", tt.claimedAmount))
			keeper.SetCommitments(ctx, commitments)

			// Add committed amount
			commitments.AddCommittedTokens("denom", tt.committedTokens, 0)
			keeper.SetCommitments(ctx, commitments)

			// Test BurnEdenBoost
			err := keeper.BurnEdenBoost(ctx, addr, "denom", tt.deductAmount)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
