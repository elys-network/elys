package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	// Create a test context and keeper
	testapp := app.InitElysTestApp(true)

	ctx := testapp.BaseApp.NewContext(false, tmproto.Header{})
	k := testapp.CommitmentKeeper
	require.NotNil(t, k)

	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:       "eden",
			VestingDenom:    "uelys",
			EpochIdentifier: "tenseconds",
			NumEpochs:       10,
			VestNowFactor:   sdk.NewInt(90),
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
	}

	k.SetParams(ctx, params)

	p := k.GetParams(ctx)
	require.EqualValues(t, params, p)
}

func TestEncodeDecodeParams(t *testing.T) {
	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:       "eden",
			VestingDenom:    "uelys",
			EpochIdentifier: "tenseconds",
			NumEpochs:       10,
			VestNowFactor:   sdk.NewInt(90),
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
	}

	encoded, err := types.ModuleCdc.MarshalJSON(&params)
	require.NoError(t, err)

	var decoded types.Params
	err = types.ModuleCdc.UnmarshalJSON(encoded, &decoded)
	require.NoError(t, err)

	require.EqualValues(t, params, decoded)
}

func TestGetParamsNew(t *testing.T) {
	// Create a test context and keeper
	testapp := app.InitElysTestApp(true)

	ctx := testapp.BaseApp.NewContext(false, tmproto.Header{})
	k := testapp.CommitmentKeeper
	require.NotNil(t, k)

	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:       "eden",
			VestingDenom:    "uelys",
			EpochIdentifier: "tenseconds",
			NumEpochs:       10,
			VestNowFactor:   sdk.NewInt(90),
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
	}

	k.SetParams(ctx, params)

	// Create a new context to test GetParams
	newCtx := testapp.BaseApp.NewContext(false, tmproto.Header{})
	p := k.GetParams(newCtx)
	require.EqualValues(t, params, p)
}

func TestGetVestingInfo(t *testing.T) {
	// Create a test context and keeper
	testapp := app.InitElysTestApp(true)

	ctx := testapp.BaseApp.NewContext(false, tmproto.Header{})
	k := testapp.CommitmentKeeper
	require.NotNil(t, k)

	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:       "eden",
			VestingDenom:    "uelys",
			EpochIdentifier: "tenseconds",
			NumEpochs:       10,
			VestNowFactor:   sdk.NewInt(90),
		},
		{
			BaseDenom:       "test",
			VestingDenom:    "test_vesting",
			EpochIdentifier: "tenseconds",
			NumEpochs:       10,
			VestNowFactor:   sdk.NewInt(90),
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
	}

	k.SetParams(ctx, params)

	// Test GetVestingInfo with existing base denom
	vestingInfo := k.GetVestingInfo(ctx, "eden")
	require.NotNil(t, vestingInfo)
	require.Equal(t, "eden", vestingInfo.BaseDenom)
	require.Equal(t, "uelys", vestingInfo.VestingDenom)

	// Test GetVestingInfo with non-existing base denom
	vestingInfo = k.GetVestingInfo(ctx, "nonexistent")
	require.Nil(t, vestingInfo)

	// Test GetVestingInfo with another existing base denom
	vestingInfo = k.GetVestingInfo(ctx, "test")
	require.NotNil(t, vestingInfo)
	require.Equal(t, "test", vestingInfo.BaseDenom)
	require.Equal(t, "test_vesting", vestingInfo.VestingDenom)
}
