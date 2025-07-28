package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v7/x/commitment/types"
	"github.com/stretchr/testify/require"
)

func TestCommitments_AddCommittedTokens(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.AddCommittedTokens("lp/1", math.NewInt(100), 100)
	commitments.AddCommittedTokens("lp/1", math.NewInt(100), 150)
	commitments.AddCommittedTokens("lp/1", math.NewInt(100), 200)
	commitments.AddCommittedTokens("lp/2", math.NewInt(100), 100)

	require.Len(t, commitments.CommittedTokens, 2)
	require.Len(t, commitments.CommittedTokens[0].Lockups, 3)
	require.Equal(t, commitments.CommittedTokens[0].Lockups[0].Amount.String(), "100")
	require.Equal(t, commitments.CommittedTokens[0].Lockups[0].UnlockTimestamp, uint64(100))
	require.Len(t, commitments.CommittedTokens[1].Lockups, 1)

	commitments.AddCommittedTokens("lp/3", math.NewInt(1000), 100)
	commitments.AddCommittedTokens("lp/3", math.NewInt(2000), 120)
	commitments.AddCommittedTokens("lp/3", math.NewInt(3000), 130)

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

	commitments.AddCommittedTokens("lp/1", math.NewInt(100), 100)
	commitments.AddCommittedTokens("lp/1", math.NewInt(100), 150)
	commitments.AddCommittedTokens("lp/1", math.NewInt(100), 200)
	commitments.AddCommittedTokens("lp/2", math.NewInt(100), 100)

	err := commitments.DeductFromCommitted("lp/1", math.NewInt(100), 100, false)
	require.NoError(t, err)

	err = commitments.DeductFromCommitted("lp/1", math.NewInt(100), 100, false)
	require.Error(t, err)

	err = commitments.DeductFromCommitted("lp/2", math.NewInt(200), 100, false)
	require.Error(t, err)
}

func TestLockupAmount_WithdrawCommited(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.AddCommittedTokens("lp/1", math.NewInt(1000), 1)
	commitments.AddCommittedTokens("lp/1", math.NewInt(5000), 2)
	commitments.AddCommittedTokens("lp/1", math.NewInt(3000), 4)

	err := commitments.DeductFromCommitted("lp/1", math.NewInt(9000), 3, false)
	require.Error(t, err)
}

// Test_Commitments_IsEmpty tests the IsEmpty method of the Commitments type
func Test_Commitments_IsEmpty(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	require.True(t, commitments.IsEmpty())

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:   "lp/1",
		Amount:  math.NewInt(100),
		Lockups: nil,
	})

	require.False(t, commitments.IsEmpty())

	commitments.CommittedTokens = nil
	commitments.Claimed = sdk.Coins{sdk.NewInt64Coin("lp/1", 100)}

	require.False(t, commitments.IsEmpty())

	commitments.Claimed = nil
	commitments.VestingTokens = append(commitments.VestingTokens, &types.VestingTokens{
		Denom:                "lp/1",
		TotalAmount:          math.NewInt(100),
		ClaimedAmount:        math.NewInt(0),
		NumBlocks:            100,
		StartBlock:           0,
		VestStartedTimestamp: 0,
	})

	require.False(t, commitments.IsEmpty())
}

// Test_Commitments_GetCommittedAmountForDenom tests the GetCommittedAmountForDenom method of the Commitments type
func Test_Commitments_GetCommittedAmountForDenom(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:   "lp/1",
		Amount:  math.NewInt(100),
		Lockups: nil,
	})

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:   "lp/2",
		Amount:  math.NewInt(200),
		Lockups: nil,
	})

	require.Equal(t, commitments.GetCommittedAmountForDenom("lp/1"), math.NewInt(100))
	require.Equal(t, commitments.GetCommittedAmountForDenom("lp/2"), math.NewInt(200))
	require.Equal(t, commitments.GetCommittedAmountForDenom("lp/3"), math.NewInt(0))
}

// Test_Commitments_GetCommittedLockUpsForDenom tests the GetCommittedLockUpsForDenom method of the Commitments type
func Test_Commitments_GetCommittedLockUpsForDenom(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:   "lp/1",
		Amount:  math.NewInt(100),
		Lockups: nil,
	})

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:  "lp/2",
		Amount: math.NewInt(200),
		Lockups: []types.Lockup{
			{
				Amount:          math.NewInt(100),
				UnlockTimestamp: 100,
			},
			{
				Amount:          math.NewInt(200),
				UnlockTimestamp: 200,
			},
		},
	})

	require.Nil(t, commitments.GetCommittedLockUpsForDenom("lp/1"))
	require.NotNil(t, commitments.GetCommittedLockUpsForDenom("lp/2"))
	require.Len(t, commitments.GetCommittedLockUpsForDenom("lp/2"), 2)
}

// Test_Commitments_GetCommittedLockUpsForDenomNil tests the GetCommittedLockUpsForDenom method of the Commitments type
// when there are no lockups for the specified denom
func Test_Commitments_GetCommittedLockUpsForDenomNil(t *testing.T) {
	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:   "lp/1",
		Amount:  math.NewInt(100),
		Lockups: nil,
	})

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:   "lp/2",
		Amount:  math.NewInt(200),
		Lockups: nil,
	})

	require.Nil(t, commitments.GetCommittedLockUpsForDenom("lp/3"))
}

// Test_Commitments_CommittedTokensLocked tests the CommittedTokensLocked method of the Commitments type
func Test_Commitments_CommittedTokensLocked(t *testing.T) {
	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil).
		WithBlockTime(time.Unix(0, 0))

	commitments := types.Commitments{
		Creator:         "",
		CommittedTokens: []*types.CommittedTokens{},
		VestingTokens:   []*types.VestingTokens{},
	}

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:   "lp/1",
		Amount:  math.NewInt(100),
		Lockups: nil,
	})

	commitments.CommittedTokens = append(commitments.CommittedTokens, &types.CommittedTokens{
		Denom:  "lp/2",
		Amount: math.NewInt(200),
		Lockups: []types.Lockup{
			{
				Amount:          math.NewInt(100),
				UnlockTimestamp: 100,
			},
			{
				Amount:          math.NewInt(200),
				UnlockTimestamp: 200,
			},
		},
	})

	locked, unlocked := commitments.CommittedTokensLocked(ctx)
	require.Equal(t, locked, sdk.Coins{sdk.NewInt64Coin("lp/2", 300)})
	require.Equal(t, unlocked, sdk.Coins{sdk.NewInt64Coin("lp/1", 100), sdk.NewInt64Coin("lp/2", 200)})
}
