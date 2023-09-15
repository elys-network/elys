package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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
			BaseDenom:       ptypes.Eden,
			VestingDenom:    ptypes.Elys,
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
			BaseDenom:       ptypes.Eden,
			VestingDenom:    ptypes.Elys,
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
			BaseDenom:       ptypes.Eden,
			VestingDenom:    ptypes.Elys,
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
			BaseDenom:       ptypes.Eden,
			VestingDenom:    ptypes.Elys,
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
	vestingInfo := k.GetVestingInfo(ctx, ptypes.Eden)
	require.NotNil(t, vestingInfo)
	require.Equal(t, ptypes.Eden, vestingInfo.BaseDenom)
	require.Equal(t, ptypes.Elys, vestingInfo.VestingDenom)

	// Test GetVestingInfo with non-existing base denom
	vestingInfo = k.GetVestingInfo(ctx, "nonexistent")
	require.Nil(t, vestingInfo)

	// Test GetVestingInfo with another existing base denom
	vestingInfo = k.GetVestingInfo(ctx, "test")
	require.NotNil(t, vestingInfo)
	require.Equal(t, "test", vestingInfo.BaseDenom)
	require.Equal(t, "test_vesting", vestingInfo.VestingDenom)
}
