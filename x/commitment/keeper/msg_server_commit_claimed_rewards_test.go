package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/app"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/v5/x/commitment/keeper"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitClaimedRewardsWithEden(t *testing.T) {
	// Create a test context and keeper
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Define the test data
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()
	denom := ptypes.Eden
	initialUnclaimed := sdkmath.NewInt(500)
	commitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed tokens
	rewardsClaimed := sdk.NewCoin(denom, initialUnclaimed)
	initialCommitments := types.Commitments{
		Creator: creator,
		Claimed: sdk.Coins{rewardsClaimed},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: denom, CommitEnabled: true})

	// Call the CommitClaimedRewards function
	msg := types.MsgCommitClaimedRewards{
		Creator: creator,
		Amount:  commitAmount,
		Denom:   denom,
	}
	_, err := msgServer.CommitClaimedRewards(ctx, &msg)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments := keeper.GetCommitments(ctx, sdk.MustAccAddressFromBech32(creator))
	assert.Equal(t, creator, commitments.Creator, "Incorrect creator")
	assert.Len(t, commitments.CommittedTokens, 1, "Incorrect number of committed tokens")
	assert.Equal(t, denom, commitments.CommittedTokens[0].Denom, "Incorrect denom")
	assert.Equal(t, commitAmount, commitments.CommittedTokens[0].Amount, "Incorrect amount")
}

func TestCommitClaimedRewardsWithEdenB(t *testing.T) {
	// Create a test context and keeper
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Define the test data
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()
	denom := ptypes.EdenB
	initialUnclaimed := sdkmath.NewInt(500)
	commitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed tokens
	rewardsClaimed := sdk.NewCoin(denom, initialUnclaimed)
	initialCommitments := types.Commitments{
		Creator: creator,
		Claimed: sdk.Coins{rewardsClaimed},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: denom, CommitEnabled: true})

	// Call the CommitClaimedRewards function
	msg := types.MsgCommitClaimedRewards{
		Creator: creator,
		Amount:  commitAmount,
		Denom:   denom,
	}
	_, err := msgServer.CommitClaimedRewards(ctx, &msg)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments := keeper.GetCommitments(ctx, sdk.MustAccAddressFromBech32(creator))
	assert.Equal(t, creator, commitments.Creator, "Incorrect creator")
	assert.Len(t, commitments.CommittedTokens, 1, "Incorrect number of committed tokens")
	assert.Equal(t, denom, commitments.CommittedTokens[0].Denom, "Incorrect denom")
	assert.Equal(t, commitAmount, commitments.CommittedTokens[0].Amount, "Incorrect amount")
}

// TestCommitClaimedRewardsWithInvalidDenom tests the CommitClaimedRewards function with an invalid denom
func TestCommitClaimedRewardsWithInvalidDenom(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Create a new account
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()

	// Create a commit claimed rewards message with an invalid denom
	msg := types.MsgCommitClaimedRewards{
		Creator: creator,
		Amount:  sdkmath.NewInt(100),
		Denom:   "invalid",
	}

	// Execute the CommitClaimedRewards function
	_, err := msgServer.CommitClaimedRewards(ctx, &msg)
	require.Error(t, err, "should throw an error when using an invalid denom")
	require.True(t, assetprofiletypes.ErrAssetProfileNotFound.Is(err), "error should be asset profile not found")
}

// TestCommitClaimedRewardsWithEdenDisabled tests the CommitClaimedRewards function with Eden disabled
func TestCommitClaimedRewardsWithEdenDisabled(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Create a new account
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: false})

	// Create a commit claimed rewards message with Eden
	msg := types.MsgCommitClaimedRewards{
		Creator: creator,
		Amount:  sdkmath.NewInt(100),
		Denom:   ptypes.Eden,
	}

	// Execute the CommitClaimedRewards function
	_, err := msgServer.CommitClaimedRewards(ctx, &msg)
	require.Error(t, err, "should throw an error when Eden is disabled")
	require.True(t, types.ErrCommitDisabled.Is(err), "error should be commit disabled")
}
