package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	simapp "github.com/elys-network/elys/v4/app"
	"github.com/elys-network/elys/v4/x/commitment/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKeeper_SpendableCoins tests the SpendableCoins function
func TestKeeper_SpendableCoins(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], edenToken)
	require.NoError(t, err)

	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Test SpendableCoins
	spendableCoins := keeper.SpendableCoins(ctx, creator)
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000000))), spendableCoins)
}

// TestKeeper_AddEdenEdenBOnAccount tests the AddEdenEdenBOnAccount function
func TestKeeper_AddEdenEdenBOnAccount(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], edenToken)
	require.NoError(t, err)

	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Test AddEdenEdenBOnAccount
	_, _ = keeper.AddEdenEdenBOnAccount(ctx, creator, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)), sdk.NewCoin(ptypes.EdenB, sdkmath.NewInt(50))))

	// Check the updated commitments
	commitments = keeper.GetCommitments(ctx, creator)
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(200)), sdk.NewCoin(ptypes.EdenB, sdkmath.NewInt(50))), commitments.Claimed)
}

// TestKeeper_AddEdenEdenBOnModule tests the AddEdenEdenBOnModule function
func TestKeeper_AddEdenEdenBOnModule(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	creator := addr[0]

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, edenToken)
	require.NoError(t, err)

	moduleName := types.ModuleName

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Test AddEdenEdenBOnModule
	_, _ = keeper.AddEdenEdenBOnModule(
		ctx, moduleName,
		sdk.NewCoins(
			sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)),
			sdk.NewCoin(ptypes.EdenB, sdkmath.NewInt(150)),
		),
	)

	// Check the updated commitments
	commitments = keeper.GetCommitments(ctx, creator)
	assert.Equal(t,
		sdk.NewCoins(
			sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(150)),
		),
		commitments.Claimed,
	)
}

// TestKeeper_SubEdenEdenBOnModule_InsufficientClaimedTokens tests the SubEdenEdenBOnModule function with insufficient claimed tokens error
func TestKeeper_SubEdenEdenBOnModule_InsufficientClaimedTokens(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	creator := addr[0]

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, edenToken)
	require.NoError(t, err)

	moduleName := types.ModuleName

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: authtypes.NewModuleAddress(moduleName).String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
			{
				Denom:  ptypes.EdenB,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{
			sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(150)),
			sdk.NewCoin(ptypes.EdenB, sdkmath.NewInt(150)),
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Test SubEdenEdenBOnModule
	_, _, err = keeper.SubEdenEdenBOnModule(
		ctx, moduleName,
		sdk.NewCoins(
			sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)),
			sdk.NewCoin(ptypes.EdenB, sdkmath.NewInt(50)),
		),
	)
	require.NoError(t, err)

	// Check the updated commitments
	commitments = keeper.GetCommitments(ctx, creator)
	assert.Equal(t,
		sdk.NewCoins(),
		commitments.Claimed,
	)
}

// TestKeeper_Logger tests the Logger function
func TestKeeper_Logger(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	logger := app.Logger()

	keeper.Logger(ctx).Info("test")
	logger.Info("test")
}

// TestKeeper_BankKeeper tests the BankKeeper function
func TestKeeper_BankKeeper(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	keeper := app.CommitmentKeeper

	assert.NotNil(t, keeper.BankKeeper())
}

// TestKeeper_SetHooks_Panic tests the SetHooks function with a nil argument
func TestKeeper_SetHooks_Panic(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	keeper := app.CommitmentKeeper

	assert.Panics(t, func() {
		keeper.SetHooks(nil)
	})
}

// TestKeeper_MintCoins tests the MintCoins function
func TestKeeper_MintCoins(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Mint 100ueden and uelys
	tokens := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000000)))

	err := keeper.MintCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Check the updated commitments
	commitments := keeper.GetCommitments(ctx, addr[0])
	assert.Equal(t, sdk.NewCoins(), commitments.Claimed)
}

