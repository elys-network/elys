package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/app"

	commitmentkeeper "github.com/elys-network/elys/v5/x/commitment/keeper"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKeeper_VestNow tests the VestNow function with VestNowEnabled set to false
func TestVestNowDisabled(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)
	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := ptypes.Eden

	// Test scenario: VestNow should fail
	msg := &types.MsgVestNow{
		Creator: creator,
		Amount:  sdkmath.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}

// TestKeeper_VestNow tests the VestNow function with invalid denom
func TestVestNowInvalidDenom(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)
	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := "invalid_denom"

	vestingInfos := []types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdkmath.ZeroInt(),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos:  vestingInfos,
		EnableVestNow: true,
	}

	keeper.SetParams(ctx, params)

	// Test scenario: VestNow should fail
	msg := &types.MsgVestNow{
		Creator: creator,
		Amount:  sdkmath.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}

// TestKeeper_VestNow tests the VestNow function with vest now factor set to zero
func TestVestNowInvalidAmount(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)
	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := ptypes.Eden
	initialClaimed := sdkmath.NewInt(5000)
	initialCommitted := sdkmath.NewInt(10000)

	vestingInfos := []types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdkmath.ZeroInt(),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos:  vestingInfos,
		EnableVestNow: true,
	}

	keeper.SetParams(ctx, params)

	// Set up initial commitments object with sufficient claimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: initialCommitted,
	}

	initialCommitments := types.Commitments{
		Creator: creator,
		Claimed: sdk.Coins{
			{
				Denom:  denom,
				Amount: initialClaimed,
			},
		},
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Test scenario: VestNow should fail
	msg := &types.MsgVestNow{
		Creator: creator,
		Amount:  sdkmath.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}

func TestVestNow(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)
	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := ptypes.Eden
	initialClaimed := sdkmath.NewInt(5000)
	initialCommitted := sdkmath.NewInt(10000)

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
		VestingInfos:  vestingInfos,
		EnableVestNow: true,
	}

	keeper.SetParams(ctx, params)

	// Set up initial commitments object with sufficient claimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: initialCommitted,
	}

	initialCommitments := types.Commitments{
		Creator: creator,
		Claimed: sdk.Coins{
			{
				Denom:  denom,
				Amount: initialClaimed,
			},
		},
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Test scenario 1: Withdraw within claimed balance
	msg := &types.MsgVestNow{
		Creator: creator,
		Amount:  sdkmath.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.NoError(t, err)

	updatedCommitments := keeper.GetCommitments(ctx, creatorAddr)

	claimedBalance := updatedCommitments.GetClaimedForDenom(denom)
	assert.Equal(t, sdkmath.NewInt(2000), claimedBalance)

	// Check if the vested tokens were received
	creatorBalance := app.BankKeeper.GetBalance(ctx, creatorAddr, vestingInfos[0].VestingDenom)
	require.Equal(t, sdkmath.NewInt(33), creatorBalance.Amount, "tokens were not vested correctly")

	// Test scenario 2: Withdraw more than claimed balance but within total balance
	msg = &types.MsgVestNow{
		Creator: creator,
		Amount:  sdkmath.NewInt(7000),
		Denom:   denom,
	}

	_, err = msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}
