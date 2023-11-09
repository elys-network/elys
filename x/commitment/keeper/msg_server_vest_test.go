package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/app"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestVest(t *testing.T) {
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}

	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:       ptypes.Eden,
			VestingDenom:    ptypes.Elys,
			EpochIdentifier: "tenseconds",
			NumEpochs:       10,
			VestNowFactor:   sdk.NewInt(90),
			NumMaxVestings:  10,
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
		Amount:  sdk.NewInt(100),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdk.NewInt(50),
			},
		},
		RewardsUnclaimed: []*types.RewardsUnclaimed{
			{
				Denom:  ptypes.Eden,
				Amount: sdk.NewInt(150),
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Execute the Vest function
	_, err := msgServer.Vest(ctx, vestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were added to commitments
	newCommitments, found := keeper.GetCommitments(ctx, vestMsg.Creator)
	require.True(t, found, "commitments not found")
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not added")

	// Check if the unclaimed tokens were updated correctly
	rewardUnclaimed := newCommitments.GetUnclaimedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdk.NewInt(50), rewardUnclaimed, "unclaimed tokens were not updated correctly")

	// Check if the committed tokens were updated correctly
	committedToken := newCommitments.GetCommittedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdk.NewInt(50), committedToken, "committed tokens were not updated correctly")

	_, err = msgServer.Vest(ctx, vestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were added to commitments
	newCommitments, found = keeper.GetCommitments(ctx, vestMsg.Creator)
	require.True(t, found, "commitments not found")
	require.Len(t, newCommitments.VestingTokens, 2, "vesting tokens were not added")

	// Check if the unclaimed tokens were updated correctly
	rewardUnclaimed = newCommitments.GetUnclaimedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdk.NewInt(0), rewardUnclaimed, "unclaimed tokens were not updated correctly")

	// Check if the committed tokens were updated correctly
	committedToken = newCommitments.GetCommittedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdk.NewInt(0), committedToken, "committed tokens were not updated correctly")
}

func TestExceedVesting(t *testing.T) {
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}

	vestingInfos := []*types.VestingInfo{
		{
			BaseDenom:       ptypes.Eden,
			VestingDenom:    ptypes.Elys,
			EpochIdentifier: "tenseconds",
			NumEpochs:       10,
			VestNowFactor:   sdk.NewInt(90),
			NumMaxVestings:  1,
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
		Amount:  sdk.NewInt(100),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdk.NewInt(50),
			},
		},
		RewardsUnclaimed: []*types.RewardsUnclaimed{
			{
				Denom:  ptypes.Eden,
				Amount: sdk.NewInt(150),
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Execute the Vest function
	_, err := msgServer.Vest(ctx, vestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were added to commitments
	newCommitments, found := keeper.GetCommitments(ctx, vestMsg.Creator)
	require.True(t, found, "commitments not found")
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not added")

	// Check if the unclaimed tokens were updated correctly
	rewardUnclaimed := newCommitments.GetUnclaimedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdk.NewInt(50), rewardUnclaimed, "unclaimed tokens were not updated correctly")

	// Check if the committed tokens were updated correctly
	committedToken := newCommitments.GetCommittedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdk.NewInt(50), committedToken, "committed tokens were not updated correctly")

	// Execute the Vest again and it should endfunction
	_, err = msgServer.Vest(ctx, vestMsg)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrExceedMaxVestings, "creator: %s", vestMsg.Creator).Error(), "exceed vesting not worked properly")
}
