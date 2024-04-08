package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestExternalIncentive(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	externalIncentives := []types.ExternalIncentive{
		{
			Id:             0,
			RewardDenom:    "reward1",
			PoolId:         1,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdk.OneInt(),
		},
		{
			Id:             1,
			RewardDenom:    "reward1",
			PoolId:         1,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdk.OneInt(),
		},
		{
			Id:             2,
			RewardDenom:    "reward1",
			PoolId:         2,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdk.OneInt(),
		},
	}
	require.Equal(t, app.MasterchefKeeper.GetExternalIncentiveIndex(ctx), uint64(0))
	for _, externalIncentive := range externalIncentives {
		err := app.MasterchefKeeper.SetExternalIncentive(ctx, externalIncentive)
		require.NoError(t, err)
	}
	require.Equal(t, app.MasterchefKeeper.GetExternalIncentiveIndex(ctx), uint64(3))
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
