package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitClaimedRewards(t *testing.T) {
	// Create a test context and keeper
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	// Define the test data
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()
	denom := "test_denom"
	initialUnclaimed := sdk.NewInt(500)
	commitAmount := sdk.NewInt(100)

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
