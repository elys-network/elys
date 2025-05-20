package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/stretchr/testify/require"

	simapp "github.com/elys-network/elys/v4/app"

	"github.com/elys-network/elys/v4/x/commitment/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
)

// TestKeeper_CommitmentVestingInfo tests the CommitmentVestingInfo method
func TestKeeper_CommitmentVestingInfo(t *testing.T) {
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

	actualRes, err := mk.CommitmentVestingInfo(ctx, &types.QueryCommitmentVestingInfoRequest{
		Address: creator.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, actualRes)

	expectedRes := &types.QueryCommitmentVestingInfoResponse{
		Total: sdkmath.NewInt(50),
		VestingDetails: []types.VestingDetails{
			{
				Id:              "0",
				TotalVesting:    sdkmath.NewInt(100),
				Claimed:         sdkmath.NewInt(50),
				VestedSoFar:     sdkmath.NewInt(-10),
				RemainingBlocks: 11,
			},
		},
	}

	require.Equal(t, expectedRes, actualRes)
}

// TestKeeper_CommitmentVestingInfoNilRequest tests the CommitmentVestingInfo method with nil request
func TestKeeper_CommitmentVestingInfoNilRequest(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk := app.CommitmentKeeper

	_, err := mk.CommitmentVestingInfo(ctx, nil)
	require.Error(t, err)
}
