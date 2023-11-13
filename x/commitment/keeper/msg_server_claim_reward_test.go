package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"

	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaimReward(t *testing.T) {
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := ptypes.Eden
	initialUnclaimed := sdk.NewInt(50)
	initialCommitted := sdk.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	rewardsUnclaimed := sdk.Coin{
		Denom:  denom,
		Amount: initialUnclaimed,
	}

	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: initialCommitted,
	}

	initialCommitments := types.Commitments{
		Creator:          creator,
		RewardsUnclaimed: sdk.Coins{rewardsUnclaimed},
		CommittedTokens:  []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: denom, WithdrawEnabled: true})

	// Test scenario 1: Withdraw within unclaimed balance
	msg := &types.MsgClaimReward{
		Creator: creator,
		Amount:  sdk.NewInt(30),
		Denom:   denom,
	}

	_, err := msgServer.ClaimReward(ctx, msg)
	require.NoError(t, err)

	updatedCommitments, found := keeper.GetCommitments(ctx, creator)
	require.True(t, found)

	unclaimedBalance := updatedCommitments.GetRewardUnclaimedForDenom(denom)
	assert.Equal(t, sdk.NewInt(20), unclaimedBalance)

	// Check if the withdrawn tokens were received
	require.Equal(t, "30ueden", updatedCommitments.Claimed.String(), "tokens were not claimed correctly")

	// Test scenario 2: Withdraw more than unclaimed reward
	msg = &types.MsgClaimReward{
		Creator: creator,
		Amount:  sdk.NewInt(70),
		Denom:   denom,
	}

	_, err = msgServer.ClaimReward(ctx, msg)
	require.Error(t, err)
}
