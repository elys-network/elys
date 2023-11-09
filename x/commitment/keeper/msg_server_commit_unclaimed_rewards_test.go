package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitUnclaimedRewards(t *testing.T) {
	// Create a test context and keeper
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	// Define the test data
	creator := "test_creator"
	denom := "test_denom"
	initialUnclaimed := sdk.NewInt(500)
	commitAmount := sdk.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed tokens
	rewardsUnclaimed := types.RewardsUnclaimed{
		Denom:  denom,
		Amount: initialUnclaimed,
	}

	initialCommitments := types.Commitments{
		Creator:          creator,
		RewardsUnclaimed: []*types.RewardsUnclaimed{&rewardsUnclaimed},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: denom, CommitEnabled: true})

	// Call the CommitUnclaimedRewards function
	msg := types.MsgCommitUnclaimedRewards{
		Creator: creator,
		Amount:  commitAmount,
		Denom:   denom,
	}
	_, err := msgServer.CommitUnclaimedRewards(ctx, &msg)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments, found := keeper.GetCommitments(ctx, creator)
	assert.True(t, found, "Commitments not found in the store")

	// Check if the committed tokens have the expected values
	assert.Equal(t, creator, commitments.Creator, "Incorrect creator")
	assert.Len(t, commitments.CommittedTokens, 1, "Incorrect number of committed tokens")
	assert.Equal(t, denom, commitments.CommittedTokens[0].Denom, "Incorrect denom")
	assert.Equal(t, commitAmount, commitments.CommittedTokens[0].Amount, "Incorrect amount")
}
