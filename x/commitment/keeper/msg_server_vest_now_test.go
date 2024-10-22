package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"

	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKeeper_VestNow tests the VestNow function with VestNowEnabled set to false
func TestVestNowDisabled(t *testing.T) {
	app := app.InitElysTestApp(true)

	// Disable VestNow for test
	commitmentkeeper.VestNowEnabled = false

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)
	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := ptypes.Eden

	// Test scenario: VestNow should fail
	msg := &types.MsgVestNow{
		Creator: creator,
		Amount:  sdk.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}

// TestKeeper_VestNow tests the VestNow function with invalid denom
func TestVestNowInvalidDenom(t *testing.T) {
	app := app.InitElysTestApp(true)

	// Enable VestNow for test
	commitmentkeeper.VestNowEnabled = true

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)
	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := "invalid_denom"

	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdk.ZeroInt(),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
	}

	keeper.SetParams(ctx, params)

	// Test scenario: VestNow should fail
	msg := &types.MsgVestNow{
		Creator: creator,
		Amount:  sdk.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}

// TestKeeper_VestNow tests the VestNow function with vest now factor set to zero
func TestVestNowInvalidAmount(t *testing.T) {
	app := app.InitElysTestApp(true)

	// Enable VestNow for test
	commitmentkeeper.VestNowEnabled = true

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)
	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := ptypes.Eden
	initialClaimed := sdk.NewInt(5000)
	initialCommitted := sdk.NewInt(10000)

	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdk.ZeroInt(),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
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
		Amount:  sdk.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}

func TestVestNow(t *testing.T) {
	app := app.InitElysTestApp(true)

	// Enable VestNow for test
	commitmentkeeper.VestNowEnabled = true

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)
	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := ptypes.Eden
	initialClaimed := sdk.NewInt(5000)
	initialCommitted := sdk.NewInt(10000)

	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdk.NewInt(90),
			NumMaxVestings: 10,
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
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
		Amount:  sdk.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.NoError(t, err)

	updatedCommitments := keeper.GetCommitments(ctx, creatorAddr)

	claimedBalance := updatedCommitments.GetClaimedForDenom(denom)
	assert.Equal(t, sdk.NewInt(2000), claimedBalance)

	// Check if the vested tokens were received
	creatorBalance := app.BankKeeper.GetBalance(ctx, creatorAddr, vestingInfos[0].VestingDenom)
	require.Equal(t, sdk.NewInt(33), creatorBalance.Amount, "tokens were not vested correctly")

	// Test scenario 2: Withdraw more than claimed balance but within total balance
	msg = &types.MsgVestNow{
		Creator: creator,
		Amount:  sdk.NewInt(7000),
		Denom:   denom,
	}

	_, err = msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}
