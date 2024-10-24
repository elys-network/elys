package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

// TestKeeper_ClaimVesting tests the ClaimVesting function
func TestKeeper_ClaimVesting(t *testing.T) {
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	// set block height
	ctx = ctx.WithBlockHeight(10)

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

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")
	acc := app.AccountKeeper.GetAccount(ctx, creator)
	if acc == nil {
		acc = app.AccountKeeper.NewAccountWithAddress(ctx, creator)
		app.AccountKeeper.SetAccount(ctx, acc)
	}
	// Create a claim vesting message
	claimVestMsg := &types.MsgClaimVesting{
		Sender: creator.String(),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:         ptypes.Elys,
				TotalAmount:   sdk.NewInt(100),
				ClaimedAmount: sdk.NewInt(1),
				NumBlocks:     100,
				StartBlock:    0,
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Execute the CancelVest function
	_, err := msgServer.ClaimVesting(ctx, claimVestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were updated correctly
	newCommitments := keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not updated correctly")
	require.Equal(t, sdk.NewInt(100), newCommitments.VestingTokens[0].TotalAmount, "total amount was not updated correctly")
	require.Equal(t, sdk.NewInt(10), newCommitments.VestingTokens[0].ClaimedAmount, "claimed amount was not updated correctly")
}
