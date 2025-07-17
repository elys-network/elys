package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/x/commitment/types"
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

// TestKeeper_IterateCommitmentsWithHandlerFn tests the IterateCommitments function with handlerFn returning true
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
			claimedAmount:   math.NewInt(100),
			committedTokens: math.NewInt(100),
			deductAmount:    math.NewInt(0),
			expectedError:   false,
		},
		{
			name:            "deduct amount is greater than claimed amount",
			claimedAmount:   math.NewInt(100),
			committedTokens: math.NewInt(0),
			deductAmount:    math.NewInt(100),
			expectedError:   false,
		},
		{
			name:            "deduct amount is greater than claimed amount with no committed tokens",
			claimedAmount:   math.NewInt(200),
			committedTokens: math.NewInt(0),
			deductAmount:    math.NewInt(200),
			expectedError:   false,
		},
		{
			name:            "deduct amount is less than claimed amount",
			claimedAmount:   math.NewInt(100),
			committedTokens: math.NewInt(100),
			deductAmount:    math.NewInt(50),
			expectedError:   false,
		},
		{
			name:            "deduct amount is equal to claimed amount",
			claimedAmount:   math.NewInt(100),
			committedTokens: math.NewInt(100),
			deductAmount:    math.NewInt(100),
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

			comm := keeper.GetCommitments(ctx, addr)
			totalPrev := comm.Claimed.AmountOf("denom").Add(comm.GetCommittedAmountForDenom("denom"))

			// Test BurnEdenBoost
			err := keeper.BurnEdenBoost(ctx, addr, "denom", tt.deductAmount)
			comm = keeper.GetCommitments(ctx, addr)
			total := comm.Claimed.AmountOf("denom").Add(comm.GetCommittedAmountForDenom("denom"))
			assert.Equal(t, tt.deductAmount.String(), totalPrev.Sub(total).String())
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
