package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"

	aptypes "github.com/elys-network/elys/x/assetprofile/types"
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

	creatorAddr, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Define the test data
	creator := creatorAddr.String()
	denom := ptypes.Eden
	initialUnclaimed := sdk.NewInt(50)
	rewardsByElysUnclaimed := sdk.NewInt(1)
	rewardsByEdenUnclaimed := sdk.NewInt(1)
	rewardsByEdenbUnclaimed := sdk.NewInt(1)
	rewardsByUsdcUnclaimed := sdk.NewInt(1)

	initialCommitted := sdk.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: initialCommitted,
	}

	initialCommitments := types.Commitments{
		Creator: creator,
		RewardsUnclaimed: sdk.Coins{sdk.Coin{
			Denom:  denom,
			Amount: initialUnclaimed,
		}},
		RewardsByElysUnclaimed: sdk.Coins{sdk.Coin{
			Denom:  denom,
			Amount: rewardsByElysUnclaimed,
		}},
		RewardsByEdenUnclaimed: sdk.Coins{sdk.Coin{
			Denom:  denom,
			Amount: rewardsByEdenUnclaimed,
		}},
		RewardsByEdenbUnclaimed: sdk.Coins{sdk.Coin{
			Denom:  denom,
			Amount: rewardsByEdenbUnclaimed,
		}},
		RewardsByUsdcUnclaimed: sdk.Coins{sdk.Coin{
			Denom:  denom,
			Amount: rewardsByUsdcUnclaimed,
		}},
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{Denom: denom, BaseDenom: denom, WithdrawEnabled: true})

	// Test scenario 1: Withdraw LP mining rewards
	err := app.CommitmentKeeper.RecordClaimReward(ctx, creator, denom, sdk.NewInt(46), types.EarnType_ALL_PROGRAM)
	require.NoError(t, err)

	updatedCommitments := keeper.GetCommitments(ctx, creator)
	unclaimedBalance := updatedCommitments.GetLPMiningSubBucketRewardUnclaimedForDenom(denom)
	assert.Equal(t, sdk.NewInt(4), unclaimedBalance)
	require.Equal(t, "46ueden", updatedCommitments.Claimed.String(), "tokens were not claimed correctly")

	// Test scenario 2: Withdraw within unclaimed balance
	err = app.CommitmentKeeper.RecordClaimReward(ctx, creator, denom, sdk.NewInt(2), types.EarnType_ALL_PROGRAM)
	require.NoError(t, err)

	updatedCommitments = keeper.GetCommitments(ctx, creator)
	unclaimedBalance = updatedCommitments.GetRewardUnclaimedForDenom(denom)
	assert.Equal(t, sdk.NewInt(2), unclaimedBalance)
	require.Equal(t, "48ueden", updatedCommitments.Claimed.String(), "tokens were not claimed correctly")

	// Test scenario 3: Withdraw more than unclaimed reward
	err = app.CommitmentKeeper.RecordClaimReward(ctx, creator, denom, sdk.NewInt(70), types.EarnType_ALL_PROGRAM)
	require.Error(t, err)
}
