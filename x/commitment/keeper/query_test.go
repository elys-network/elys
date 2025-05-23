package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v5/app"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/require"
)

// TestKeeper_ShowCommitments tests the ShowCommitments function
func TestKeeper_ShowCommitments(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk := app.CommitmentKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000000000))

	// Define the test data
	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:                ptypes.Eden,
				TotalAmount:          sdkmath.NewInt(100),
				ClaimedAmount:        sdkmath.NewInt(50),
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
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk := app.CommitmentKeeper

	_, err := mk.ShowCommitments(ctx, nil)
	require.Error(t, err)
}

// TestKeeper_NumberOfCommitments tests the NumberOfCommitments function
func TestKeeper_NumberOfCommitments(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk := app.CommitmentKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000000000))

	// Define the test data
	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:                ptypes.Eden,
				TotalAmount:          sdkmath.NewInt(100),
				ClaimedAmount:        sdkmath.NewInt(50),
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
		Number: 3, // set to 3 because end block from estaking
	}

	require.Equal(t, expectedRes, actualRes)
}

// TestKeeper_NumberOfCommitmentsNilRequest tests the case where the request is nil
func TestKeeper_NumberOfCommitmentsNilRequest(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk := app.CommitmentKeeper

	_, err := mk.NumberOfCommitments(ctx, nil)
	require.Error(t, err)
}

// TestKeeper_CommittedTokensLocked tests the CommittedTokensLocked function
func TestKeeper_CommittedTokensLocked(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk := app.CommitmentKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000000000))

	// Define the test data
	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:                ptypes.Eden,
				TotalAmount:          sdkmath.NewInt(100),
				ClaimedAmount:        sdkmath.NewInt(50),
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
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk := app.CommitmentKeeper

	_, err := mk.CommittedTokensLocked(ctx, nil)
	require.Error(t, err)
}
