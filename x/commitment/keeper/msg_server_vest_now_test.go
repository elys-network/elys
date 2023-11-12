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
	initialUnclaimed := sdk.NewInt(5000)
	initialCommitted := sdk.NewInt(10000)

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

	// Test scenario 1: Withdraw within unclaimed balance
	msg := &types.MsgVestNow{
		Creator: creator,
		Amount:  sdk.NewInt(3000),
		Denom:   denom,
	}

	_, err := msgServer.VestNow(ctx, msg)
	require.NoError(t, err)

	updatedCommitments, found := keeper.GetCommitments(ctx, creator)
	require.True(t, found)

	unclaimedBalance := updatedCommitments.GetRewardUnclaimedForDenom(denom)
	assert.Equal(t, sdk.NewInt(2000), unclaimedBalance)
	// Check if the vested tokens were received
	creatorBalance := app.BankKeeper.GetBalance(ctx, creatorAddr, vestingInfos[0].VestingDenom)
	require.Equal(t, sdk.NewInt(33), creatorBalance.Amount, "tokens were not vested correctly")

	// Test scenario 2: Withdraw more than unclaimed balance but within total balance
	msg = &types.MsgVestNow{
		Creator: creator,
		Amount:  sdk.NewInt(7000),
		Denom:   denom,
	}

	_, err = msgServer.VestNow(ctx, msg)
	require.NoError(t, err)

	updatedCommitments, found = keeper.GetCommitments(ctx, creator)
	require.True(t, found)

	unclaimedBalance = updatedCommitments.GetRewardUnclaimedForDenom(denom)
	assert.Equal(t, sdk.NewInt(0), unclaimedBalance)

	committedBalance := updatedCommitments.GetCommittedAmountForDenom(denom)
	assert.Equal(t, sdk.NewInt(5000), committedBalance)

	// Check if the vested tokens were received
	creatorBalance = app.BankKeeper.GetBalance(ctx, creatorAddr, vestingInfos[0].VestingDenom)

	require.Equal(t, sdk.NewInt(110), creatorBalance.Amount, "tokens were not vested correctly")

	// Test scenario 3: Withdraw more than total balance
	msg = &types.MsgVestNow{
		Creator: creator,
		Amount:  sdk.NewInt(10000),
		Denom:   denom,
	}

	_, err = msgServer.VestNow(ctx, msg)
	require.Error(t, err)
}
