package keeper_test

import (
	"fmt"
	"testing"

	errorsmod "cosmossdk.io/errors"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	simapp "github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestLiquidVestWithExceed(t *testing.T) {
	app := simapp.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], edenToken)
	require.NoError(t, err)

	creator := addr[0]
	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)
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

	// Create a vesting message
	vestMsg := &types.MsgVestLiquid{
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
		Claimed: sdk.Coins{
			{
				Denom:  ptypes.Eden,
				Amount: sdk.NewInt(150),
			},
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true})

	// Execute the Vest function
	_, err = msgServer.VestLiquid(ctx, vestMsg)
	require.NoError(t, err)

	// Check if the vesting tokens were added to commitments
	newCommitments := keeper.GetCommitments(ctx, creator)
	require.Len(t, newCommitments.VestingTokens, 1, "vesting tokens were not added")

	// Check if the claimed tokens were updated correctly
	claimed := newCommitments.GetClaimedForDenom(vestMsg.Denom)
	require.Equal(t, sdk.NewInt(150).String(), claimed.String(), "claimed tokens were not updated correctly")

	// Check if the committed tokens were updated correctly
	committedToken := newCommitments.GetCommittedAmountForDenom(vestMsg.Denom)
	require.Equal(t, sdk.NewInt(50).String(), committedToken.String(), "committed tokens were not updated correctly")

	edenCoin := app.BankKeeper.GetBalance(ctx, addr[0], ptypes.Eden)
	require.Equal(t, edenCoin.Amount, sdk.ZeroInt())

	_, err = msgServer.VestLiquid(ctx, vestMsg)
	require.Equal(t, err.Error(), errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, fmt.Sprintf("unable to send deposit tokens: %v", edenToken)).Error())
}
