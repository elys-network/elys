package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func Test_validateParams(t *testing.T) {
	params := types.DefaultParams()

	// default params have no error
	require.NoError(t, params.Validate())

	// validate mincommision
	params.RewardPortionForLps = sdk.NewDecWithPrec(12, 1)
	require.Error(t, params.Validate())

	lpIncentive := types.IncentiveInfo{
		// reward amount in eden for 1 year
		EdenAmountPerYear: sdk.NewInt(10000000000000),
		// starting block height of the distribution
		DistributionStartBlock: sdk.ZeroInt(),
		// distribution duration - block number per year
		TotalBlocksPerYear: sdk.NewInt(10512000),
		// maximum eden allocation per day that won't exceed 30% apr
		MaxEdenPerAllocation: sdk.NewInt(27397238400),
		// current epoch in block number
		CurrentEpochInBlocks: sdk.NewInt(0),
	}

	params.LpIncentives = &lpIncentive
	require.Error(t, params.Validate())
}
