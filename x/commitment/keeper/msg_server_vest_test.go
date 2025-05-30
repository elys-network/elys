package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/app"
	commitmentkeeper "github.com/elys-network/elys/v6/x/commitment/keeper"
	"github.com/elys-network/elys/v6/x/commitment/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestVest(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
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

	// Create a vesting message
	vestMsg := &types.MsgVest{
		Creator: creator.String(),
		Denom:   ptypes.Eden,
		Amount:  sdkmath.NewInt(100),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(150),
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Execute the Vest function
	_, err := msgServer.Vest(ctx, vestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were added to commitments
	newCommitments := keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not added")

	// Check if the claimed tokens were updated correctly
	claimed := newCommitments.GetClaimedForDenom(vestMsg.Denom)
	require.Equal(t, sdkmath.NewInt(50), claimed, "claimed tokens were not updated correctly")

	// Check if the committed tokens were updated correctly
	committedToken := newCommitments.GetCommittedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdkmath.NewInt(50), committedToken, "committed tokens were not updated correctly")

	_, err = msgServer.Vest(ctx, vestMsg)
	require.Error(t, err)
}

func TestExceedVesting(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
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

	vestingInfos := []types.VestingInfo{
		{
			BaseDenom:      ptypes.Eden,
			VestingDenom:   ptypes.Elys,
			NumBlocks:      10,
			VestNowFactor:  sdkmath.NewInt(90),
			NumMaxVestings: 1,
		},
	}

	params := types.Params{
		VestingInfos: vestingInfos,
	}

	keeper.SetParams(ctx, params)

	// Create a vesting message
	vestMsg := &types.MsgVest{
		Creator: creator.String(),
		Denom:   ptypes.Eden,
		Amount:  sdkmath.NewInt(100),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(150),
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Execute the Vest function
	_, err := msgServer.Vest(ctx, vestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were added to commitments
	newCommitments := keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not added")

	// Check if the claimed tokens were updated correctly
	claimed := newCommitments.GetClaimedForDenom(vestMsg.Denom)
	require.Equal(t, sdkmath.NewInt(50), claimed, "claimed tokens were not updated correctly")

	// Check if the committed tokens were updated correctly
	committedToken := newCommitments.GetCommittedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdkmath.NewInt(50), committedToken, "committed tokens were not updated correctly")

	// Execute the Vest again and it should endfunction
	_, err = msgServer.Vest(ctx, vestMsg)
	require.EqualError(t, err, errorsmod.Wrapf(types.ErrExceedMaxVestings, "creator: %s", vestMsg.Creator).Error(), "exceed vesting not worked properly")
}
