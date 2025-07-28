package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v7/app"

	assetprofiletypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/v7/x/commitment/keeper"
	"github.com/elys-network/elys/v7/x/commitment/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_UnstakeWithElys(t *testing.T) {
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
	_, err = msgServer.Unstake(ctx, &types.MsgUnstake{
		Creator:          creator,
		Asset:            denom,
		Amount:           uncommitAmount,
		ValidatorAddress: "cosmosvaloper1x8efhljzvs52u5xa6m7crcwes7v9u0nlwdgw30",
	})
	require.Error(t, err)
}

// TestKeeper_UnstakeWithElysNegativeAmount tests the Unstake function with elys tokens and negative amount
func TestKeeper_UnstakeWithElysNegativeAmount(t *testing.T) {
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
	_, err = msgServer.Unstake(ctx, &types.MsgUnstake{
		Creator:          "invalid_address",
		Asset:            denom,
		Amount:           uncommitAmount,
		ValidatorAddress: "cosmosvaloper1x8efhljzvs52u5xa6m7crcwes7v9u0nlwdgw30",
	})
	require.Error(t, err)
}

// TestKeeper_UnstakeWithUncommitInvalidValidatorAddress tests the Unstake function with uncommitting tokens and invalid validator address
func TestKeeper_UnstakeWithUncommitInvalidValidatorAddress(t *testing.T) {
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
	_, err = msgServer.Unstake(ctx, &types.MsgUnstake{
		Creator:          creator,
		Asset:            denom,
		Amount:           uncommitAmount,
		ValidatorAddress: "invalid_validator_address",
	})
	require.Error(t, err)
}

func TestKeeper_UnstakeWithUncommit(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(false)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := ptypes.Eden
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
	_, err = msgServer.Unstake(ctx, &types.MsgUnstake{
		Creator:          creator,
		Asset:            denom,
		Amount:           uncommitAmount,
		ValidatorAddress: "cosmosvaloper1x8efhljzvs52u5xa6m7crcwes7v9u0nlwdgw30",
	})
	require.Error(t, err)
}

// TestKeeper_UnstakeWithUncommitNegativeAmount tests the Unstake function with uncommitting tokens and negative amount
func TestKeeper_UnstakeWithUncommitNegativeAmount(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(false)

	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Generate 1 random account with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Define the test data
	creator := addr[0].String()
	denom := ptypes.Eden
	initialCommitted := sdkmath.NewInt(100)
	uncommitAmount := sdkmath.NewInt(-100)

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
	_, err = msgServer.Unstake(ctx, &types.MsgUnstake{
		Creator:          creator,
		Asset:            denom,
		Amount:           uncommitAmount,
		ValidatorAddress: "cosmosvaloper1x8efhljzvs52u5xa6m7crcwes7v9u0nlwdgw30",
	})
	require.Error(t, err)
}
