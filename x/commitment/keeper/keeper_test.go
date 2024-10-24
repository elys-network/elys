package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKeeper_SpendableCoins tests the SpendableCoins function
func TestKeeper_SpendableCoins(t *testing.T) {
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

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdk.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdk.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Test SpendableCoins
	spendableCoins := keeper.SpendableCoins(ctx, creator)
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(100)), sdk.NewCoin(ptypes.Elys, sdk.NewInt(1000000))), spendableCoins)
}

// TestKeeper_AddEdenEdenBOnAccount tests the AddEdenEdenBOnAccount function
func TestKeeper_AddEdenEdenBOnAccount(t *testing.T) {
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

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdk.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdk.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Test AddEdenEdenBOnAccount
	_ = keeper.AddEdenEdenBOnAccount(ctx, creator, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(50)), sdk.NewCoin(ptypes.EdenB, sdk.NewInt(50))))

	// Check the updated commitments
	commitments = keeper.GetCommitments(ctx, creator)
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(200)), sdk.NewCoin(ptypes.EdenB, sdk.NewInt(50))), commitments.Claimed)
}

// TestKeeper_AddEdenEdenBOnModule tests the AddEdenEdenBOnModule function
func TestKeeper_AddEdenEdenBOnModule(t *testing.T) {
	app := simapp.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	creator := addr[0]

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(100)))

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
				Amount: sdk.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdk.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Test AddEdenEdenBOnModule
	_ = keeper.AddEdenEdenBOnModule(
		ctx, moduleName,
		sdk.NewCoins(
			sdk.NewCoin(ptypes.Eden, sdk.NewInt(50)),
			sdk.NewCoin(ptypes.EdenB, sdk.NewInt(150)),
		),
	)

	// Check the updated commitments
	commitments = keeper.GetCommitments(ctx, creator)
	assert.Equal(t,
		sdk.NewCoins(
			sdk.NewCoin(ptypes.Eden, sdk.NewInt(150)),
		),
		commitments.Claimed,
	)
}

// TestKeeper_SubEdenEdenBOnModule tests the SubEdenEdenBOnModule function with insufficient claimed tokens error
func TestKeeper_SubEdenEdenBOnModule_InsufficientClaimedTokens(t *testing.T) {
	app := simapp.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	creator := addr[0]

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(100)))

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
				Amount: sdk.NewInt(50),
			},
			{
				Denom:  ptypes.EdenB,
				Amount: sdk.NewInt(50),
			},
		},
		Claimed: sdk.Coins{
			sdk.NewCoin(ptypes.Eden, sdk.NewInt(150)),
			sdk.NewCoin(ptypes.EdenB, sdk.NewInt(150)),
		},
	}
	keeper.SetCommitments(ctx, commitments)

	// Test SubEdenEdenBOnModule
	_, err = keeper.SubEdenEdenBOnModule(
		ctx, moduleName,
		sdk.NewCoins(
			sdk.NewCoin(ptypes.Eden, sdk.NewInt(50)),
			sdk.NewCoin(ptypes.EdenB, sdk.NewInt(50)),
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
