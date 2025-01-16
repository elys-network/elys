package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	// Create a test context and keeper
	testapp := app.InitElysTestApp(true, t)

	ctx := testapp.BaseApp.NewContext(true)
	k := testapp.CommitmentKeeper
	require.NotNil(t, k)

	vestingInfos := []types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdkmath.NewInt(90),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos:     vestingInfos,
		TotalEdenSupply:  math.ZeroInt(),
		TotalEdenbSupply: math.ZeroInt(),
	}

	k.SetParams(ctx, params)

	p := k.GetParams(ctx)
	require.EqualValues(t, params, p)
}

func TestEncodeDecodeParams(t *testing.T) {
	vestingInfos := []types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdkmath.NewInt(90),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos:     vestingInfos,
		TotalCommitted:   sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 10)},
		TotalEdenSupply:  math.ZeroInt(),
		TotalEdenbSupply: math.ZeroInt(),
	}
	encoded, err := codec.NewLegacyAmino().MarshalJSON(&params)
	require.NoError(t, err)

	var decoded types.Params
	err = codec.NewLegacyAmino().UnmarshalJSON(encoded, &decoded)
	require.NoError(t, err)

	require.EqualValues(t, params, decoded)
}

func TestGetParamsNew(t *testing.T) {
	// Create a test context and keeper
	testapp := app.InitElysTestApp(true, t)

	ctx := testapp.BaseApp.NewContext(true)
	k := testapp.CommitmentKeeper
	require.NotNil(t, k)

	vestingInfos := []types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdkmath.NewInt(90),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos:     vestingInfos,
		TotalEdenSupply:  math.ZeroInt(),
		TotalEdenbSupply: math.ZeroInt(),
	}

	k.SetParams(ctx, params)

	// Create a new context to test GetParams
	newCtx := testapp.BaseApp.NewContext(true)
	p := k.GetParams(newCtx)
	require.EqualValues(t, params, p)
}

func TestGetVestingInfo(t *testing.T) {
	// Create a test context and keeper
	testapp := app.InitElysTestApp(true, t)

	ctx := testapp.BaseApp.NewContext(true)
	k := testapp.CommitmentKeeper
	require.NotNil(t, k)

	vestingInfos := []types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdkmath.NewInt(90),
			NumMaxVestings: 10,
		},
		{
			BaseDenom:      "test",
			VestingDenom:   "test_vesting",
			NumBlocks:      10,
			VestNowFactor:  sdkmath.NewInt(90),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
	}

	k.SetParams(ctx, params)

	// Test GetVestingInfo with existing base denom
	vestingInfo, _ := k.GetVestingInfo(ctx, ptypes.Eden)
	require.NotNil(t, vestingInfo)
	require.Equal(t, ptypes.Eden, vestingInfo.BaseDenom)
	require.Equal(t, ptypes.Elys, vestingInfo.VestingDenom)

	// Test GetVestingInfo with non-existing base denom
	vestingInfo, _ = k.GetVestingInfo(ctx, "nonexistent")
	require.Nil(t, vestingInfo)

	// Test GetVestingInfo with another existing base denom
	vestingInfo, _ = k.GetVestingInfo(ctx, "test")
	require.NotNil(t, vestingInfo)
	require.Equal(t, "test", vestingInfo.BaseDenom)
	require.Equal(t, "test_vesting", vestingInfo.VestingDenom)
}
