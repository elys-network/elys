package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestExternalIncentive(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	externalIncentives := []types.ExternalIncentive{
		{
			Id:             0,
			RewardDenom:    "reward1",
			PoolId:         1,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdkmath.OneInt(),
			Apr:            sdkmath.LegacyZeroDec(),
		},
		{
			Id:             1,
			RewardDenom:    "reward1",
			PoolId:         1,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdkmath.OneInt(),
			Apr:            sdkmath.LegacyZeroDec(),
		},
		{
			Id:             2,
			RewardDenom:    "reward1",
			PoolId:         2,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdkmath.OneInt(),
			Apr:            sdkmath.LegacyZeroDec(),
		},
	}
	for _, externalIncentive := range externalIncentives {
		app.MasterchefKeeper.SetExternalIncentive(ctx, externalIncentive)
	}
	for _, externalIncentive := range externalIncentives {
		info, found := app.MasterchefKeeper.GetExternalIncentive(ctx, externalIncentive.Id)
		require.True(t, found)
		require.Equal(t, info, externalIncentive)
	}
	externalIncentivesStored := app.MasterchefKeeper.GetAllExternalIncentives(ctx)
	require.Len(t, externalIncentivesStored, 3)

	app.MasterchefKeeper.RemoveExternalIncentive(ctx, externalIncentives[0].Id)
	externalIncentivesStored = app.MasterchefKeeper.GetAllExternalIncentives(ctx)
	require.Len(t, externalIncentivesStored, 2)
}
