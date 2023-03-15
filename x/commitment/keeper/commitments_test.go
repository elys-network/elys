package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestKeeper_SetGetRemoveCommitments(t *testing.T) {
	k, ctx := keepertest.CommitmentKeeper(t)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	commitments := types.Commitments{
		Creator: addr.String(),
	}

	// Test SetCommitments
	k.SetCommitments(ctx, commitments)

	// Test GetCommitments
	retrievedCommitments, found := k.GetCommitments(ctx, addr.String())
	require.True(t, found)
	assert.Equal(t, commitments, retrievedCommitments)

	// Test RemoveCommitments
	k.RemoveCommitments(ctx, addr.String())

	// Test that commitments are removed
	_, found = k.GetCommitments(ctx, addr.String())
	assert.False(t, found)
}
