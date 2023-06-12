package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/elys-network/elys/x/epochs/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func Test_validateParams(t *testing.T) {
	params := types.DefaultParams()

	// default params have no error
	require.NoError(t, params.Validate())

	// validate mincommision
	params.RewardPortionForLps = sdk.NewDecWithPrec(12, 1)
	require.Error(t, params.Validate())

	params.CommunityTax = sdk.NewDecWithPrec(1, 0)
	require.Error(t, params.Validate())

	lpIncentive := types.IncentiveInfo{
		// reward amount
		Amount: sdk.NewInt(10000),
		// epoch identifier
		EpochIdentifier: etypes.WeekEpochID,
		// start_time of the distribution
		StartTime: time.Now(),
		// distribution duration
		NumEpochs:    0,
		CurrentEpoch: 0,
		EdenBoostApr: 100,
	}

	params.LpIncentives = append(params.LpIncentives, lpIncentive)
	require.Error(t, params.Validate())
}
