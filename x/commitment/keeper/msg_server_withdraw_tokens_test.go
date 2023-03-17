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

func TestWithdrawTokens(t *testing.T) {
	t.Skip() // TODO test keeper needs bank keeper
	// Create a test context and keeper
	keeper, ctx := keepertest.CommitmentKeeper(t)
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Define the test data
	creator := "test_creator"
	denom := "test_denom"
	initialUncommitted := sdk.NewInt(50)
	initialCommitted := sdk.NewInt(100)

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

	// Test scenario 1: Withdraw within uncommitted balance
	msg := &types.MsgWithdrawTokens{
		Creator: creator,
		Amount:  sdk.NewInt(30),
		Denom:   denom,
	}

	_, err := msgServer.WithdrawTokens(ctx, msg)
	require.NoError(t, err)

	updatedCommitments, found := keeper.GetCommitments(ctx, creator)
	require.True(t, found)

	uncommittedBalance := updatedCommitments.GetUncommittedAmountForDenom(denom)
	assert.Equal(t, sdk.NewInt(20), uncommittedBalance)

	// Test scenario 2: Withdraw more than uncommitted balance but within total balance
	msg = &types.MsgWithdrawTokens{
		Creator: creator,
		Amount:  sdk.NewInt(70),
		Denom:   denom,
	}

	_, err = msgServer.WithdrawTokens(ctx, msg)
	require.NoError(t, err)

	updatedCommitments, found = keeper.GetCommitments(ctx, creator)
	require.True(t, found)

	uncommittedBalance = updatedCommitments.GetUncommittedAmountForDenom(denom)
	assert.Equal(t, sdk.NewInt(0), uncommittedBalance)

	committedBalance := updatedCommitments.GetCommittedAmountForDenom(denom)
	assert.Equal(t, sdk.NewInt(50), committedBalance)

	// Test scenario 3: Withdraw more than total balance
	msg = &types.MsgWithdrawTokens{
		Creator: creator,
		Amount:  sdk.NewInt(100),
		Denom:   denom,
	}

	_, err = msgServer.WithdrawTokens(ctx, msg)
	require.Error(t, err)
}
