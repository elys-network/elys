package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v4/app"

	assetprofiletypes "github.com/elys-network/elys/v4/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/v4/x/commitment/keeper"
	"github.com/elys-network/elys/v4/x/commitment/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUncommitTokens(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, ctx)
	simapp.SetupAssetProfile(app, ctx)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := "testdenom"
	initialCommitted := sdkmath.NewInt(100)
	uncommitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: initialCommitted,
	}

	params := app.CommitmentKeeper.GetParams(ctx)
	params.TotalCommitted = sdk.NewCoins(sdk.NewCoin(denom, initialCommitted))
	app.CommitmentKeeper.SetParams(ctx, params)

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

	// Call the UncommitTokens function
	err = keeper.UncommitTokens(ctx, addr[0], denom, uncommitAmount, false)
	require.NoError(t, err)

	// Check if the committed tokens have been added to the store
	commitments := keeper.GetCommitments(ctx, sdk.MustAccAddressFromBech32(creator))

	// Check if the committed tokens have the expected values
	assert.Equal(t, len(commitments.CommittedTokens), 0, "Incorrect creator")

	rewardUnclaimed := sdk.NewCoins(sdk.NewCoin(denom, uncommitAmount))

	edenCoin := app.BankKeeper.GetBalance(ctx, addr[0], denom)
	require.Equal(t, sdk.Coins{edenCoin}, rewardUnclaimed)
}

// TestUncommitTokensDenomNotFound tests the UncommitTokens function when the asset profile entry for the denom is not found
func TestUncommitTokensDenomNotFound(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(false)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := "testdenom"
	uncommitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: uncommitAmount,
	}

	initialCommitments := types.Commitments{
		Creator:         creator,
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Call the UncommitTokens function
	err := keeper.UncommitTokens(ctx, addr[0], denom, uncommitAmount, false)
	require.Error(t, err)
}

// TestUncommitTokensWithdrawDisabled tests the UncommitTokens function when the withdraw is disabled for the denom
func TestUncommitTokensWithdrawDisabled(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := "testdenom"
	uncommitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: uncommitAmount,
	}

	initialCommitments := types.Commitments{
		Creator:         creator,
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: denom, CommitEnabled: true, WithdrawEnabled: false})

	// Call the UncommitTokens function
	err := keeper.UncommitTokens(ctx, addr[0], denom, uncommitAmount, false)
	require.Error(t, err)
}

// TestUncommitTokensEdenBTriggersHookError tests the UncommitTokens function when the denom is EdenB and the hook returns an error
func TestUncommitTokensEdenBTriggersHookError(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := ptypes.EdenB
	uncommitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: uncommitAmount,
	}

	initialCommitments := types.Commitments{
		Creator:         creator,
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: denom, CommitEnabled: true, WithdrawEnabled: true})

	// Call the UncommitTokens function
	err := keeper.UncommitTokens(ctx, addr[0], denom, uncommitAmount, false)
	require.Error(t, err)
}

// TestMsgServerUncommitTokensNoDelegationError tests the UncommitTokens function in the keeper through the MsgServer triggers error
func TestMsgServerUncommitTokensNoDelegationError(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := ptypes.Eden
	uncommitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: uncommitAmount,
	}

	initialCommitments := types.Commitments{
		Creator:         creator,
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: denom, CommitEnabled: true, WithdrawEnabled: true})

	// Set up the MsgServer
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Call the UncommitTokens function through the MsgServer
	_, err := msgServer.UncommitTokens(ctx, &types.MsgUncommitTokens{
		Creator: creator,
		Denom:   denom,
		Amount:  uncommitAmount,
	})
	require.Error(t, err)
}

// TestMsgServerUncommitTokensUnsupportedUncommitTokenError tests the UncommitTokens function in the keeper through the MsgServer triggers error
func TestMsgServerUncommitTokensUnsupportedUncommitTokenError(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := "testdenom"
	uncommitAmount := sdkmath.NewInt(100)

	// Set up initial commitments object with sufficient unclaimed & committed tokens
	committedTokens := types.CommittedTokens{
		Denom:  denom,
		Amount: uncommitAmount,
	}

	initialCommitments := types.Commitments{
		Creator:         creator,
		CommittedTokens: []*types.CommittedTokens{&committedTokens},
	}

	keeper.SetCommitments(ctx, initialCommitments)

	// Set up the MsgServer
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Call the UncommitTokens function through the MsgServer
	_, err := msgServer.UncommitTokens(ctx, &types.MsgUncommitTokens{
		Creator: creator,
		Denom:   denom,
		Amount:  uncommitAmount,
	})
	require.Error(t, err)
}

// TestMsgServerUncommitTokensInvalidAddressError tests the UncommitTokens function in the keeper through the MsgServer triggers error
func TestMsgServerUncommitTokensInvalidAddressError(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Set up the MsgServer
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Call the UncommitTokens function through the MsgServer
	_, err := msgServer.UncommitTokens(ctx, &types.MsgUncommitTokens{
		Creator: "invalid_address",
		Denom:   ptypes.Eden,
		Amount:  sdkmath.NewInt(100),
	})
	require.Error(t, err)
}
