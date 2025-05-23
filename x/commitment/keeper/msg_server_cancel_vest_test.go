package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/app"
	commitmentkeeper "github.com/elys-network/elys/v5/x/commitment/keeper"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCancelVest_CompleteAmount(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

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
		VestingInfos: vestingInfos,
	}

	keeper.SetParams(ctx, params)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}
	// Create a cancel vesting message
	cancelVestMsg := &types.MsgCancelVest{
		Creator: creator.String(),
		Denom:   ptypes.Eden,
		Amount:  sdkmath.NewInt(100),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:         ptypes.Elys,
				TotalAmount:   sdkmath.NewInt(100),
				ClaimedAmount: sdkmath.NewInt(0),
				NumBlocks:     100,
				StartBlock:    0,
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Increase the block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 25)

	// From 100 tokens, some of tokens are already vested and that will be claimed,
	// so CancelVest should cancel the remaining amount without any issue
	// Execute the CancelVest function
	_, err := msgServer.CancelVest(ctx, cancelVestMsg)
	require.NoError(t, err)

	newCommitments := keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 0, "vesting tokens should be empty after cancelling all remaining amount")

	// No vesting tokens, so should throw an error
	_, err = msgServer.CancelVest(ctx, cancelVestMsg)
	require.Error(t, err)
	require.True(t, types.ErrInsufficientVestingTokens.Is(err), "Error should be insufficient vesting tokens")
}

func TestCancelVest(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

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
		VestingInfos: vestingInfos,
	}

	keeper.SetParams(ctx, params)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}
	// Create a cancel vesting message
	cancelVestMsg := &types.MsgCancelVest{
		Creator: creator.String(),
		Denom:   ptypes.Eden,
		Amount:  sdkmath.NewInt(25),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:         ptypes.Elys,
				TotalAmount:   sdkmath.NewInt(100),
				ClaimedAmount: sdkmath.NewInt(0),
				NumBlocks:     100,
				StartBlock:    0,
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Increase the block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 25)

	// Execute the CancelVest function
	_, err := msgServer.CancelVest(ctx, cancelVestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were updated correctly
	newCommitments := keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not updated correctly")
	require.Equal(t, sdkmath.NewInt(50), newCommitments.VestingTokens[0].TotalAmount, "total amount was not updated correctly")
	require.Equal(t, sdkmath.NewInt(0), newCommitments.VestingTokens[0].ClaimedAmount, "claimed amount was not updated correctly")
	// check if the unclaimed tokens were updated correctly
	require.Equal(t, sdkmath.NewInt(25), newCommitments.GetClaimedForDenom(ptypes.Eden))

	// Try to cancel an amount that exceeds the unvested amount, should cancel all the remaining amount
	cancelVestMsg.Amount = sdkmath.NewInt(100)
	_, err = msgServer.CancelVest(ctx, cancelVestMsg)
	require.NoError(t, err)
	newCommitments = keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 0, "vesting tokens should be empty after cancelling all remaining amount")
}

func TestCancelVest_WithPreviousClaimed(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

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
		VestingInfos: vestingInfos,
	}

	keeper.SetParams(ctx, params)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}
	// Create a cancel vesting message
	cancelVestMsg := &types.MsgCancelVest{
		Creator: creator.String(),
		Denom:   ptypes.Eden,
		Amount:  sdkmath.NewInt(25),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:         ptypes.Elys,
				TotalAmount:   sdkmath.NewInt(100),
				ClaimedAmount: sdkmath.NewInt(20),
				NumBlocks:     100,
				StartBlock:    0,
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)
	// Increase the block height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 25)

	// Execute the CancelVest function
	_, err := msgServer.CancelVest(ctx, cancelVestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were updated correctly
	newCommitments := keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not updated correctly")
	// vested so far already claimed before cancelling
	require.Equal(t, sdkmath.NewInt(50), newCommitments.VestingTokens[0].TotalAmount, "total amount was not updated correctly")
	require.Equal(t, sdkmath.NewInt(0), newCommitments.VestingTokens[0].ClaimedAmount, "claimed amount was not updated correctly")
	// check if the unclaimed tokens were updated correctly
	require.Equal(t, sdkmath.NewInt(25), newCommitments.GetClaimedForDenom(ptypes.Eden))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 25)
	cancelVestMsg.Amount = sdkmath.NewInt(20)

	// Execute the CancelVest function
	_, err = msgServer.CancelVest(ctx, cancelVestMsg)
	require.NoError(t, err)

	newCommitments = keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not updated correctly")
	require.Equal(t, sdkmath.NewInt(14), newCommitments.VestingTokens[0].TotalAmount, "total amount was not updated correctly")
	require.Equal(t, sdkmath.NewInt(0), newCommitments.VestingTokens[0].ClaimedAmount, "claimed amount was not updated correctly")
	require.Equal(t, int64(50), newCommitments.VestingTokens[0].NumBlocks, "NumBlocks not updated correctly")
	// check if the unclaimed tokens were updated correctly, 25 from previous cancel + 20 from this cancel
	require.Equal(t, sdkmath.NewInt(45), newCommitments.GetClaimedForDenom(ptypes.Eden))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 25)

	// Try to cancel total remaining amount
	cancelVestMsg.Amount = sdkmath.NewInt(7)
	_, err = msgServer.CancelVest(ctx, cancelVestMsg)
	require.NoError(t, err)

	newCommitments = keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 0, "vesting tokens should be empty after cancelling all remaining amount")

	// check if the unclaimed tokens were updated correctly, 25+20 from previous cancel + 7 from this cancel
	require.Equal(t, sdkmath.NewInt(52), newCommitments.GetClaimedForDenom(ptypes.Eden))
}

// TestCancelVestIncorrectDenom tests the CancelVest function with an incorrect denom
func TestCancelVestIncorrectDenom(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}
	// Create a cancel vesting message
	cancelVestMsg := &types.MsgCancelVest{
		Creator: creator.String(),
		Denom:   "incorrect",
		Amount:  sdkmath.NewInt(25),
	}

	// Execute the CancelVest function
	_, err := msgServer.CancelVest(ctx, cancelVestMsg)
	require.Error(t, err, "should throw an error when trying to cancel tokens with an incorrect denom")
	require.True(t, types.ErrInvalidDenom.Is(err), "error should be invalid denom")
}

// TestCancelVestNoVestingInfo tests the CancelVest function with no vesting info
func TestCancelVestNoVestingInfo(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}
	// Create a cancel vesting message
	cancelVestMsg := &types.MsgCancelVest{
		Creator: creator.String(),
		Denom:   ptypes.Eden,
		Amount:  sdkmath.NewInt(25),
	}

	// Execute the CancelVest function
	_, err := msgServer.CancelVest(ctx, cancelVestMsg)
	require.Error(t, err, "should throw an error when trying to cancel tokens with no vesting info")
	fmt.Println(err.Error())
	require.True(t, types.ErrInsufficientVestingTokens.Is(err), "Error should be insufficient vesting tokens")
}
