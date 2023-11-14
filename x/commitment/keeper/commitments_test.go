package keeper_test

import (
	"testing"

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
	retrievedCommitments := keeper.GetCommitments(ctx, addr.String())
	assert.Equal(t, commitments, retrievedCommitments)

	// Test RemoveCommitments
	keeper.RemoveCommitments(ctx, addr.String())

	// Test that commitments are removed
	commitments = keeper.GetCommitments(ctx, addr.String())
	assert.True(t, commitments.IsEmpty())
}
