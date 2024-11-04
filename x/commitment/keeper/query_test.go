package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

// TestKeeper_ShowCommitments tests the ShowCommitments function
func TestKeeper_ShowCommitments(t *testing.T) {
	app := app.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk := app.CommitmentKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Define the test data
	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:                ptypes.Eden,
				TotalAmount:          sdk.NewInt(100),
				ClaimedAmount:        sdk.NewInt(50),
				NumBlocks:            10,
				StartBlock:           1,
				VestStartedTimestamp: 1,
			},
		},
	}

	mk.SetCommitments(ctx, commitments)

	actualRes, err := mk.ShowCommitments(ctx, &types.QueryShowCommitmentsRequest{
		Creator: creator.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, actualRes)

	expectedRes := &types.QueryShowCommitmentsResponse{
		Commitments: commitments,
	}

	require.Equal(t, expectedRes, actualRes)
}

// TestKeeper_ShowCommitmentsNilRequest tests the case where the request is nil
func TestKeeper_ShowCommitmentsNilRequest(t *testing.T) {
	app := app.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk := app.CommitmentKeeper

	_, err := mk.ShowCommitments(ctx, nil)
	require.Error(t, err)
}

// TestKeeper_NumberOfCommitments tests the NumberOfCommitments function
func TestKeeper_NumberOfCommitments(t *testing.T) {
	app := app.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk := app.CommitmentKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Define the test data
	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:                ptypes.Eden,
				TotalAmount:          sdk.NewInt(100),
				ClaimedAmount:        sdk.NewInt(50),
				NumBlocks:            10,
				StartBlock:           1,
				VestStartedTimestamp: 1,
			},
		},
	}

	mk.SetCommitments(ctx, commitments)

	actualRes, err := mk.NumberOfCommitments(ctx, &types.QueryNumberOfCommitmentsRequest{})
	require.NoError(t, err)
	require.NotNil(t, actualRes)

	expectedRes := &types.QueryNumberOfCommitmentsResponse{
		Number: 1,
	}

	require.Equal(t, expectedRes, actualRes)
}

// TestKeeper_NumberOfCommitmentsNilRequest tests the case where the request is nil
func TestKeeper_NumberOfCommitmentsNilRequest(t *testing.T) {
	app := app.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk := app.CommitmentKeeper

	_, err := mk.NumberOfCommitments(ctx, nil)
	require.Error(t, err)
}

// TestKeeper_CommittedTokensLocked tests the CommittedTokensLocked function
func TestKeeper_CommittedTokensLocked(t *testing.T) {
	app := app.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk := app.CommitmentKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Define the test data
	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:                ptypes.Eden,
				TotalAmount:          sdk.NewInt(100),
				ClaimedAmount:        sdk.NewInt(50),
				NumBlocks:            10,
				StartBlock:           1,
				VestStartedTimestamp: 1,
			},
		},
	}

	mk.SetCommitments(ctx, commitments)

	actualRes, err := mk.CommittedTokensLocked(ctx, &types.QueryCommittedTokensLockedRequest{
		Address: creator.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, actualRes)

	expectedRes := &types.QueryCommittedTokensLockedResponse{
		Address:         creator.String(),
		LockedCommitted: sdk.NewCoins(),
		TotalCommitted:  sdk.NewCoins(),
	}

	require.Equal(t, expectedRes, actualRes)
}

// TestKeeper_CommittedTokensLockedNilRequest tests the case where the request is nil
func TestKeeper_CommittedTokensLockedNilRequest(t *testing.T) {
	app := app.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk := app.CommitmentKeeper

	_, err := mk.CommittedTokensLocked(ctx, nil)
	require.Error(t, err)
}
