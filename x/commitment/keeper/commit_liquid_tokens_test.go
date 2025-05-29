package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/app"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	"github.com/elys-network/elys/v6/x/commitment/types"
	"github.com/stretchr/testify/require"
)

func TestCommitLiquidTokens(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Set assetprofile entry for denom
	denom := "amm/pool/1"
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: denom, CommitEnabled: true})

	// Add initial funds to creator's account
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewInt(200)))
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, coins)
	require.NoError(t, err)
	balance := app.BankKeeper.GetBalance(ctx, creator, denom)
	require.Equal(t, coins.AmountOf(denom), balance.Amount, "creator balance did not initialize")

	require.NoError(t, err)

	// Execute the DepositTokens function
	err = keeper.CommitLiquidTokens(ctx, creator, denom, sdkmath.NewInt(100), 0)
	require.NoError(t, err)

	// Check if the tokens were deposited and unclaimed balance was updated
	commitments := keeper.GetCommitments(ctx, creator)

	committedBalance := commitments.GetCommittedAmountForDenom(denom)
	require.Equal(t, sdkmath.NewInt(100), committedBalance, "committed balance did not update correctly")

	// Check if the deposited tokens were deducted from creator balance
	remainingCoins := app.BankKeeper.GetBalance(ctx, creator, denom)
	require.Equal(t, sdkmath.NewInt(100), remainingCoins.Amount, "tokens were not deducted correctly")
}
