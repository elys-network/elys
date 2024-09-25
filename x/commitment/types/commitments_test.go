package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
)

func TestCommitments_AddCommittedTokens(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.AddCommittedTokens("lp/1", sdk.NewInt(100), 100)
	commitments.AddCommittedTokens("lp/1", sdk.NewInt(100), 150)
	commitments.AddCommittedTokens("lp/1", sdk.NewInt(100), 200)
	commitments.AddCommittedTokens("lp/2", sdk.NewInt(100), 100)

	require.Len(t, commitments.CommittedTokens, 2)
	require.Len(t, commitments.CommittedTokens[0].Lockups, 3)
	require.Equal(t, commitments.CommittedTokens[0].Lockups[0].Amount.String(), "100")
	require.Equal(t, commitments.CommittedTokens[0].Lockups[0].UnlockTimestamp, uint64(100))
	require.Len(t, commitments.CommittedTokens[1].Lockups, 1)

	commitments.AddCommittedTokens("lp/3", sdk.NewInt(1000), 100)
	commitments.AddCommittedTokens("lp/3", sdk.NewInt(2000), 120)
	commitments.AddCommittedTokens("lp/3", sdk.NewInt(3000), 130)

	require.Equal(t, commitments.CommittedTokens[2].Lockups[0].Amount.String(), "1000")
	require.Equal(t, commitments.CommittedTokens[2].Lockups[1].Amount.String(), "2000")
	require.Equal(t, commitments.CommittedTokens[2].Lockups[2].Amount.String(), "3000")

}

func TestCommitments_WithdrawCommitedTokens(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.AddCommittedTokens("lp/1", sdk.NewInt(100), 100)
	commitments.AddCommittedTokens("lp/1", sdk.NewInt(100), 150)
	commitments.AddCommittedTokens("lp/1", sdk.NewInt(100), 200)
	commitments.AddCommittedTokens("lp/2", sdk.NewInt(100), 100)

	err := commitments.DeductFromCommitted("lp/1", sdk.NewInt(100), 100, false)
	require.NoError(t, err)

	err = commitments.DeductFromCommitted("lp/1", sdk.NewInt(100), 100, false)
	require.Error(t, err)

	err = commitments.DeductFromCommitted("lp/2", sdk.NewInt(200), 100, false)
	require.Error(t, err)
}

func TestLockupAmount_WithdrawCommited(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.AddCommittedTokens("lp/1", sdk.NewInt(1000), 1)
	commitments.AddCommittedTokens("lp/1", sdk.NewInt(5000), 2)
	commitments.AddCommittedTokens("lp/1", sdk.NewInt(3000), 4)

	err := commitments.DeductFromCommitted("lp/1", sdk.NewInt(9000), 3, false)
	require.Error(t, err)
}
