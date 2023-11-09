package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
)

func TestCommitments_AddCommitedTokens(t *testing.T) {
	commitments := types.Commitments{
		Creator:          "",
		CommittedTokens:  []*types.CommittedTokens{},
		RewardsUnclaimed: []*types.RewardsUnclaimed{},
		VestingTokens:    []*types.VestingTokens{},
	}

	commitments.AddCommitedTokens("lp/1", sdk.NewInt(100), 100)
	commitments.AddCommitedTokens("lp/1", sdk.NewInt(100), 150)
	commitments.AddCommitedTokens("lp/1", sdk.NewInt(100), 200)
	commitments.AddCommitedTokens("lp/2", sdk.NewInt(100), 100)

	require.Len(t, commitments.CommittedTokens, 2)
	require.Len(t, commitments.CommittedTokens[0].Lockups, 3)
	require.Equal(t, commitments.CommittedTokens[0].Lockups[0].Amount.String(), "100")
	require.Equal(t, commitments.CommittedTokens[0].Lockups[0].UnlockTimestamp, uint64(100))
	require.Len(t, commitments.CommittedTokens[1].Lockups, 1)
}
func TestCommitments_WithdrawCommitedTokens(t *testing.T) {
	commitments := types.Commitments{
		Creator:          "",
		CommittedTokens:  []*types.CommittedTokens{},
		RewardsUnclaimed: []*types.RewardsUnclaimed{},
		VestingTokens:    []*types.VestingTokens{},
	}

	commitments.AddCommitedTokens("lp/1", sdk.NewInt(100), 100)
	commitments.AddCommitedTokens("lp/1", sdk.NewInt(100), 150)
	commitments.AddCommitedTokens("lp/1", sdk.NewInt(100), 200)
	commitments.AddCommitedTokens("lp/2", sdk.NewInt(100), 100)

	err := commitments.DeductFromCommitted("lp/1", sdk.NewInt(100), 100)
	require.NoError(t, err)

	err = commitments.DeductFromCommitted("lp/1", sdk.NewInt(100), 100)
	require.Error(t, err)

	err = commitments.DeductFromCommitted("lp/2", sdk.NewInt(200), 100)
	require.Error(t, err)
}
