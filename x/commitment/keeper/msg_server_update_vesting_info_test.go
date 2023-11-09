package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/elys-network/elys/app"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddVestingInfo(t *testing.T) {
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	govAddress := sdk.AccAddress(address.Module("gov"))

	// Define the test data
	signer := govAddress.String()

	// Call the Update Vesting Info function
	msg := types.MsgUpdateVestingInfo{
		Authority:       signer,
		BaseDenom:       "test_denom",
		VestingDenom:    "test_denom",
		EpochIdentifier: "day",
		NumEpochs:       10,
		VestNowFactor:   10,
		NumMaxVestings:  10,
	}
	_, err := msgServer.UpdateVestingInfo(ctx, &msg)
	require.NoError(t, err)

	// Check if the vesting info has been added to the store
	vestingInfo, _ := keeper.GetVestingInfo(ctx, "test_denom")
	assert.Equal(t, vestingInfo.BaseDenom, "test_denom", "Incorrect denom")
}

func TestUpdateVestingInfo(t *testing.T) {
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	govAddress := sdk.AccAddress(address.Module("gov"))

	// Define the test data
	signer := govAddress.String()
	vestingInfo := &types.VestingInfo{
		BaseDenom:       "test_denom",
		VestingDenom:    "test_denom",
		EpochIdentifier: "month",
		NumEpochs:       10,
		VestNowFactor:   sdk.NewInt(10),
		NumMaxVestings:  10,
	}

	params := keeper.GetParams(ctx)
	params.VestingInfos = append(params.VestingInfos, vestingInfo)
	keeper.SetParams(ctx, params)

	vestingInfo, _ = keeper.GetVestingInfo(ctx, "test_denom")
	assert.Equal(t, vestingInfo.BaseDenom, "test_denom", "Incorrect denom")

	// Call the UpdateVestingInfo function
	msg := types.MsgUpdateVestingInfo{
		Authority:       signer,
		BaseDenom:       "test_denom",
		VestingDenom:    "test_denom",
		EpochIdentifier: "day",
		NumEpochs:       10,
		VestNowFactor:   10,
		NumMaxVestings:  10,
	}
	_, err := msgServer.UpdateVestingInfo(ctx, &msg)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	vestingInfo, _ = keeper.GetVestingInfo(ctx, "test_denom")
	assert.Equal(t, vestingInfo.EpochIdentifier, "day", "Incorrect epoch identifier")
}
