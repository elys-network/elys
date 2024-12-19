package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	simapp "github.com/elys-network/elys/app"

	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper_Stake(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(false)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := ptypes.Elys
	initialCommitted := sdkmath.NewInt(100)
	uncommitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: initialCommitted,
	}

	initialCommitments := types.Commitments{
		Creator:         creator,
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: denom, CommitEnabled: true, WithdrawEnabled: true})

	// Add coins on commitment module
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(denom, initialCommitted)})
	require.NoError(t, err)

	// Call the Stake function
	_, err = msgServer.Stake(ctx, &types.MsgStake{
		Creator:          creator,
		Asset:            denom,
		Amount:           uncommitAmount,
		ValidatorAddress: "cosmosvaloper1x8efhljzvs52u5xa6m7crcwes7v9u0nlwdgw30",
	})
	require.Error(t, err)
}

func TestKeeper_Stake_commit(t *testing.T) {
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

	// Call the Stake function
	_, err := msgServer.Stake(ctx, &types.MsgStake{
		Creator:          "creator",
		Asset:            denom,
		Amount:           commitAmount,
		ValidatorAddress: "",
	})
	require.Error(t, err)

	_, err = msgServer.Stake(ctx, &types.MsgStake{
		Creator:          creator,
		Asset:            denom,
		Amount:           commitAmount.Add(sdkmath.NewInt(1000000)),
		ValidatorAddress: "",
	})
	require.Error(t, err)

	_, err = msgServer.Stake(ctx, &types.MsgStake{
		Creator:          creator,
		Asset:            denom,
		Amount:           commitAmount,
		ValidatorAddress: "",
	})
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments := keeper.GetCommitments(ctx, sdk.MustAccAddressFromBech32(creator))
	assert.Equal(t, creator, commitments.Creator, "Incorrect creator")
	assert.Len(t, commitments.CommittedTokens, 1, "Incorrect number of committed tokens")
	assert.Equal(t, denom, commitments.CommittedTokens[0].Denom, "Incorrect denom")
	assert.Equal(t, commitAmount, commitments.CommittedTokens[0].Amount, "Incorrect amount")
}
