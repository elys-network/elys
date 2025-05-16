package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/elys-network/elys/app"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddVestingInfo(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	govAddress := sdk.AccAddress(address.Module("gov"))

	// Define the test data
	signer := govAddress.String()

	// Call the Update Vesting Info function
	msg := types.MsgUpdateVestingInfo{
		Authority:      signer,
		BaseDenom:      "test_denom",
		VestingDenom:   "test_denom",
		NumBlocks:      10,
		VestNowFactor:  10,
		NumMaxVestings: 10,
	}
	_, err := msgServer.UpdateVestingInfo(ctx, &msg)
	require.NoError(t, err)

	// Check if the vesting info has been added to the store
	vestingInfo, _ := keeper.GetVestingInfo(ctx, "test_denom")
	assert.Equal(t, vestingInfo.BaseDenom, "test_denom", "Incorrect denom")
}

func TestUpdateVestingInfo(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	govAddress := sdk.AccAddress(address.Module("gov"))

	// Define the test data
	signer := govAddress.String()
	vestingInfo := types.VestingInfo{
		BaseDenom:      "test_denom",
		VestingDenom:   "test_denom",
		NumBlocks:      10,
		VestNowFactor:  sdkmath.NewInt(10),
		NumMaxVestings: 10,
	}

	params := keeper.GetParams(ctx)
	params.VestingInfos = append(params.VestingInfos, vestingInfo)
	keeper.SetParams(ctx, params)

	vestingInfoResult, _ := keeper.GetVestingInfo(ctx, "test_denom")
	assert.Equal(t, vestingInfoResult.BaseDenom, "test_denom", "Incorrect denom")

	// Call the UpdateVestingInfo function
	msg := types.MsgUpdateVestingInfo{
		Authority:      signer,
		BaseDenom:      "test_denom",
		VestingDenom:   "test_denom",
		NumBlocks:      10,
		VestNowFactor:  10,
		NumMaxVestings: 10,
	}
	_, err := msgServer.UpdateVestingInfo(ctx, &msg)
	require.NoError(t, err)
}

// TestKeeper_UpdateVestingInfoWithWrongGovAddress tests the UpdateVestingInfo function with wrong gov address
func TestKeeper_UpdateVestingInfoWithWrongGovAddress(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Define the test data
	signer := sdk.AccAddress(address.Module("test")).String()

	// Call the UpdateVestingInfo function
	msg := types.MsgUpdateVestingInfo{
		Authority:      signer,
		BaseDenom:      "test_denom",
		VestingDenom:   "test_denom",
		NumBlocks:      10,
		VestNowFactor:  10,
		NumMaxVestings: 10,
	}
	_, err := msgServer.UpdateVestingInfo(ctx, &msg)
	require.Error(t, err)
}

// TestKeeper_UpdateVestingInfoWithNegativeNumBlocks tests the UpdateVestingInfo function with negative num blocks
func TestKeeper_UpdateVestingInfoWithNegativeNumBlocks(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	govAddress := sdk.AccAddress(address.Module("gov"))

	// Define the test data
	signer := govAddress.String()

	// Call the UpdateVestingInfo function
	msg := types.MsgUpdateVestingInfo{
		Authority:      signer,
		BaseDenom:      "test_denom",
		VestingDenom:   "test_denom",
		NumBlocks:      -10,
		VestNowFactor:  10,
		NumMaxVestings: 10,
	}
	_, err := msgServer.UpdateVestingInfo(ctx, &msg)
	require.Error(t, err)
}
