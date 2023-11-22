package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCancelVest(t *testing.T) {
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

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
		Amount:  sdk.NewInt(25),
	}

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		VestingTokens: []*types.VestingTokens{
			{
				Denom:           ptypes.Eden,
				TotalAmount:     sdk.NewInt(100),
				UnvestedAmount:  sdk.NewInt(100),
				EpochIdentifier: epochstypes.DayEpochID,
				NumEpochs:       100,
				CurrentEpoch:    0,
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Execute the CancelVest function
	_, err := msgServer.CancelVest(ctx, cancelVestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were updated correctly
	newCommitments := keeper.GetCommitments(ctx, cancelVestMsg.Creator)
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not updated correctly")
	require.Equal(t, sdk.NewInt(75), newCommitments.VestingTokens[0].TotalAmount, "total amount was not updated correctly")
	require.Equal(t, sdk.NewInt(75), newCommitments.VestingTokens[0].UnvestedAmount, "unvested amount was not updated correctly")
	// check if the unclaimed tokens were updated correctly
	require.Equal(t, sdk.NewInt(25), newCommitments.GetClaimedForDenom(ptypes.Eden))

	// Try to cancel an amount that exceeds the unvested amount
	cancelVestMsg.Amount = sdk.NewInt(101)
	_, err = msgServer.CancelVest(ctx, cancelVestMsg)
	require.Error(t, err, "should throw an error when trying to cancel more tokens than available")
	require.True(t, types.ErrInsufficientVestingTokens.Is(err), "error should be insufficient vesting tokens")
}