// TestKeeper_MintCoins_EmptyAmount tests the MintCoins function with an empty amount
func TestKeeper_MintCoins_EmptyAmount(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Mint 0ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(0)))

	err := keeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
}

// TestKeeper_BurnCoins tests the BurnCoins function
func TestKeeper_BurnCoins(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Mint 100ueden and uelys
	tokens := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000000)))

	err := keeper.MintCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Burn 50ueden and uelys
	tokens = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(500000)))

	err = keeper.BurnCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Check the updated commitments
	commitments := keeper.GetCommitments(ctx, addr[0])
	assert.Equal(t, sdk.NewCoins(), commitments.Claimed)
}

// TestKeeper_BurnCoins_EmptyAmount tests the BurnCoins function with empty amount
func TestKeeper_BurnCoins_EmptyAmount(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Mint 100ueden and uelys
	tokens := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000000)))

	err := keeper.MintCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Burn 50ueden and uelys
	tokens = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)))

	err = keeper.BurnCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Check the updated commitments
	commitments := keeper.GetCommitments(ctx, addr[0])
	assert.Equal(t, sdk.NewCoins(), commitments.Claimed)
}

// TestKeeper_SendCoinsFromModuleToModule tests the SendCoinsFromModuleToModule function
func TestKeeper_SendCoinsFromModuleToModule(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 2, sdkmath.NewInt(1000000))

	// Mint 100ueden and uelys
	tokens := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000000)))

	err := keeper.MintCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Send 50ueden and uelys from module to module
	tokens = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(500000)))

	err = keeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.ModuleName, tokens)
	require.NoError(t, err)

	// Check the updated commitments
	commitments := keeper.GetCommitments(ctx, addr[0])
	assert.Equal(t, sdk.NewCoins(), commitments.Claimed)
}

// TestKeeper_SendCoinsFromModuleToModule_EmptyAmount tests the SendCoinsFromModuleToModule function with empty amount
func TestKeeper_SendCoinsFromModuleToModule_EmptyAmount(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 2, sdkmath.NewInt(1000000))

	// Mint 100ueden and uelys
	tokens := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000000)))

	err := keeper.MintCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Send 50ueden and uelys from module to module
	tokens = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)))

	err = keeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.ModuleName, tokens)
	require.NoError(t, err)

	// Check the updated commitments
	commitments := keeper.GetCommitments(ctx, addr[0])
	assert.Equal(t, sdk.NewCoins(), commitments.Claimed)
}

// TestKeeper_SendCoinsFromModuleToAccount tests the SendCoinsFromModuleToAccount function
func TestKeeper_SendCoinsFromModuleToAccount(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 2, sdkmath.NewInt(1000000))

	// Mint 100ueden and uelys
	tokens := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000000)))

	err := keeper.MintCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Send 50ueden and uelys from module to account
	tokens = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(500000)))

	err = keeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], tokens)
	require.NoError(t, err)

	// Check the updated commitments
	commitments := keeper.GetCommitments(ctx, addr[0])
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50))), commitments.Claimed)
}

// TestKeeper_SendCoinsFromModuleToAccount_EmptyAmount tests the SendCoinsFromModuleToAccount function with empty amount
func TestKeeper_SendCoinsFromModuleToAccount_EmptyAmount(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 2, sdkmath.NewInt(1000000))

	// Mint 100ueden and uelys
	tokens := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000000)))

	err := keeper.MintCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	// Send 50ueden and uelys from module to account
	tokens = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50)))

	err = keeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], tokens)
	require.NoError(t, err)

	// Check the updated commitments
	commitments := keeper.GetCommitments(ctx, addr[0])
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(50))), commitments.Claimed)
}

// TestKeeper_SendCoinsFromAccountToModule tests the SendCoinsFromAccountToModule function
func TestKeeper_SendCoinsFromAccountToModule(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 2, sdkmath.NewInt(1000000))

	// Send 50uelys from account to module
	tokens := sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(50)))

	err := keeper.SendCoinsFromAccountToModule(ctx, addr[0], types.ModuleName, tokens)
	require.NoError(t, err)
}
